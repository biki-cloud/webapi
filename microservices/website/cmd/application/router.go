package application

import (
	"net/http"
	"path/filepath"
	pkgHttpHandlers "webapi/pkg/http/handlers"

	pkgOs "webapi/pkg/os"
)

func (app *Application) Routes() *http.ServeMux {
	router := http.NewServeMux()

	// ユーザがこのハンドラにアクセスした場合は全てのサーバにアクセスし、全てのプログラムを表示する。
	router.HandleFunc("/user/top", app.Top)

	router.HandleFunc("/user/exec/", app.Exec)

	currentDir, err := pkgOs.GetCurrentDir()
	if err != nil {
		panic(err.Error())
	}
	router.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(filepath.Join(currentDir, "ui/static")))))

	// このサーバが生きているかを判断するのに使用するハンドラ
	router.HandleFunc("/health", pkgHttpHandlers.HealthHandler)

	return router
}
