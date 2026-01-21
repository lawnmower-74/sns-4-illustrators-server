# よく使うコマンド

```bash
# go.mod 作成
docker run --rm -v ${PWD}:/server -w /server golang:1.25 sh -c "go mod init sns-4-illustrators-server"
```

```bash
# コード内で利用するライブラリの自動追加・削除
docker run --rm -v ${PWD}:/server -w /server golang:1.25 sh -c "go mod tidy"
```

```bash
# .air.toml 作成
docker run --rm -v ${PWD}:/server -w /server golang:1.25 sh -c "go install github.com/air-verse/air@latest && air init"
```

```bash
# 通常起動
docker-compose up --build
```

```bash
# Airを介さずにmain.goを実行したい場合
docker-compose run --rm server go run main.go
```

```bash
# buildのみを実行し、成否を確認したい場合
docker-compose run --rm server go build -o test_bin .
```

```bash
# 起動中のコンテナ内で pprof を実行（パフォーマンス・負荷状況をグラフ形式で確認できる）
docker exec -it コンテナID go tool pprof -http=0.0.0.0:8081 http://localhost:8080/debug/pprof/profile?seconds=10

# コンテナIDは以下で確認可能
docker ps

# 1. seconds 秒内にアップロードなど計測したい処理を実行
# 2. http://localhost:8081 にアクセスし、ダウンロード
# 3. https://www.speedscope.app/ にアップ（こっちのほうが見やすい）
```


# 開発時に発生したエラーと対処

## 1. Airが古いバイナリを読み込んでしまい、変更が反映されない

### 原因

Go（v1.18～）では、ビルド時にGit情報を埋め込む仕様になっている  
→ Docker（Linux）側からホスト（Windows）の .git フォルダを読み取る際に、所有権エラー（Exit 128）が発生  

これにより、go build が内部で失敗（Exit 1）し、新しい実行ファイルが生成されなかった。  

結果、  
Airは「最新のビルドが失敗した」場合、プロセスを止めず、前回成功した時点の古いバイナリをそのまま動かし続けたため、変更が反映されなかった。

#### 補足（ホットリロードの流れ）

1. コード変更

2. Airが検知し、build を実行。`./tmp/main` が最新版に上書きされる

    ```.air.tmol
    [build]
    cmd = "go build -o ./tmp/main ."
    ```

3. 自動で再起動されるため、いちいち「コンテナ停止・再起動」しなくていい

### 特定方法

```bash
# buildのみを実行し、原因を特定
$ docker-compose run --rm server go build -o test_bin .

[+] Running 1/1
 ✔ Container 5d8815513f4fc87bf267ebddb0e3e3354cd35d00a8cc4e592c9374cce36c6b95-db-1  Started                                                0.2s 
error obtaining VCS status: exit status 128           # VCS（git）の入手ステータス = 128（所有権エラー）
        Use -buildvcs=false to disable VCS stamping.  # 回避法
```

### 対処

```.air.tmol
cmd = "go build -buildvcs=false -o ./tmp/main ."
```

とし、Git情報の参照をスキップさせることでエラーを回避


## 2. PostmanからのPOSTが通らない

### 原因

Postman（Windows）から Docker（Linux）への IPv6 接続ができていない

### 特定方法

```powershell
# ※PostmanはWindowsOS上で動くアプリのため、PowerShellで確認
> resolve-dnsname localhost

Name                                           Type   TTL   Section    IPAddress
----                                           ----   ---   -------    ---------
localhost                                      AAAA   1200  Question   ::1
localhost                                      A      1200  Question   127.0.0.1
```

WindowsOSでは `localhost = ::1(IPv6)` として優先的に解釈している  
つまり、  
Postman → Docker は IPv6 で接続しようとしている（実際は接続できずに応答なし）  

ちなみに、  

