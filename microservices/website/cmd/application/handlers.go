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

	minURL := minimumServerSelector.GetMinimumMemoryServer(serverMemoryMap)

	apigwAddr := minURL

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
		// Help はプログラムの場所においてあるhelp.txtなどを読み込み、/user/topに表示するための値。
		// その時、プログラム開発者次第でわかりやすくヘルプが作成できるようにhelp.txtなどに書かれた
		// 文字列をHTMLに変換している。なお、help.txtに普通の文字列を記入した場合は普通に表示されるため
		// 余裕がある人だけHTMLで記述することができる。
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

	var dataToHTML htmlData

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
		dataToHTML.ProInfos = append(dataToHTML.ProInfos, p)
	}

	w.Header().Add("Content-Type", "text/html")
	currentDir, err := pkgOs.GetCurrentDir()
	serveHTML := filepath.Join(currentDir, "ui/html", "top.html")
	absHTML, err := filepath.Abs(serveHTML)
	if err != nil {
		app.ServerError(w, err)
	}

	// 名前順でソートする
	sort.Slice(dataToHTML.ProInfos, func(i, j int) bool { return dataToHTML.ProInfos[i].Name < dataToHTML.ProInfos[j].Name })

	app.RenderTemplate(w, absHTML, dataToHTML)
}

// Exec ファイルなどをドラッグし、実行ボタンを押してリクエストをするためのページ
func (app *Application) Exec(w http.ResponseWriter, r *http.Request) {

	currentDir, err := pkgOs.GetCurrentDir()
	if err != nil {
		app.ServerError(w, err)
	}
	serveHTML := filepath.Join(currentDir, "ui/html", "exec.html")

	r.ParseForm()

	// top.htmlを表示する段階でプログラムの実行サーバのURIは決まっていて
	// それをURIをformの一部に載せるのでそれをここで受け取り、またexec.htmlに渡して、
	// ajaxでそのURIのサーバへPOSTする
	execServerURI := r.FormValue("execServerURI")
	proName := r.FormValue("proName")

	// execServerURIにアクセスしてプログラム情報を取得
	getProgramInfoURL := execServerURI + "/program/all"
	resp, _ := http.Get(getProgramInfoURL)
	programInfoJSON, _ := ioutil.ReadAll(resp.Body)
	// programInfoJSONはこのようなJSON文字列になる
	// {
	// 	"AddZipExt": {
	// 		"help": "take any file.\noutput file that is added zip extension.\n\u003ch1\u003ethis is h1 tag\u003c/h1\u003e"
	// 	},
	//  "xxxxxxxxx": {
	// 		"help": "xxxxxxxxx"
	//  },
	// }

	fmt.Println(string(programInfoJSON))

	// programInfoJSONをマップに変換するための構造体
    var programInfoMap map[string]map[string]string

	// JSONのバイト文字列を定義したマップへ入れ込む。
    json.Unmarshal(programInfoJSON, &programInfoMap)

	var help string = programInfoMap[proName]["help"]
	var detailedHelp string = programInfoMap[proName]["detailedHelp"]

	// htmlに渡すための値を保持する構造体
	type dataToHTML struct {
		Name          string
		Help          template.HTML
		DetailedHelp  template.HTML
		ExecServerURI string
	}

	d := dataToHTML{Name: proName, ExecServerURI: execServerURI, Help: template.HTML(help), DetailedHelp: template.HTML(detailedHelp)}

	app.RenderTemplate(w, serveHTML, d)
	// TODO: topの方は軽い説明にし、execの方で具体的な説明にする。
}
