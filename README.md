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
- [マイクロサービスの作成方法](#マイクロサービスの作成方法)


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

## マイクロサービスの作成方法
1. どのようなマイクロサービスにするかを決定する。
入力ファイルに.jsonの拡張子を付与し、出力するマイクロサービスを作成する。

2. コマンドを決定する
```shell
python3 convert_to_json.py <input file> <output dir> 
```

3. 実装する
```python
import os
import shutil
import sys

infile = sys.argv[1]
output_dir = sys.argv[2]

outfile = os.path.join(output_dir, os.path.basename(infile)) + ".json"
shutil.move(infile, outfile)
```
