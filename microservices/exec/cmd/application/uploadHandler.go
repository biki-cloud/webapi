/*
アップロードされるハンドラーを定義したファイル。
*/

package application

import (
	"fmt"
	"net/http"
	"path/filepath"

	"webapi/microservices/exec/pkg/msgs"
	pkgHttpUpload "webapi/pkg/http/upload"
	pkgInt "webapi/pkg/int"
)

// APIUpload はファイルをアップロードするためのハンドラー。
func (app *Application) APIUpload(w http.ResponseWriter, r *http.Request) {
	uploadDir := filepath.Join(app.Cfg.FileServer.Dir, app.Cfg.FileServer.UploadDir)
	maxUploadSize := int64(pkgInt.MBToByte(int(app.Cfg.MaxUploadSizeMB)))

	_, err := pkgHttpUpload.UploadHelper(w, r, uploadDir, maxUploadSize)
	if err != nil {
		err = fmt.Errorf("APIExec: %v, err msg: %v", msgs.UploadFileSizeExceedError(app.Cfg.MaxUploadSizeMB), err.Error())
		app.ServerError(w, err)
	}

	_, err = fmt.Fprintf(w, msgs.UploadSuccess)
	if err != nil {
		app.ServerError(w, err)
	}
}
