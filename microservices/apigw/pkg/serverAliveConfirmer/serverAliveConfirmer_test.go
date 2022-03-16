package serverAliveConfirmer_test

import (
	"log"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"webapi/microservices/apigw/pkg/serverAliveConfirmer"
	os2 "webapi/pkg/os"
	"webapi/test"
)

var currentDir string
var confirmer serverAliveConfirmer.Confirmer

func init() {
	c, err := os2.GetCurrentDir()
	if err != nil {
		log.Fatalln(err.Error())
	}
	currentDir = c
	confirmer = serverAliveConfirmer.New()
}

func tearDown() {
	os.RemoveAll(filepath.Join(currentDir, "fileserver"))
	os.RemoveAll(filepath.Join(currentDir, "log.txt"))
}

func TestIsAlive(t *testing.T) {
	addrs, _, err := test.GetStartedServers(1)
	if err != nil {
		log.Fatalln(err)
	}

	t.Run("exec is alive.", func(t *testing.T) {
		testIsAlive(t, addrs[0], "/health", true)
	})

	t.Run("exec is not alive.", func(t *testing.T) {
		testIsAlive(t, "http://127.0.0.1:8052", "/user/top", false)
	})

	t.Cleanup(func() {
		tearDown()
	})
}

func testIsAlive(t *testing.T, addr, endPoint string, expect bool) {
	t.Helper()
	alive, err := confirmer.IsAlive(addr, endPoint)
	if err != nil {
		t.Errorf("err occur: %v \n", err.Error())
	}
	if alive != expect {
		t.Errorf("got: %v, want false.", expect)
	}
}

func TestGetAliveServers(t *testing.T) {
	addrs, _, err := test.GetStartedServers(3)
	if err != nil {
		log.Fatalln(err)
	}

	t.Run("get alive servers 1", func(t *testing.T) {
		testGetAliveServers(t, addrs, "/health", addrs)
	})

	t.Cleanup(func() {
		tearDown()
	})
}

func testGetAliveServers(t *testing.T, servers []string, endPoint string, expectServers []string) {
	t.Helper()
	aliveServers, err := serverAliveConfirmer.GetAliveServers(servers, endPoint, confirmer)
	if err != nil {
		t.Errorf(err.Error())
	}

	if !reflect.DeepEqual(aliveServers, expectServers) {
		t.Errorf("got: %v, want: %v \n", aliveServers, expectServers)
	}
}
