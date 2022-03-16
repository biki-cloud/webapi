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
	myPort string
)

func main() {
	// k8sの場合は80で固定にし、ローカルの場合は指定できるようにする
	// ローカルの場合のポートの指定はコマンドラインで行う
	flag.StringVar(&myPort, "port", "80", "server port")
	flag.Parse()
	os.Setenv("MY_PORT", myPort)

	cfg := env.New()
	a := application.New()

	// localでも実行はポートを指定できるが、k8sの場合は80で固定にする
	var serverURI string
	if os.Getenv("ENV") == "k8s" {
		serverURI = ":80"
	} else {
		serverURI = ":" + cfg.ProgramServerPort
	}

	srv := application.NewServer(serverURI, a)

	err := srv.ListenAndServe()
	if err != nil {
		log.Fatalln(err)
	}
}
