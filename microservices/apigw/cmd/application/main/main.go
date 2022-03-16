package main

import (
	"flag"
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
	flag.StringVar(&myPort, "port", "80", "server port")
	flag.Parse()
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
