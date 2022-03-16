package application_test

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
	apigwApp "webapi/microservices/apigw/cmd/application"
	execApp "webapi/microservices/exec/cmd/application"
	pkgOs "webapi/pkg/os"
)

/*
ゲートウェイのハンドラーのテストはまずサーバを立てて、
そこから初めてハンドラーのテストをしなければならない。
またテストサーバのIPとconfのservers.jsonのサーバIPを同じにしなければならない。
*/

var (
	currentDir    string
	deletes       []string
	ts1, ts2, ts3 *httptest.Server
	err           error
	app           *apigwApp.Application
)

func init() {
	// サーバを立てるとカレントディレクトリにfileserverディレクトリとlog.txt
	// ができるのでそれを削除する。
	c, err := pkgOs.GetCurrentDir()
	if err != nil {
		log.Fatalln(err.Error())
	}
	currentDir = c
	serverSet()
}

func serverSet() {
	// apigwのenvの環境変数の初期値(LOCAL_EXEC_SERVERS)の初期値を見て下の
	// portsを設定する
	ports := []string{"9001", "9002", "9003"}
	for _, p := range ports {
		p := p
		go func() {
			a := execApp.New()
			srv := execApp.NewServer(":"+p, a)
			if err := srv.ListenAndServe(); err != nil {
				panic(err.Error())
			}
			time.Sleep(1 * time.Second)
		}()
	}

	app = apigwApp.New()
}

func tearDown() {
	os.RemoveAll(filepath.Join(currentDir, "fileserver"))
	os.Remove(filepath.Join(currentDir, "log.txt"))
}

func TestGetMinimumMemoryServerHandler(t *testing.T) {

	r, _ := http.NewRequest(http.MethodGet, "/program-server/memory/minimum", nil)
	w := httptest.NewRecorder()

	app.GetMinimumMemoryServerHandler(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("got %v, want %v, body: %v \n", w.Code, http.StatusOK, w.Body.String())
	}

	type j struct {
		Url string `json:"url"`
	}

	var d j
	err = json.Unmarshal(w.Body.Bytes(), &d)
	if err != nil {
		t.Errorf(err.Error())
	}

	if d.Url == "" {
		t.Errorf("d.Url is empty.")
	}

	t.Cleanup(func() {
		tearDown()
	})
}

func TestGetMinimumMemoryAndHasProgram(t *testing.T) {
	r, _ := http.NewRequest(http.MethodGet, "/program-server/minimumMemory-and-hasProgram/convertToJson", nil)
	w := httptest.NewRecorder()

	app.GetMinimumMemoryAndHasProgram(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("got %v, want %v", w.Code, http.StatusOK)
	}

	type j struct {
		Url string `json:"url"`
	}

	var d j
	err = json.Unmarshal(w.Body.Bytes(), &d)
	if err != nil {
		t.Errorf(err.Error() + w.Body.String())
	}

	if d.Url == "" {
		t.Errorf("d.Url is empty.")
	}

	t.Cleanup(func() {
		tearDown()
	})
}

func TestGetAliveServersHandler(t *testing.T) {

	r, _ := http.NewRequest(http.MethodGet, "/program-server/alive", nil)
	w := httptest.NewRecorder()

	app.GetAliveServersHandler(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("got %v, want %v", w.Code, http.StatusOK)
	}

	type data struct {
		AliveServers []string `json:"AliveServers"`
	}
	var d data
	err = json.Unmarshal(w.Body.Bytes(), &d)
	if err != nil {
		t.Errorf(err.Error())
	}

	if len(d.AliveServers) == 1 {
		t.Errorf("AliveServers is empty.")
	}

	t.Cleanup(func() {
		tearDown()
	})
}

func TestGetAllProgramsHandler(t *testing.T) {
	r, _ := http.NewRequest(http.MethodGet, "/program-server/program/all", nil)
	w := httptest.NewRecorder()

	app.GetAllProgramsHandler(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("got %v, want %v", w.Code, http.StatusOK)
	}

	b := w.Body.String()
	if !strings.Contains(b, "convertToJson") {
		t.Errorf("%v doesn't contain convertToJson", b)
	}

	t.Cleanup(func() {
		tearDown()
	})
}
