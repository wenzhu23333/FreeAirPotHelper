package main

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

type Config struct {
	Port               int    `yaml:"port"`
	SubconverterPrefix string `yaml:"subconverterPrefix"`
	VideoID            string `yaml:"videoID"`
	QueueCapacity      int    `yaml:"queueCapacity"`
	NodeNum            int    `yaml:"nodeNum"`
	Token              string `yaml:"token"`
	VideoQuality       string `yaml:"videoQuality"`
	SubscribeURL       string `yaml:"subscribeURL"`
	NodeUpdateInterval int    `yaml:"nodeUpdateInterval"`
	SubUpdateInterval  int    `yaml:"subUpdateInterval"`

	//CertFile        string `yaml:"certFile"`
	//KeyFile         string `yaml:"keyFile"`
}

func readConfig(filename string) *Config {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println(err)
	}

	// 解析配置文件
	config := &Config{}
	err = yaml.Unmarshal(data, config)
	if err != nil {
		fmt.Println(err)
	}
	return config
}
