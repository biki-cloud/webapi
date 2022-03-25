/*
プログラムサーバの開始
*/

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"webapi/microservices/exec/cmd/application"
	"webapi/microservices/exec/env"
)

var (
	// サーバが起動するポート
	myPort string
)

func main() {
	// k8sの場合は80で固定にし、ローカルの場合は指定できるようにする
	// ローカルは通常の場合コマンドラインでポートの指定を行う
	flag.StringVar(&myPort, "port", "80", "server port")
	flag.CommandLine.Usage = func() {
		o := flag.CommandLine.Output()
		fmt.Fprintf(o, "\nUsage: \n  %s -port <サーバを起動させるポート番号> \n", flag.CommandLine.Name())
		fmt.Fprintf(o, "\n\n"+
			"Description:  \n  "+
			"登録プログラムを実行させるサーバ\n\n  "+
			"以下の環境変数を変更することでサーバの設定項目を変更することができます。 \n  "+
			"WORKED_DIR_KEEP_SEC(プログラム実行が終了したディレクトリを何秒保持するか) default = 600  \n  "+
			"EXECUTE_TIMEOUT_SEC(プログラム実行が始まって何秒経過したらタイムアウトとするか) default = 100 \n  "+
			"STDOUT_BUFFER_SIZE(プログラムの標準出力を何バイト出力するか) default = 1000000 \n  "+
			"STDERR_BUFFER_SIZE(プログラムの標準エラー出力を何バイト出力するか) default = 1000000 \n  "+
			"MAX_UPLOAD_SIZE_MB(入力ファイルサイズの上限を何メガにするか) default = 300 \n  "+
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

	flag.Parse()
	os.Setenv("MY_PORT", myPort)

	e := env.New()
	a := application.New()

	var serverURI string
	if os.Getenv("ENV") == "k8s" {
		serverURI = ":80"
	} else {
		serverURI = ":" + e.ProgramServerPort
	}

	srv := application.NewServer(serverURI, a)

	err := srv.ListenAndServe()
	if err != nil {
		log.Fatalln(err)
	}
}
