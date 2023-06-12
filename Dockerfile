# Build Stage
FROM golang:1.18-alpine as builder

WORKDIR /app

ENV GOPROXY https://goproxy.cn

COPY . .

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o freeAP

# Production Stage
FROM jrottenberg/ffmpeg as final

WORKDIR /tmp/workdir/

RUN apt-get update && apt-get install -y curl

RUN curl -O -L https://github.com/tindy2013/subconverter/releases/download/v0.7.2/subconverter_linux64.tar.gz

RUN tar -xzvf subconverter_linux64.tar.gz && rm subconverter_linux64.tar.gz

RUN chmod -R +x ./subconverter

RUN curl -O -L https://github.com/Dreamacro/clash/releases/download/v1.16.0/clash-linux-amd64-v1.16.0.gz

RUN gzip -d clash-linux-amd64-v1.16.0.gz && mv clash-linux-amd64-v1.16.0 clash

RUN chmod +x ./clash

#RUN curl -o Country.mmdb -L https://git.io/GeoLite2-Country.mmdb

RUN mkdir configs

COPY --from=builder /app/freeAP /tmp/workdir/

RUN chmod +x ./freeAP

COPY entrypoint.sh /tmp/workdir/

RUN chmod +x ./entrypoint.sh

ENTRYPOINT ["sh", "entrypoint.sh"]


