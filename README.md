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
1. どのようなマイクロサービスにするかを決定する。 <br>
実装内容: 入力ファイルに.jsonの拡張子を付与し、出力するマイクロサービスを作成する。

<br>

2. プロジェクトを作成
```shell
mkdir ConvertToJson
cd ConvertToJson
touch convert_to_json.py
```

<br>

3. コマンドを決定する
```shell
python3 convert_to_json.py <input file> <output dir> 
```

<br>

4. 実装する
```python
import os
import shutil
import sys

infile = sys.argv[1]
output_dir = sys.argv[2]

outfile = os.path.join(output_dir, os.path.basename(infile)) + ".json"
shutil.move(infile, outfile)
```

<br>

5. ヘルプを書く
```shell
cat help.txt
入力ファイルに.jsonを付与し、出力ディレクトリに移動させる。
```

## マイクロサービスをexecサービスへ登録する方法
1. exec/config/programConfig.jsonを編集する
```json
{
  "programs": [
    {
      "name": "ConvertToJson",
      "command": "python3 programs/ConvertToJson/convert_to_json.py INPUTFILE OUTPUTDIR",
      "helpPath": "programs/ConvertToJson/help.txt"
    }
  ]
}
```

<br>

2. [マイクロサービスの作成方法](#マイクロサービスの作成方法)で作成したプロジェクトをexecディレクトリへ格納する。
```shell
mv ConvertToJson exec/programs/ConvertToJson
```
