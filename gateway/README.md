
## クライアントとサービス間の通信

クライアントとバックエンドサービスの通信は全てAPI Gatewayを経由する。
API Gatewayを経由することでクライアントとバックエンドサービスを疎結合とする。
スマートフォンやPCなどプラットフォーム間で要求するリクエストとレスポンスに差異が生じる場合は[BFF](https://learn.microsoft.com/en-us/azure/architecture/patterns/backends-for-frontends)の採用を検討する。

## サービス間通信の方式

サービス間の通信は可能な限りメッセージング(RabbitMQ)は採用する。
マイクロサービス間の通信が必要な場合、非同期通信が可能であれば極力非同期通信とする。

### rpc(gRPC)を採用する場合の懸念

https://microservices.io/patterns/communication-style/rpi.html

* 可用性の低下
* サービス間の密結合

https://particular.net/blog/rpc-vs-messaging-which-is-faster

> But usually, the async processing model results in a more parallel-processing result, 
> resulting in higher throughput for your system than the synchronous blocking RPC model.

訳文

> ただし、通常、非同期処理モデルはより並列処理の結果をもたらし、
> 同期ブロッキング RPC モデルよりもシステムのスループットが高くなります。

* トラフィックの増大によりシステムのスループットの低下が懸念される

## 開発用メモ

### 環境構築手順

1. Go: Install/Update Toolsを実行する

### gRPCのセットアップ

```
# install gRPC tools
sudo apt upate
sudo apt install -y protobuf-compiler
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.27
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
```

```
# generate protobuf files
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative logs.proto
```