```bash
$ curl -v http://localhost:8080

* Host localhost:8080 was resolved.
* IPv6: ::1
* IPv4: 127.0.0.1
*   Trying [::1]:8080...
* Connected to localhost (::1) port 8080
> GET / HTTP/1.1
> Host: localhost:8080
> User-Agent: curl/8.5.0
> Accept: */*
>
< HTTP/1.1 404 Not Found
< Content-Type: text/plain
< Date: Tue, 20 Jan 2026 09:06:29 GMT
< Content-Length: 18
<
* Connection #0 to host localhost left intact
404 page not found
```

ここから WSL → Docker の同じLinux環境同士だと、IPv6 で接続できていることがわかる

次に、Windows でどのポートを何が使っているのか確認する

```powershell
# ※Windowsのプロセスを確認するため、PowerShellで確認
> netstat -ano | findstr :8080
  TCP         0.0.0.0:8080           0.0.0.0:0              LISTENING       25296
  TCP         127.0.0.1:8080         127.0.0.1:62577        ESTABLISHED     25296
  TCP         127.0.0.1:8080         127.0.0.1:64711        ESTABLISHED     25296
  TCP         127.0.0.1:62577        127.0.0.1:8080         ESTABLISHED     44544
  TCP         127.0.0.1:64711        127.0.0.1:8080         ESTABLISHED     44544
  TCP         [::]:8080              [::]:0                 LISTENING       25296
  TCP         [::1]:8080             [::]:0                 LISTENING       24584
  TCP         [::1]:8080             [::]:0                 LISTENING       24584
  TCP         [::1]:8080             [::1]:49695            ESTABLISHED     24584
  TCP         [::1]:8080             [::1]:50456            ESTABLISHED     24584
  TCP         [::1]:8080             [::1]:51161            ESTABLISHED     24584
  TCP         [::1]:8080             [::1]:63813            ESTABLISHED     24584
  TCP         [::1]:8080             [::1]:64467            ESTABLISHED     24584
  TCP         [::1]:8080             [::1]:64697            ESTABLISHED     24584
  TCP         [::1]:8080             [::1]:65060            ESTABLISHED     24584
  TCP         [::1]:49695            [::1]:8080             FIN_WAIT_1      44544
  TCP         [::1]:50456            [::1]:8080             FIN_WAIT_1      35136
  TCP         [::1]:51161            [::1]:8080             FIN_WAIT_1      42024
  TCP         [::1]:63813            [::1]:8080             FIN_WAIT_1      42024
  TCP         [::1]:64467            [::1]:8080             FIN_WAIT_1      42024
  TCP         [::1]:64697            [::1]:8080             FIN_WAIT_1      44544
  TCP         [::1]:65060            [::1]:8080             FIN_WAIT_1      35136
```

タスクマネージャーの詳細から、上記のPIDが何なのか特定

- IPv4（0.0.0.0:8080）→ PID：`25296` = com.docker.backend.exe
- IPv6（[::]:8080）→ PID：`24584` = wslrelay.exe（※WSLのネットワーク機能の一部）

よって原因は、  

1. Postman が IPv6 (::1) で 8080 にリクエスト

2. wslrelay.exe (PID: 24584) がパケットを受信

3. wslrelay.exe は「WSLの中で 8080 を待っている奴」に届けようとするが、Goアプリは Docker コンテナ（別のネットワーク空間）にいるため、WSLのシステム内では発見できない

4. 届け先が見つからない wslrelay.exe は応答を返せず、Postman はタイムアウトするまで待ち続ける

### 対処

Postmanで、IPv4 を指定する

```
http://localhost:8080/upload
            ↓
http://127.0.0.1:8080/upload
```


## 3. ホットリロードが効かない

### 原因

OS → Dockerコンテナ への変更通知がうまく伝わっていない

### 対処

```bash
poll = true           # 自分からファイルを確認しに行くモードをオンに変更
poll_interval = 500   # 0.5秒おきにファイルの更新をチェックするよう変更
```