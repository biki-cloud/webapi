package minimumServerSelector_test

import (
	"log"
	"os"
	"path/filepath"
	"testing"
	"webapi/microservices/apigw/env"
	mg "webapi/microservices/apigw/pkg/memoryGetter"
	"webapi/microservices/apigw/pkg/minimumServerSelector"
	sc "webapi/microservices/apigw/pkg/serverAliveConfirmer"
	os2 "webapi/pkg/os"
	"webapi/test"
)

var (
	memoryGetter                mg.Getter
	currentDir                  string
	minimumMemoryServerSelector = minimumServerSelector.New()
	serverAliveConfirmer        = sc.New()
	// CreateDummyServer関数のcleanUpで削除するファイルたち
	deletes []string
	cfg     = env.New()
	addrs   []string
)

func init() {
	// サーバを立てるとカレントディレクトリにfileserverディレクトリとlog.txt
	// ができるのでそれを削除する。
	c, err := os2.GetCurrentDir()
	if err != nil {
		log.Fatalln(err.Error())
	}
	currentDir = c
	memoryGetter = mg.New()
	addrs, _, err = test.GetStartedServers(3)
	if err != nil {
		log.Fatalln(err)
	}
}

func tearDown() {
	os.RemoveAll(filepath.Join(currentDir, "fileserver"))
	os.RemoveAll(filepath.Join(currentDir, "log.txt"))
}

func TestSelect(t *testing.T) {

	selectedServer, err := minimumMemoryServerSelector.Select(addrs, serverAliveConfirmer, memoryGetter, "/health/memory", "/health")
	if err != nil {
		t.Errorf("err occur: %v \n", err.Error())
	}

	contain := false
	for _, addr := range addrs {
		if selectedServer == addr {
			contain = true
		}
	}
	if !contain {
		t.Errorf("selected exec: %v doesn't contain servers: %v", selectedServer, addrs)
	}

}

func TestGetMinimumMemoryServer(t *testing.T) {
	serverInfoMap := map[string]uint64{
		"http://127.0.0.1:8083": 3000,
		"http://127.0.0.1:8081": 2597,
		"http://127.0.0.1:8082": 2700,
	}

	wantAddr := "http://127.0.0.1:8081"
	addr := minimumServerSelector.GetMinimumMemoryServer(serverInfoMap)
	if addr != wantAddr {
		t.Errorf("GetMinimumMemoryServer(): %v, want: %v \n", addr, wantAddr)
	}

	t.Cleanup(func() {
		tearDown()
	})
}

func TestGetServerMemoryMap(t *testing.T) {
	serverInfoMap, err := minimumServerSelector.GetServerMemoryMap(addrs, "/health/memory", memoryGetter)
	if err != nil {
		t.Errorf(err.Error())
	}

	if len(serverInfoMap) != 3 {
		t.Errorf("len(serverInfoMap) is not 3. got: %v \n", len(serverInfoMap))
	}

	for _, memory := range serverInfoMap {
		if memory < 1 {
			t.Errorf("memory(%v) is not more than 1. \n", memory)
		}
	}
	t.Cleanup(func() {
		tearDown()
	})
}
