package selectServer_test

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
	"webapi/cli/pkg/selectServer"
	apigwApp "webapi/microservices/apigw/cmd/application"
	execApp "webapi/microservices/exec/cmd/application"
	pkgOs "webapi/pkg/os"
)

var (
	currentDir      string
	execServerPorts []string
	apigwServerPort string
)

func init() {
	c, err := pkgOs.GetCurrentDir()
	if err != nil {
		log.Fatalln(err.Error())
	}
	currentDir = c

	// 下の３つのポートはapigwの環境変数のLOCAL_EXEC_SERVERSの初期値に合わせなければならない
	// serve exec server
	execServerPorts = []string{"9001", "9002", "9003"}
	for _, p := range execServerPorts {
		var done chan error
		go func() {
			srv := execApp.NewServer(":"+p, execApp.New())
			if err := srv.ListenAndServe(); err != nil {
				panic(err)
			}
		}()
		select {
		case err := <-done:
			if err != nil {
				log.Fatalln(err.Error())
			}
		case <-time.After(1 * time.Second):
			fmt.Println("serve success")
		}
	}

	// APIGWサーバを起動
	apigwServerPort = "8001"
	var done chan error
	go func() {
		srv := apigwApp.NewServer(":"+apigwServerPort, apigwApp.New())
		done <- srv.ListenAndServe()
	}()
	select {
	case err := <-done:
		if err != nil {
			log.Fatalln(err.Error())
		}
	case <-time.After(1 * time.Second):
		fmt.Println("serve success")
	}
}

func tearDown() {
	os.RemoveAll(filepath.Join(currentDir, "fileserver"))
}

func contains(li []string, s string) bool {
	for _, a := range li {
		if strings.Contains(s, a) {
			return true
		}
	}
	return false
}

func TestSelectServer(t *testing.T) {
	selector := selectServer.New()

	loadBalancerURL := "http://127.0.0.1:" + apigwServerPort
	programName := "convertToJson"

	url := loadBalancerURL + "/program-server/minimumMemory-and-hasProgram/" + programName
	serverURL, err := selector.Select(url)
	if err != nil {
		t.Errorf("err from Select(): %v \n", err.Error())
	}

	if !contains(execServerPorts, serverURL) {
		t.Errorf("%v doesn't contain of %v \n", serverURL, execServerPorts)
	}

	t.Cleanup(func() {
		tearDown()
	})
}
