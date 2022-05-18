## EMS (Easy Micro Service)

EMS is a platform where everybody can develop microservices.
Users can concentrate on developing their microservices.
Usually, EMS works on the k8s.

## Contents
- MicroServices
    - [website](#website)
    - [apigw](#apigw)
      - [apigw REST API](#apigw-REST-API)
    - [exec](#exec)
      - [exec REST API](#exec-REST-API)
    - [cli](#cli)
- [Test](#Test)
- [How to create microservices](#How-to-create-microservices)


## website
website is a web application that users use registered microservices of EMS.

## apigw
apigw works load balancing to exec service and manage registered microservices of exec service.

### apigw REST API
```go
package application

import (
  "net/http"

  pkgHttpHandlers "webapi/pkg/http/handlers"
)

func (app *Application) Routes() *http.ServeMux {
  router := http.NewServeMux()

  // コマンドラインからはここにアクセスし、メモリ使用量が一番低いサーバのURLを返す。
  router.HandleFunc("/program-server/memory/minimum", app.GetMinimumMemoryServerHandler)

  // コマンドラインからここにアクセスし、プログラムがあるかつメモリ使用量が一番低いサーバのURLを返す。
  router.HandleFunc("/program-server/minimumMemory-and-hasProgram/", app.GetMinimumMemoryAndHasProgram)

  // 現在稼働しているサーバを返すAPI
  router.HandleFunc("/program-server/alive", app.GetAliveServersHandler)

  // 生きている全てのサーバのプログラムを取得してJSONで表示するAPI
  router.HandleFunc("/program-server/program/all", app.GetAllProgramsHandler)

  // このサーバが生きているかを判断するのに使用するハンドラ
  router.HandleFunc("/health", pkgHttpHandlers.HealthHandler)

  // このサーバプログラムのメモリ状態をJSONで表示するAPI
  router.HandleFunc("/health/memory", pkgHttpHandlers.GetRuntimeHandler)

  return router
}
```

## exec
exec works registering created microservices and executing microservices.

### exec REST API
```go
package application

import (
  "net/http"

  pkgHttpHandlers "webapi/pkg/http/handlers"
)

// Routes ハンドラをセットしたrouterを返す。
func (app *Application) Routes() *http.ServeMux {
  r := http.NewServeMux()

  // ファイルサーバーの機能のハンドラ
  // Env.FileServer.Dir以下のファイルをwebから見ることができる。
  fileServer := "/" + app.Cfg.FileServer.Dir + "/"
  r.Handle(fileServer, http.StripPrefix(fileServer, http.FileServer(http.Dir(app.Cfg.FileServer.Dir))))

  // 登録プログラムを実行させるAPI
  r.HandleFunc("/api/exec/", app.APIExec)

  // ファイルをアップロードするAPI
  r.HandleFunc("/api/upload", app.APIUpload)

  // このサーバプログラムのメモリ状態をJSONで表示するAPI
  r.HandleFunc("/health/memory", pkgHttpHandlers.GetRuntimeHandler)

  // プログラムサーバに登録してあるプログラム一覧をJSONで表示するAPI
  r.HandleFunc("/program/all", app.AllHandler)

  // このサーバが生きているかを判断するのに使用するハンドラ
  r.HandleFunc("/health", pkgHttpHandlers.HealthHandler)

  // コンテンツをダウンロードするためのAPI
  r.HandleFunc("/download/", app.Download)

  return r
}
```

## cli
コマンドラインでexecに登録したマイクロサービスを利用できる

### How to use
```shell
# 一番シンプルな実行方法 
cli -name <プロラム名> -i <入力ファイル> -o <出力ファイル> 
   
# パラメータを付加させる場合, -pの後の文字列をダブルクォーテーションで囲む必要がある。中の文字列の構成は登録プログラムの仕様に依存する。 
cli -name <プロラム名> -i <入力ファイル> -o <出力ファイル> -p "<パラメータ１,パラメータ２>" 
   
# 実行結果をJSONで受け取る場合 
cli -j -name <プログラム名> -i <入力ファイル> -o <出力ファイル> 
 
# プログラムの処理過程を表示しながら実行する場合 
cli -l -name <プログラム名> -i <入力ファイル> -o <出力ファイル>
```


## Test
```shell
go test ./...
```

## How to create microservices
1. Determine microservices that we want to create <br>
Content of implement: We want to create microservices that take file, then output file that is added extension of ".json" to them.

<br>

2. Create project
```shell
mkdir ConvertToJson
cd ConvertToJson
touch convert_to_json.py
```

<br>

3. Determine command to execute
```shell
python3 convert_to_json.py <input file> <output dir> 
```

<br>

4. Coding
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

5. Write help about coded program
```shell
cat help.txt
take file, then output file that is added extension of ".json" to them.
```

## How to register microservices to exec service.
1. Edit exec/config/programConfig.json
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

2. [How to create microservices](#how-to-create-microservices)
```shell
mv ConvertToJson exec/programs/ConvertToJson
```
