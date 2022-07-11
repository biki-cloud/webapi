package application

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"sort"
	pkgOs "webapi/pkg/os"

	"webapi/microservices/apigw/pkg/memoryGetter"
	"webapi/microservices/apigw/pkg/minimumServerSelector"
	"webapi/microservices/apigw/pkg/programHasServers"
	"webapi/microservices/apigw/pkg/serverAliveConfirmer"
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
	// APIゲートウェイサーバから生きているサーバを取得、かつメモリ消費が一番少ないサーバを選択する。
	confirmer := serverAliveConfirmer.New()
	apigwServers, err := serverAliveConfirmer.GetAliveServers(app.Env.APIGateWayServers, "/health", confirmer)
	if len(apigwServers) == 0 {
		app.ServerError(w, fmt.Errorf("生きているAPIゲートウェイサーバはありませんでした。"))
	}
	if err != nil {
		app.ServerError(w, err)
	}

	mg := memoryGetter.New()
	// 生きているサーバにアクセスしていき、メモリ状況を取得、一番消費メモリが少ないサーバを取得する
	serverMemoryMap, err := minimumServerSelector.GetServerMemoryMap(apigwServers, "/health/memory", mg)
	if err != nil {
		log.Fatalf("APIゲートウェイサーバにてエラーが発生しました。err: %v\n", err)
	}

	minUrl := minimumServerSelector.GetMinimumMemoryServer(serverMemoryMap)

	apigwAddr := minUrl

	// 登録プログラムの情報を取得する
	command := fmt.Sprintf("curl %v/program-server/program/all", apigwAddr)
	programsJSON, stderr, err := pkgOs.SimpleExec(command)
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
		Help    template.HTML `json:"help"`
		Command string `json:"command"`
	}
	type proInfo struct {
		Name    string
		Help    template.HTML
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

		pg := programHasServers.New()
		phs, err := pg.Get(t.AliveExecServers, "/program/all", proName)
		if err != nil {
			app.ServerError(w, err)
		}

		// プログラムを保持しているサーバたちの中で一番使用メモリが少ないサーバを選択する。
		// この処理を行うと、execサーバが多くなってくるとwebsiteの/user/topにアクセスした時にローディング時間が長くなる。
		// なので対策としてはプログラムを保持しているサーバをメモリの使用量で判断するのではなく、ランダムで一台選択するようにする。
		mmss := minimumServerSelector.New()
		mg := memoryGetter.New()
		sac := serverAliveConfirmer.New()
		url, err := mmss.Select(phs, sac, mg, "/health/memory", "/health")
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
	currentDir, err := pkgOs.GetCurrentDir()
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

	currentDir, err := pkgOs.GetCurrentDir()
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
	help := r.FormValue("help")

	type data struct {
		Name          string
		Help          string
		ExecServerURI string
	}

	d := data{Name: proName, ExecServerURI: execServerURI, Help: help}

	app.RenderTemplate(w, serveHtml, d)
}
