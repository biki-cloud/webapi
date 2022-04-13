package main

import (
	"flag"
	"fmt"
	"os"

	"webapi/microservices/website/cmd/application"
	"webapi/microservices/website/env"
)

var (
	myPort string
)

func main() {
	// k8sの場合は80で固定にし、ローカルの場合は指定できるようにする
	// ローカルの場合のポートの指定はコマンドラインで行う
	flag.StringVar(&myPort, "port", "80", "サーバが起動するポートを指定")
	flag.CommandLine.Usage = func() {
		o := flag.CommandLine.Output()
		fmt.Fprintf(o, "\nUsage: \n  %s -port <サーバを起動させるポート番号> \n", flag.CommandLine.Name())
		fmt.Fprintf(o, "\n\n"+
			"Description:  \n  "+
			"Execサーバの登録プログラムを実行できるWEBページ\n\n  "+
			"実行する前にAPIGWサーバのアドレスを環境変数にセットしてください。値は環境に応じて変更してください。\n  "+
			"Linux  : export LOCAL_APIGW_SERVERS=http://127.0.0.1:8001,http://127.0.0.1:8002,http://127.0.0.1:8003 \n  "+
			"Windows: SET LOCAL_APIGW_SERVERS=http://127.0.0.1:8001,http://127.0.0.1:8002,http://127.0.0.1:8003 \n\n  "+
			"実行する前にExecサーバのアドレスを環境変数にセットしてください。値は環境に応じて変更してください。\n  "+
			"Linux  : export LOCAL_EXEC_SERVERS=http://127.0.0.1:9001,http://127.0.0.1:9002,http://127.0.0.1:9003 \n  "+"Windows: SET LOCAL_EXEC_SERVERS=http://127.0.0.1:9001,http://127.0.0.1:9002,http://127.0.0.1:9003 \n\n  "+
			"\nOptions: \n")
		flag.PrintDefaults()
		fmt.Fprintf(o, "\nUpdated date 2022.4.6 by morituka. \n\n")
	}
	flag.Parse()

	// 引数がなければヘルプを表示する
	if len(os.Args) == 1 {
		flag.CommandLine.Usage()
		os.Exit(1)
	}

	flag.Parse()
	os.Setenv("MY_PORT", myPort)

	// 設定値を表示させる
	env.Print(os.Stdout)

	serverURI := ":" + myPort
	app := application.New()
	srv := application.NewServer(serverURI, app)

	err := srv.ListenAndServe()
	if err != nil {
		app.ErrorLog.Fatalln(err)
	}
}
