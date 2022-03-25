package env_test

import (
	"reflect"
	"testing"

	"webapi/cli/env"
	pkgOs "webapi/pkg/os"
)

func TestNewEnv(t *testing.T) {
	err := pkgOs.SetEnvIfNotExists("APIGW_SERVER_URIS", "http://127.0.0.1:8001,http://127.0.0.1:8002,http://127.0.0.1:8003")
	if err != nil {
		t.Fatalf("Err: %v \n", err.Error())
	}

	e := env.New()

	servers := e.APIGWServerURIs
	wantServers := []string{
		"http://127.0.0.1:8001",
		"http://127.0.0.1:8002",
		"http://127.0.0.1:8003",
	}
	if !reflect.DeepEqual(servers, wantServers) {
		t.Errorf("e.APIGWServerURIs: %v, want: %v \n", servers, wantServers)
	}
}
