/*
生きているサーバの全てのプログラム一覧を取得する。
{
	"programName" : {
		"command": "xxxxx",
		"help": "xxxxxxx"
	},
	......
}
*/

package getAllPrograms

import (
	"encoding/json"
	"fmt"
	"webapi/pkg/os"
)

type Getter interface {
	// Get は生きているサーバたちのリストを入れて、全てのサーバにアクセスし、
	// 全てのプログラム情報を取得しmapで返す。
	Get(aliveServers []string, endPoint string) (map[string]interface{}, error)
}

func New() Getter {
	return &getter{}
}

type getter struct{}

func (a *getter) Get(aliveServers []string, endPoint string) (map[string]interface{}, error) {
	allServerMaps := map[string]interface{}{}

	// aliveServersにアクセスしていき、プログラム情報を取得、
	// 全てのプログラム情報を取得し、allServerMapsに格納する。
	for _, server := range aliveServers {
		url := server + endPoint

		// サーバにアクセス
		stdout, _, err := os.SimpleExec(fmt.Sprintf("curl %v", url))
		if err != nil {
			return nil, fmt.Errorf("Get: %v", err)
		}

		// レスポンスをパース
		mapForResBody := map[string]interface{}{}
		err = json.Unmarshal([]byte(stdout), &mapForResBody)
		if err != nil {
			return nil, fmt.Errorf("Get: %v", err)
		}

		// mapForResBodyに一つのサーバから取り出したプログラム一覧(map[string]interface{}{})がある。
		for proName, proInfo := range mapForResBody {
			// allServerMapsにプログラムネームのキーが入っていなかったら追加する
			if _, ok := allServerMaps[proName]; !ok {
				m, _ := proInfo.(map[string]interface{})
				allServerMaps[proName] = m
			}
		}
	}

	return allServerMaps, nil
}
