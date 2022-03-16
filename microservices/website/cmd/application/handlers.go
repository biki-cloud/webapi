package application

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"sort"
	mg "webapi/microservices/apigw/pkg/memoryGetter"
	"webapi/microservices/apigw/pkg/minimumServerSelector"
	gp "webapi/microservices/apigw/pkg/programHasServers"
	"webapi/microservices/apigw/pkg/serverAliveConfirmer"
	"webapi/pkg/os"
)

// mapToStruct はmapからstructに変換する。
func mapToStruct(m interface{}, val interface{}) error {
	tmp, err := json.Marshal(m)
	if err != nil {
		return fmt.Errorf("mapToStruct: %v", err)
	}
	err = json.Unmarshal(tmp, val)
	if err != nil {
		return fmt.Errorf("mapToStruct: %v", err)
	}
	return nil
}

// Top はユーザがwebにアクセスした場合はに
// 全サーバにアクセスし、全てのプログラムリストを取得,
// プログラムがあるサーバなおかつメモリ使用量が最小のサーバのURLをセットし、
// webpageに表示する.ボタンが押されたらそのサーバのプログラム実行準備画面に飛ぶ
func (app *Application) Top(w http.ResponseWriter, r *http.Request) {
	// config.jsonのAPIゲートウェイサーバからメモリ消費が一番少ないサーバを選択する。
	// 生きているサーバのリストを取得
	confirmer := serverAliveConfirmer.New()
	apigwServers, err := serverAliveConfirmer.GetAliveServers(app.Env.APIGateWayServers, "/health", confirmer)
	if len(apigwServers) == 0 {
		app.ServerError(w, fmt.Errorf("生きているAPIゲートウェイサーバはありませんでした。"))
	}
	if err != nil {
		app.ServerError(w, err)
	}

	memoryGetter := mg.New()
	// 生きているサーバにアクセスしていき、メモリ状況を取得、一番消費メモリが少ないサーバを取得する
	serverMemoryMap, err := minimumServerSelector.GetServerMemoryMap(apigwServers, "/health/memory", memoryGetter)
	if err != nil {
		log.Fatalf("APIゲートウェイサーバにてエラーが発生しました。err: %v\n", err)
	}

	minUrl := minimumServerSelector.GetMinimumMemoryServer(serverMemoryMap)

	apigwAddr := minUrl

	// 登録プログラムの情報を取得する
	command := fmt.Sprintf("curl %v/program-server/program/all", apigwAddr)
	programsJSON, stderr, err := os.SimpleExec(command)
	if err != nil || programsJSON == "" {
		app.ServerError(w, fmt.Errorf("err: %v, stderr: %v", err.Error(), stderr))
	}

	// mapに変換
	var programsMap map[string]interface{}
	err = json.Unmarshal([]byte(programsJSON), &programsMap)
	if err != nil {
		app.ServerError(w, err)
	}

	type tmpProInfo struct {
		Help    string `json:"help"`
		Command string `json:"command"`
	}
	type proInfo struct {
		Name    string
		Help    string
		Command string
		URL     string
	}
	// proInfoのリストをhtmlに与えるだけで良いが今後さらに項目を渡す場合に備えて構造体を定義しておく
	type htmlData struct {
		ProInfos []proInfo
	}

	var dataToHtml htmlData

	// execサーバのリストを取得する
	resp, err := http.Get(apigwAddr + "/program-server/alive")
	if err != nil {
		app.ServerError(w, err)
	}

	type tmp struct {
		AliveExecServers []string `json:"AliveServers"`
	}

	var t tmp
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		app.ServerError(w, err)
	}
	err = json.Unmarshal(body, &t)
	if err != nil {
		app.ServerError(w, err)
	}

	// 生きているサーバたちにプログラムを持っているか聞いていき、持っていた何台かがあった場合
	// 一番消費メモリの少ないサーバを選択する。
	// htmlに渡すhtmlData(struct)をこのループで完成させる。
	for proName, programInfoMap := range programsMap {
		var tmpInfo tmpProInfo
		err := mapToStruct(programInfoMap, &tmpInfo)
		if err != nil {
			app.ServerError(w, err)
		}

		programHasServersGetter := gp.New()
		programHasServers, err := programHasServersGetter.Get(t.AliveExecServers, "/program/all", proName)
		if err != nil {
			app.ServerError(w, err)
		}

		// プログラムを保持しているサーバたちの中で一番使用メモリが少ないサーバを選択する。
		minimumMemoryServerSelector := minimumServerSelector.New()
		memoryGetter := mg.New()
		sc := serverAliveConfirmer.New()
		url, err := minimumMemoryServerSelector.Select(programHasServers, sc, memoryGetter, "/health/memory", "/health")
		if err != nil {
			app.ServerError(w, err)
		}

		p := proInfo{
			Name:    proName,
			Help:    tmpInfo.Help,
			Command: tmpInfo.Command,
			URL:     url,
		}
		dataToHtml.ProInfos = append(dataToHtml.ProInfos, p)
	}

	w.Header().Add("Content-Type", "text/html")
	currentDir, err := os.GetCurrentDir()
	serveHtml := filepath.Join(currentDir, "ui/html", "top.html")
	absHtml, err := filepath.Abs(serveHtml)
	if err != nil {
		app.ServerError(w, err)
	}

	// 名前順でソートする
	sort.Slice(dataToHtml.ProInfos, func(i, j int) bool { return dataToHtml.ProInfos[i].Name < dataToHtml.ProInfos[j].Name })

	app.RenderTemplate(w, absHtml, dataToHtml)
}

func (app *Application) Exec(w http.ResponseWriter, r *http.Request) {

	currentDir, err := os.GetCurrentDir()
	if err != nil {
		app.ServerError(w, err)
	}
	serveHtml := filepath.Join(currentDir, "ui/html", "exec.html")

	proName := r.URL.Path[len("/user/exec/"):]
	r.ParseForm()

	// top.htmlを表示する段階でプログラムの実行サーバのURIは決まっていて
	// それをURIをformの一部に載せるのでそれをここで受け取り、またexec.htmlに渡して、
	// ajaxでそのURIのサーバへPOSTする
	execServerURI := r.FormValue("execServerURI")

	type data struct {
		Name          string
		ExecServerURI string
	}

	d := data{Name: proName, ExecServerURI: execServerURI}

	app.RenderTemplate(w, serveHtml, d)
}
