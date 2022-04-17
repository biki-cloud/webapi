package k8s

import (
	"errors"
	"fmt"
	"log"
	"os"
	"webapi/microservices/apigw/pkg/serverAliveConfirmer"
	pkgRandom "webapi/pkg/random"
)

var debugLog = log.New(os.Stdout, "DEBUG ", log.LstdFlags)

func LoadBalance(workerNodeIPs []string, port string) (string, error) {
	debugLog.Printf("k8s worker node ips: %v \n", workerNodeIPs)
	// 簡易的にk8sでのノードポートへのロードバランシング機能を実装する
	// クラスターを構築しているノードから生きているノードを選択する。
	var tmpNodeIP string
	var alivingNodeIP []string
	for len(workerNodeIPs) > 0 {
		tmpNodeIP, workerNodeIPs = workerNodeIPs[0], workerNodeIPs[1:]
		uri := "http://" + tmpNodeIP + ":" + port
		debugLog.Printf("Try to check k8s worker node is alive: %v \n", uri)
		alive, err := serverAliveConfirmer.New().IsAlive(uri, "/health")
		if err != nil {
			log.Fatalf("Server IsAlive: err: %v \n", err.Error())
		}
		if alive {
			debugLog.Printf("%v is alive. \n", uri)
			alivingNodeIP = append(alivingNodeIP, tmpNodeIP)
		} else {
			debugLog.Printf("%v is not alive. \n", uri)
		}
	}

	if len(alivingNodeIP) > 0 {
		ip := pkgRandom.Choice(alivingNodeIP)
		debugLog.Printf("Selected worker node ip: %v \n", ip)
		return ip, nil

	} else {
		err := errors.New("there is no alive workerNode. ")
		return "", fmt.Errorf("k8sLoadBalance: %v \n", err)
	}
}

//
//func GetRandomAliveNodeIP(workerNodeIPs []string) (string, error) {
//	debugLog.Printf("k8s workder node ips: %v \n", workerNodeIPs)
//	for i := 0; i < 100; i++ {
//		ip := pkgRandom.Choice(workerNodeIPs)
//		debugLog.Printf("ping to %v \n", ip)
//		if pkgOs.Ping(ip) {
//			debugLog.Printf("ping success %v \n", ip)
//			return ip, nil
//		} else {
//			debugLog.Printf("ping failed %v \n", ip)
//		}
//	}
//	err := errors.New("there is no alive workerNode. ")
//	return "", fmt.Errorf("k8s.GetRandomAliveNodeIP: %v \n", err)
//}
