package main

import (
	"flag"
	"fmt"
	"github.com/kkdai/youtube/v2"
	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/qrcode"
	"image"
	_ "image/jpeg"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

var queue *Queue

func main() {
	configPath := flag.String("f", "config_freeap.yaml", "Config File Path")
	flag.Parse()
	config := readConfig(*configPath)
	fmt.Println("Config: ")
	fmt.Println("----------------------")
	queueCapacity := config.QueueCapacity
	fmt.Println("Queue Capacity: " + strconv.Itoa(queueCapacity))
	videoID := config.VideoID
	fmt.Println("Video ID: " + videoID)
	subconverterPrefix := config.SubconverterPrefix
	fmt.Println("Subconverter Prefix: " + subconverterPrefix)
	nodeNum := config.NodeNum
	fmt.Println("Node Num: " + strconv.Itoa(nodeNum))
	token := config.Token
	fmt.Println("Token: " + token)
	port := config.Port
	fmt.Println("Port: " + strconv.Itoa(port))
	//certFile := config.CertFile
	//keyFile := config.KeyFile
	//fmt.Println("Cert File:" + certFile)
	//fmt.Println("Key File: " + keyFile)
	videoQuality := config.VideoQuality
	fmt.Println("Video Quality: " + videoQuality)
	subscribeURL := config.SubscribeURL
	fmt.Println("Subscribe URL: " + subscribeURL)
	subUpdateInterval := config.SubUpdateInterval
	nodeUpdateInterval := config.NodeUpdateInterval
	fmt.Println("Node Update Interval (s): " + strconv.Itoa(nodeUpdateInterval))
	fmt.Println("Subscribe Update Interval (s): " + strconv.Itoa(subUpdateInterval))
	fmt.Println("----------------------")
	queue = NewQueue(queueCapacity)

	go createSubscribeWorker(videoID, nodeNum, subconverterPrefix, videoQuality, nodeUpdateInterval, subUpdateInterval)

	mux := http.NewServeMux()
	mux.HandleFunc(subscribeURL, subscribeHandler(token))

	server := &http.Server{
		Addr:           fmt.Sprintf(":%v", port),
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	err := server.ListenAndServe()
	if err != nil {
		//log.Println(err)
		fmt.Println(err)
	}
}

func subscribeHandler(token string) http.HandlerFunc {
	return func(w http.ResponseWriter, request *http.Request) {
		tokenURL := request.URL.Query().Get("token")
		if tokenURL != token {
			_, err := fmt.Fprintln(w, "Error !")
			if err != nil {
				return
			}
			return
		}
		subURL, _ := queue.Head()
		if subURL == "" {
			_, err := fmt.Fprintln(w, "Wait For A Moment!")
			if err != nil {
				fmt.Println(err)
			}
		} else {
			resp, err := http.Get(subURL)
			if err != nil {
				fmt.Println(err)
			}
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				fmt.Println(err)
			}
			w.Header().Set("Profile-Update-Interval", "24")
			subscribeYaml := string(body)
			//w.Header().Set("Content-Length", fmt.Sprint(len(subscribeYaml)))
			w.Header().Set("Access-Control-Allow-Origin", "*")
			_, err = fmt.Fprintln(w, subscribeYaml)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}

func createSubscribeWorker(videoID string, nodeNum int, subconverterPrefix string, videoQuality string,
	nodeUpdateInterval int, subUpdateInterval int) {
	for {
		strs := make([]string, 0, nodeNum)
		for i := 0; i < nodeNum; i++ {
			if err := GetLiveFileMP4(videoID, videoQuality); err != nil {
				fmt.Println(err)
			}
			proxyUrl, err := GetQRCodeResult("temp.mp4", "1", "output.jpg")
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println("Node " + strconv.Itoa(i+1) + ": " + proxyUrl)
			fmt.Println("")
			strs = append(strs, proxyUrl)
			time.Sleep(time.Second * time.Duration(nodeUpdateInterval))
		}
		joinUrlEncoded := url.QueryEscape(strings.Join(strs, "|"))
		subconverterUrl := subconverterPrefix + joinUrlEncoded
		//channel <- subscribeUrl
		queue.Enqueue(subconverterUrl)
		err := ioutil.WriteFile("subscribe.txt", []byte(subconverterUrl), 0644)
		if err != nil {
			fmt.Println("写入文件时发生错误:", err)
			return
		}
		fmt.Println("订阅写入成功")
		//break
		time.Sleep(time.Second * time.Duration(subUpdateInterval))
	}
}

func GetLiveFileMP4(videoID string, videoQuality string) error {
	client := youtube.Client{}

	video, err := client.GetVideo(videoID)
	if err != nil {
		return err
	}
	formats := video.Formats.Quality(videoQuality)
	stream, _, err := client.GetStream(video, &formats[0])
	if err != nil {
		return err
	}

	file, err := os.Create("temp.mp4")
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, stream)
	if err != nil {
		return err
	}
	return nil
}

func GetQRCodeResult(filename string, frameIndex string, QRCodePath string) (string, error) {
	_, err := os.Stat(QRCodePath)
	if err == nil {
		// 文件存在
		err := exec.Command("rm", QRCodePath).Run()
		if err != nil {
			fmt.Println("Failed to delete frame:", err)
			return "", err
		}
	}
	vfMessage := fmt.Sprintf("select=eq(n\\,%v)", frameIndex)
	cmd := exec.Command("ffmpeg", "-i", filename, "-vf", vfMessage, "-vframes", "1", QRCodePath)
	if err := cmd.Run(); err != nil {
		fmt.Println("Failed to extract frame:", err)
		return "", err
	}
	file, err := os.Open(QRCodePath)
	if err != nil {
		fmt.Println("Failed to open output file:", err)
		return "", err
	}
	//defer file.Close()
	img, _, err := image.Decode(file)

	if err != nil {
		fmt.Println("Failed to Decode:", err)
		return "", err
	}

	// prepare BinaryBitmap
	bmp, err := gozxing.NewBinaryBitmapFromImage(img)

	if err != nil {
		fmt.Println("Failed to Prepare:", err)
		return "", err
	}

	// decode image
	qrReader := qrcode.NewQRCodeReader()
	result, err := qrReader.Decode(bmp, nil)
	//fmt.Println(result.GetText())

	if err != nil {
		fmt.Println("Failed to Read QR:", err)
		return "", err
	}

	return result.GetText(), nil

}
