FROM golang:1.25

WORKDIR /server

# ホットリロードツール（Air）のインストール
RUN go install github.com/air-verse/air@latest

# pprofがグラフを生成するのに必要な描画エンジンのインストール
RUN apt-get update && apt-get install -y graphviz

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Airを起動
CMD ["air", "-c", ".air.toml"]