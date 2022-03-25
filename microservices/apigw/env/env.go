/*
servers.jsonを読み込み、構造体に保持する。
*/

package env

import (
	"log"
	"os"

	pkgOs "webapi/pkg/os"
	pkgRandom "webapi/pkg/random"
)

type Env struct {
	// ExecServers eg -> ["http://192.168.59.101:30007,http://192.168.59.101:30008"]
	ExecServers []string
	// GateWayServerPort example -> "80"
	GateWayServerPort string
}

// New はservers.jsonの中身をserversConfig構造体にセットし、返す
func New() *Env {
	// 構造体を初期化
	e := &Env{}

	if os.Getenv("ENV") == "k8s" {
		// for k8s
		// k8sの場合はExecサービスがノードポートごとに異なる
		// かつノードの３台の中からランダムなIPにアクセスする必要がある。

		// ワーカーノード３台の中から一台を選択
		workerNodeIPs := pkgOs.ListEnvToSlice(os.Getenv("K8S_WORKER_NODE_IPS"))
		workerNodeIP := pkgRandom.Choice(workerNodeIPs)

		// それぞれ別のexecサーバにアクセスするようにURIを作成する。
		execNodePortPorts := pkgOs.ListEnvToSlice(os.Getenv("K8S_EXEC_NODEPORT_PORTS"))
		var execServers []string
		for _, p := range execNodePortPorts {
			uri := "http://" + workerNodeIP + ":" + p
			execServers = append(execServers, uri)
		}
		e.ExecServers = execServers

		e.GateWayServerPort = os.Getenv("K8S_APIGW_POD_PORT")

	} else {
		// for local
		m := make(map[string]string)

		// 実際に動作させる際は環境変数を設定しなければならない
		// 下記の設定はテスト用
		localIP := pkgOs.GetLocalIP()
		m["LOCAL_EXEC_SERVERS"] = "http://" + localIP + ":9001,http://" + localIP + ":9002,http://" + localIP + ":9003"

		for k, v := range m {
			err := pkgOs.SetEnvIfNotExists(k, v)
			if err != nil {
				log.Fatalf("New: %v", err.Error())
			}
		}
		e.ExecServers = pkgOs.ListEnvToSlice(os.Getenv("LOCAL_EXEC_SERVERS"))
		e.GateWayServerPort = os.Getenv("LOCAL_APIGW_PORT")
	}

	return e
}
