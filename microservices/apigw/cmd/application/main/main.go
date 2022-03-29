package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"webapi/microservices/apigw/cmd/application"
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
			"Execサーバと通信し、メモリ量、登録プログラムなどを把握し、管理するサーバ\n\n  "+
			"実行する前にExecサーバのアドレスを環境変数にセットしてください。値は環境に応じて変更してください。\n  "+
			"Linux  : export LOCAL_EXEC_SERVERS=http://127.0.0.1:9001,http://127.0.0.1:9002,http://127.0.0.1:9003 \n  "+
			"Windows: SET LOCAL_EXEC_SERVERS=http://127.0.0.1:9001,http://127.0.0.1:9002,http://127.0.0.1:9003 \n\n  "+
			" \n\nOptions: \n")
		flag.PrintDefaults()
		fmt.Fprintf(o, "\nUpdated date 2022.3.25 by morituka. \n\n")
	}
	flag.Parse()

	// 引数がなければヘルプを表示する
	if len(os.Args) == 1 {
		flag.CommandLine.Usage()
		os.Exit(1)
	}

	os.Setenv("LOCAL_APIGW_PORT", myPort)

	//cfg := env.New()
	a := application.New()
	serverURI := ":" + myPort
	srv := application.NewServer(serverURI, a)

	err := srv.ListenAndServe()
	if err != nil {
		log.Fatalln(err)
	}
}
