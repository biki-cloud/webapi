/*
アップロードされるハンドラーを定義したファイル。
*/

package application

import (
	"fmt"
	"net/http"
	"path/filepath"
	msg "webapi/microservices/exec/pkg/msgs"
	http2 "webapi/pkg/http/upload"
	int2 "webapi/pkg/int"
)

// APIUpload はファイルをアップロードするためのハンドラー。
func (app *Application) APIUpload(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type,access-control-allow-origin, access-control-allow-headers")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	uploadDir := filepath.Join(app.Cfg.FileServer.Dir, app.Cfg.FileServer.UploadDir)
	maxUploadSize := int64(int2.MBToByte(int(app.Cfg.MaxUploadSizeMB)))

	_, err := http2.UploadHelper(w, r, uploadDir, maxUploadSize)
	if err != nil {
		err = fmt.Errorf("APIExec: %v, err msg: %v", msg.UploadFileSizeExceedError(app.Cfg.MaxUploadSizeMB), err.Error())
		app.ServerError(w, err)
	}

	_, err = fmt.Fprintf(w, msg.UploadSuccess)
	if err != nil {
		app.ServerError(w, err)
	}
}
