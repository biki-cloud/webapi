/*
servers.jsonを読み込み、構造体に保持する。
*/

package env

import (
	"fmt"
	"io"
	"log"
	"os"

	pkgOs "webapi/pkg/os"
)

type Env struct {
	// APIGateWayServers eg -> ["http://192.168.59.101:30007,http://192.168.59.101:30008"]
	APIGateWayServers []string
	MyPort            string
}

func Print(w io.Writer) {
	e := New()
	if os.Getenv("ENV") == "k8s" {
		fmt.Fprintf(w, "K8S_APIGW_NODEPORT_PORT : %v \n", e.APIGateWayServers)
	} else {
		fmt.Fprintf(w, "MY_PORT             : %v \n", e.MyPort)
		fmt.Fprintf(w, "LOCAL_APIGW_SERVERS : %v \n", e.APIGateWayServers)
	}
}

func New() *Env {
	// initialize Env struct
	e := &Env{}

	if os.Getenv("ENV") == "k8s" {
		// for k8s
		// 全てのワーカーノードのIPにapigwのポートを付与したURIのリストを作成。
		workerNodeIPs := pkgOs.ListEnvToSlice(os.Getenv("K8S_WORKER_NODE_IPS"))
		apigwNodePortPort := os.Getenv("K8S_APIGW_NODEPORT_PORT")
		for _, ip := range workerNodeIPs {
			apigwURI := "http://" + ip + ":" + apigwNodePortPort
			e.APIGateWayServers = append(e.APIGateWayServers, apigwURI)
		}

	} else {
		// ローカルでの動作の場合はこの環境変数を設定を使用する
		// 以下の設定は環境変数がセットされていなければセットし、
		// セットされていれば何もしない
		m := make(map[string]string)
		// 実際に動作させる際は環境変数を設定しなければならない
		localIP := pkgOs.GetLocalIP()
		m["LOCAL_APIGW_SERVERS"] = "http://" + localIP + ":8001,http://" + localIP + ":8002"

		for k, v := range m {
			err := pkgOs.SetEnvIfNotExists(k, v)
			if err != nil {
				log.Fatalf("New: %v", err.Error())
			}
		}
		e.APIGateWayServers = pkgOs.ListEnvToSlice(os.Getenv("LOCAL_APIGW_SERVERS"))
		e.MyPort = os.Getenv("MY_PORT")
	}

	return e
}
