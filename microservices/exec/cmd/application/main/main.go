/*
プログラムサーバの開始
*/

package main

import (
	"flag"
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
