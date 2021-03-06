/*
プログラムサーバで使用するロガーを定義している。
*/

package log

import (
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

// GetLogger はファイルを受け取り、標準出力とログファイルへのログを書き出すロガーを返す。
func GetLogger(filePath string) *log.Logger {
	log.SetFlags(log.Ldate | log.Ltime)

	// make directory
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		log.Fatalln(err)
	}

	// make file
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln(err)
	}

	logger := log.New(file, "", log.Ldate|log.Ltime)

	mw := io.MultiWriter(os.Stdout, file)
	logger.SetOutput(mw)

	return logger
}

// RotateMiddleware はrotaterインタフェースを受け取り、ローテーションする。
func RotateMiddleware(next http.Handler, rotater Rotater) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := rotater.Rotate()
		if err != nil {
			logger.Fatalf("エラー！ログファイルのローテーションに失敗しました。err msg: %v \n", err.Error())
			return
		}
		next.ServeHTTP(w, r)
	})
}

// NullWriter logにセットすることで何も出力しない。
// log.New(new(NullWriter),.....)
type NullWriter int

func (NullWriter) Write([]byte) (int, error) { return 0, nil }
