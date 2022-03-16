/*
全てのハンドラをセットしたrouterを返すGetServeMuxを定義している。
*/

package application

import (
	"net/http"
	http2 "webapi/pkg/http/handlers"
)

// Routes ハンドラをセットしたrouterを返す。
func (app *Application) Routes() *http.ServeMux {
	r := http.NewServeMux()

	// ファイルサーバーの機能のハンドラ
	// Env.FileServer.Dir以下のファイルをwebから見ることができる。
	fileServer := "/" + app.Cfg.FileServer.Dir + "/"
	r.Handle(fileServer, http.StripPrefix(fileServer, http.FileServer(http.Dir(app.Cfg.FileServer.Dir))))

	// 登録プログラムを実行させるAPI
	r.HandleFunc("/api/exec/", app.APIExec)

	// ファイルをアップロードするAPI
	r.HandleFunc("/api/upload", app.APIUpload)

	// このサーバプログラムのメモリ状態をJSONで表示するAPI
	r.HandleFunc("/health/memory", http2.GetRuntimeHandler)

	// プログラムサーバに登録してあるプログラム一覧をJSONで表示するAPI
	r.HandleFunc("/program/all", app.AllHandler)

	// このサーバが生きているかを判断するのに使用するハンドラ
	r.HandleFunc("/health", http2.HealthHandler)

	return r
}
