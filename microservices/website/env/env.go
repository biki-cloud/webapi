/*
servers.jsonを読み込み、構造体に保持する。
*/

package env

import (
	"log"
	"os"
	pkgOs "webapi/pkg/os"
	"webapi/pkg/random"
)

type Env struct {
	APIGateWayServers []string
	// ExecServers eg -> ["http://192.168.59.101:30007,http://192.168.59.101:30008"]
	ExecServers []string
	MyPort      string
}

func New() *Env {
	// initialize Env struct
	e := &Env{}

	if os.Getenv("ENV") == "k8s" {
		// for k8s
		// k8sの場合はExecサービスがノードポートごとに異なる
		// かつノードの３台の中からランダムなIPにアクセスする必要がある。
		workerNodeIPs := pkgOs.ListEnvToSlice(os.Getenv("K8S_WORKER_NODE_IPS"))
		workerNodeIP := random.Choice(workerNodeIPs)

		// それぞれ別のexecサーバにアクセスするようにURIを作成する。
		execNodePortPorts := pkgOs.ListEnvToSlice(os.Getenv("K8S_EXEC_NODEPORT_PORTS"))
		var execServers []string
		for _, p := range execNodePortPorts {
			uri := "http://" + workerNodeIP + ":" + p
			execServers = append(execServers, uri)
		}
		e.ExecServers = execServers

		e.APIGateWayServers = []string{"http://" + workerNodeIP + ":" + os.Getenv("K8S_APIGW_NODEPORT_PORT")}

	} else {
		// for local
		// set environment variable here, if it is not setting
		m := make(map[string]string)
		// 実際に動作させる際は環境変数を設定しなければならない
		localIP := pkgOs.GetLocalIP()
		m["LOCAL_APIGW_SERVERS"] = "http://" + localIP + ":8001,http://" + localIP + ":8002"
		m["LOCAL_EXEC_SERVERS"] = "http://" + localIP + ":9001,http://" + localIP + ":9002"

		for k, v := range m {
			err := pkgOs.SetEnvIfNotExists(k, v)
			if err != nil {
				log.Fatalf("New: %v", err.Error())
			}
		}
		e.APIGateWayServers = pkgOs.ListEnvToSlice(os.Getenv("LOCAL_APIGW_SERVERS"))
		e.ExecServers = pkgOs.ListEnvToSlice(os.Getenv("LOCAL_EXEC_SERVERS"))
		e.MyPort = os.Getenv("MY_PORT")
	}

	return e
}
