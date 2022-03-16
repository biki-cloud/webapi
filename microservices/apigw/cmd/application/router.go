package application

import (
	"net/http"
	http2 "webapi/pkg/http/handlers"
)

func (app *Application) Routes() *http.ServeMux {
	router := http.NewServeMux()

	// コマンドラインからはここにアクセスし、メモリ使用量が一番低いサーバのURLを返す。
	router.HandleFunc("/program-server/memory/minimum", app.GetMinimumMemoryServerHandler)

	// コマンドラインからここにアクセスし、プログラムがあるかつメモリ使用量が一番低いサーバのURLを返す。
	router.HandleFunc("/program-server/minimumMemory-and-hasProgram/", app.GetMinimumMemoryAndHasProgram)

	// 現在稼働しているサーバを返すAPI
	router.HandleFunc("/program-server/alive", app.GetAliveServersHandler)

	// 生きている全てのサーバのプログラムを取得してJSONで表示するAPI
	router.HandleFunc("/program-server/program/all", app.GetAllProgramsHandler)

	// このサーバが生きているかを判断するのに使用するハンドラ
	router.HandleFunc("/health", http2.HealthHandler)

	// このサーバプログラムのメモリ状態をJSONで表示するAPI
	router.HandleFunc("/health/memory", http2.GetRuntimeHandler)

	return router
}
