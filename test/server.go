package test

import (
	"fmt"
	"time"

	execApp "webapi/microservices/exec/cmd/application"
	pkgHttpURL "webapi/pkg/http/url"
)

// GetStartedServers
// 使用されていないポートからなるhttpから始まるアドレス、ポートのリストを返す。
// eg.
// addrs: ["http://127.0.0.1:8001", "http://127.0.0.1:8003", "http://127.0.0.1:8003"]
// ports: ["8001", "8002", "8003"]
func GetStartedServers(numberOfServer int) (addrs, ports []string, err error) {
	for (len(addrs) < numberOfServer && len(ports) < numberOfServer) || (addrs == nil && ports == nil) {
		addr, err := pkgHttpURL.GetLoopBackURL()
		if err != nil {
			return nil, nil, err
		}
		port := pkgHttpURL.GetPortFromURL(addr)

		var done chan error
		go func() {
			// サーバを起動するときはhttp://をつけてはいけない
			a := execApp.New()
			srv := execApp.NewServer(addr, a)
			done <- srv.ListenAndServe()
		}()
		// １秒かかる前にserverStartに値が入ってきたということはhttp.ListenAndServeがエラーですぐ終了した場合。
		// １秒かかったということはhttp.ListenAndServeに成功したということ。
		select {
		// GETなどをするときはhttp://をつけなければならない
		case err := <-done:
			if err == nil {
				addrs = append(addrs, "http://"+addr)
				ports = append(ports, port)
			}
		case <-time.After(1 * time.Second):
			addrs = append(addrs, "http://"+addr)
			ports = append(ports, port)
		}
	}

	fmt.Printf("addrs: %v, ports: %v from GetStartedServers \n", addrs, ports)
	return
}
