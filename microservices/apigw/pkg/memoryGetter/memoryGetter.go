/*
サーバにアクセスし、ランタイム構造体を取得する機能を提供するパッケージ
*/

package memoryGetter

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"runtime"
)

type Getter interface {
	Get(string) (runtime.MemStats, error)
}

func New() Getter {
	return &getter{}
}

type getter struct{}

// Get はメモリ状況のjsonAPIを公開しているサーバのURLを受け取り、
// get(http)し、そのjsonをruntime.MemStatsにデコードし、runtime.MemStatsを返す
// eg url: "http://127.0.0.1:8093/health/memory"
func (g *getter) Get(url string) (runtime.MemStats, error) {
	resp, err := http.Get(url)

	if resp != nil {
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				log.Fatalln(err.Error())
			}
		}(resp.Body)
	}

	if err != nil {
		return runtime.MemStats{}, err
	}

	body, _ := ioutil.ReadAll(resp.Body)

	var d runtime.MemStats
	err = json.Unmarshal(body, &d)
	if err != nil {
		return runtime.MemStats{}, fmt.Errorf("Get: %v", err)
	}

	return d, nil
}
