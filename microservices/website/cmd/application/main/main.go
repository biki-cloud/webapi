package main

import (
	"flag"
	"os"
	"webapi/microservices/website/cmd/application"
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

	serverURI := ":" + myPort
	app := application.New()
	srv := application.NewServer(serverURI, app)

	err := srv.ListenAndServe()
	if err != nil {
		app.ErrorLog.Fatalln(err)
	}
}
