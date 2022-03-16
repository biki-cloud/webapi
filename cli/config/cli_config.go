/*
config.jsonを読み込み、中身を保持する機能を提供するパッケージ
*/

package config

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"webapi/pkg/os"
)

var serversConfigPath string

func SetConfPath(confName string) {
	currentDir, err := os.GetCurrentDir()
	if err != nil {
		panic("err msg: " + err.Error())
	}
	serversConfigPath = filepath.Join(currentDir, confName)
}

type config struct {
	APIGateWayServers []string `json:"APIGateWayServers"`
}

// New はservers.jsonの中身をserversConfig構造体にセットし、返す
func New() *config {
	if serversConfigPath == "" {
		SetConfPath("cli_config.json")
	}
	// 構造体を初期化
	conf := &config{}

	// 設定ファイルを読み込む
	cValue, err := ioutil.ReadFile(serversConfigPath)
	if err != nil {
		panic(err.Error())
	}

	// 読み込んだjson文字列をデコードし構造体にマッピング
	err = json.Unmarshal([]byte(cValue), conf)
	if err != nil {
		panic(err.Error())
	}

	return conf
}
