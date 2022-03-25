/*
config.jsonを読み込み、中身を保持する機能を提供するパッケージ
*/

package env

import (
	"os"
	pkgOs "webapi/pkg/os"
)

type Env struct {
	// APIGWServerURIs eg -> ["http://127.0.0.1:8001","http://127.0.0.1:8002"]
	APIGWServerURIs []string
}

// New はservers.jsonの中身をserversConfig構造体にセットし、返す
func New() *Env {
	e := &Env{}
	e.APIGWServerURIs = pkgOs.ListEnvToSlice(os.Getenv("APIGW_SERVER_URIS"))
	return e
}
