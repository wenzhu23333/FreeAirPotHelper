# Free AirPot Helper

Build a Free AirPot and Subscribe URL Based On YouTube Live:

https://www.youtube.com/watch?v=qmRkvKo-KbQ

## Usage

### Build with Docker
```shell
docker build -t freeap .
```

### Config
Make a new Directory
```shell
mkdir docker_config
```
Put **Clash Yaml (Rename as config_clash.yaml)** and **Country.mmdb** in it. (This is the airpot for access YouTube.)

Define the config_freeap.yaml you want.

### Deployment
```shell
docker run -d -p 80:8080 -v ./docker_config:/tmp/workdir/configs --name freeap freeap
```

### Access

Subscribe URL: http://127.0.0.1/sub?token=xx

Get Yaml for clash: 

```shell
curl -O http://127.0.0.1/sub?token=xx -o xxx.yml
```