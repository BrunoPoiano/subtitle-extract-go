FROM golang:1.24.2-alpine3.21

RUN apk add --no-cache dcron 
RUN apk add --no-cache tzdata

RUN mkdir /subextract
RUN mkdir /subextract/videos
WORKDIR /subextract

COPY start.sh main.go go.mod ./
RUN go mod download
RUN go build . 

