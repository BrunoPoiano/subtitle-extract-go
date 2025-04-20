FROM golang:1.24.2-alpine3.21

WORKDIR /app/subextract
COPY start.sh main.go go.mod ./
RUN go mod download
RUN go build .

RUN apk add --no-cache dcron
RUN apk add --no-cache tzdata
RUN apk add --no-cache ffmpeg

WORKDIR /app/subextract
RUN mkdir -p /app/subextract/videos

RUN echo " * */2 * * * /app/subextract/main >> /var/log/cron.log 2>&1" > /etc/crontabs/root

# Create empty log file for cron output
RUN touch /var/log/cron.log

# Copy startup script and make it executable
COPY start.sh /app/subextract/start.sh
RUN chmod +x /app/subextract/start.sh

# Execute the startup script when container launches
CMD ["/bin/sh", "/app/subextract/start.sh"]
