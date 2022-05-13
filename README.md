## Easy Micro Service

誰でも簡単にマイクロサービスを開発できるプラットフォーム。
K8Sの上で動作する。

## Contents
- MicroServices
    - [website](#website)
    - [apigw](#apigw)
    - [exec](#exec)
    - [cli](#cli)
- [Test](#Test)


## website
作成したマイクロサービスをブラウザで使用する

## apigw
execへのロードバランサを担当している

## exec
作成したマイクロサービスを登録し、マイクロサービスへリクエストが来れば、実行する。

## cli
コマンドラインでexecに登録したマイクロサービスを利用できる


## Test
```shell
go test ./...
```
