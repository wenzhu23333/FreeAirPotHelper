nohup ./subconverter/subconverter > subconverter_logs.txt &
cp ./configs/Country.mmdb ./Country.mmdb
nohup ./clash -d . -f ./configs/config_clash.yaml > clash_logs.txt &
export http_proxy=http://127.0.0.1:7890 && export https_proxy=http://127.0.0.1:7890
sleep 20
./freeAP -f ./configs/config_freeap.yaml > freeap_logs.txt