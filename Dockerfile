FROM golang:1.25

WORKDIR /server

# ホットリロードツール（Air）のインストール
RUN go install github.com/air-verse/air@latest

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Airを起動
CMD ["air", "-c", ".air.toml"]