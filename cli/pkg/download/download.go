/*
サーバからファイルをダウンロードする機能を提供するパッケージ
*/

package download

import (
	"path/filepath"
	"sync"
	"webapi/pkg/os"
)

var (
	currentDir string
	uploadFile string
)

type Downloader interface {
	// Download はダウンロードしたいファイルURLを入れて、outputDirへダウンロードする。
	Download(url string, outputDir string, done chan error, wg *sync.WaitGroup, mover os.Mover)
}

func New() Downloader {
	return &downloader{}
}

type downloader struct{}

// Download はダウンロードしたいファイルURLを入れて、outputDirへダウンロードする。
func (d *downloader) Download(url, outputDir string, done chan error, wg *sync.WaitGroup, mover os.Mover) {
	defer wg.Done() // 関数終了時にデクリメント
	basename := filepath.Base(url)
	newLocation := filepath.Join(outputDir, basename)
	command := "curl -L " + "-o " + newLocation + " " + url
	_, _, err := os.SimpleExec(command)
	if err != nil {
		done <- err
		return
	}
}
