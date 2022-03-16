package application

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"runtime/debug"
)

func (app *Application) ServerError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.ErrorLog.Output(2, trace)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// PrintAsJSON webページにJSONを表示する。
// jsonの元になる構造体を渡す。
func (app *Application) PrintAsJSON(w http.ResponseWriter, jsonStruct interface{}) {
	// jsonに変換
	b, err := json.MarshalIndent(jsonStruct, "", "    ")
	if err != nil {
		app.ServerError(w, err)
	}

	w.Header().Set("Content-Type", "Application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, string(b))
}

// RenderTemplate webページのhtmlにデータを渡し、表示させる
// htmlPathは絶対パスで記述する。
func (app *Application) RenderTemplate(w http.ResponseWriter, htmlPath string, data interface{}) {
	if !filepath.IsAbs(htmlPath) {
		app.ServerError(w, fmt.Errorf("serve htmlPath(%v) is not found.", htmlPath))
	}

	t, err := template.ParseFiles(htmlPath)
	if err != nil {
		app.ServerError(w, err)
	}

	if err := t.Execute(w, data); err != nil {
		app.ServerError(w, err)
	}
}
