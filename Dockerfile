FROM golang:1.22-bookworm

RUN apt-get update && \
    apt-get install -y ffmpeg mkvtoolnix && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o bakasub-server cmd/server/main.go

EXPOSE 8080

CMD ["./bakasub-server"]