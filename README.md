
```bash
# go.mod 作成
docker run --rm -v ${PWD}:/server -w /server golang:1.25 sh -c "go mod init sns-4-illustrators-server"
```

```bash
# コード内で利用するライブラリの自動追加・排除
docker run --rm -v ${PWD}:/server -w /server golang:1.25 sh -c "go mod tidy"
```

```bash
# .air.toml 作成
docker run --rm -v ${PWD}:/server -w /server golang:1.25 sh -c "go install github.com/air-verse/air@latest && air init"
```

```bash
docker-compose up --build
```