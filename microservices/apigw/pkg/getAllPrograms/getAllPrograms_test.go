package getAllPrograms_test

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"testing"
	gg "webapi/microservices/apigw/pkg/getAllPrograms"
	gs "webapi/microservices/apigw/pkg/serverAliveConfirmer"
	os2 "webapi/pkg/os"
	"webapi/test"
)

var currentDir string
var addrs []string
var ports []string

func init() {
	c, err := os2.GetCurrentDir()
	if err != nil {
		log.Fatalln(err.Error())
	}
	currentDir = c

	// http://が必要かも
	addrs, ports, err = test.GetStartedServers(3)
	if err != nil {
		panic(err)
	}
}

func tearDown() {
	os.RemoveAll(filepath.Join(currentDir, "fileserver"))
	os.Remove(filepath.Join(currentDir, "log.txt"))
}

func TestGetAllPrograms(t *testing.T) {
	serverAliveConfirmer := gs.New()
	fmt.Println(addrs)
	aliveServers, err := gs.GetAliveServers(addrs, "/health", serverAliveConfirmer)
	if err != nil {
		t.Fatalf("err from GetAliveServers(): %v \n", err.Error())
	}
	allProgramGetter := gg.New()

	endPoint := "/program/all"
	allProgramMap, err := allProgramGetter.Get(aliveServers, endPoint)
	if err != nil {
		t.Fatalf("err from Get(): %v \n", err.Error())
	}

	if _, ok := allProgramMap["convertToJson"]; !ok {
		t.Fatalf("convertToJson is not found. of %v. \n", allProgramMap)
	}

	t.Cleanup(func() {
		tearDown()
	})

}
