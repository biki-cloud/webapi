package application

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
)

// Download コンテンツをダウンロードさせるためのAPI
func (app *Application) Download(w http.ResponseWriter, r *http.Request) {
	//downloadPath eg -> fileserver/xxxxx/xxxxx/test.txt
	downloadPath := r.URL.Path[len("/download/"):]

	// Open file
	f, err := os.Open(downloadPath)
	if err != nil {
		app.ServerError(w, err)
	}
	defer f.Close()

	// octet-streamにすることでボタンを押した時に開かずにダウンロードされるようになる
	w.Header().Set("Content-type", "application/octet-stream")
	// ダウンロードした時にファイル名が付与されるようにする
	w.Header().Set("Content-Disposition", "attachment; filename="+filepath.Base(downloadPath))

	// Stream to response
	if _, err := io.Copy(w, f); err != nil {
		app.ServerError(w, err)
	}
}
