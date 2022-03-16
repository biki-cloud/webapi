package application

import (
	"net/http"
	"path/filepath"
	"webapi/pkg/os"
)

func (app *Application) Routes() *http.ServeMux {
	router := http.NewServeMux()

	// ユーザがこのハンドラにアクセスした場合は全てのサーバにアクセスし、全てのプログラムを表示する。
	router.HandleFunc("/user/top", app.Top)

	router.HandleFunc("/user/exec/", app.Exec)

	currentDir, err := os.GetCurrentDir()
	if err != nil {
		panic(err.Error())
	}
	router.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(filepath.Join(currentDir, "ui/static")))))

	return router
}
