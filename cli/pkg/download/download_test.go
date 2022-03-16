package download_test

import (
	"log"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"webapi/cli/pkg/download"
	pkgHttpUpload "webapi/pkg/http/upload"
	pkgOs "webapi/pkg/os"
	"webapi/test"
)

var (
	currentDir string
	uploadFile string
	addrs      []string
)

func init() {
	c, err := pkgOs.GetCurrentDir()
	if err != nil {
		log.Fatalln(err.Error())
	}
	currentDir = c
	uploadFile = "uploadfile"

	addrs, _, err = test.GetStartedServers(1)
}

func tearDown() {
	os.RemoveAll(filepath.Join(currentDir, "fileserver"))
	os.Remove(uploadFile)
}

func TestDownload(t *testing.T) {
	// create upload file
	err := pkgOs.CreateSpecifiedFile(uploadFile, 200)
	if err != nil {
		t.Errorf("err from CreateSpecifiedFile: %v \n", err.Error())
	}

	// upload
	uploader := pkgHttpUpload.NewUploader()
	err = uploader.Upload(addrs[0]+"/upload", uploadFile)
	if err != nil {
		t.Errorf("err from uploadHelper: %v \n", err.Error())
	}

	url := addrs[0] + "/fileserver/upload/uploadFile"

	downloader := download.New()

	// mkdir
	err = os.Mkdir("tmp", os.ModePerm)
	if err != nil {
		t.Errorf("err from Mkdir(): %v \n", err.Error())
	}

	done := make(chan error, 10)
	var wg sync.WaitGroup

	wg.Add(1)
	downloader.Download(url, "tmp", done, &wg, pkgOs.NewMover())

	t.Cleanup(func() {
		tearDown()
		os.RemoveAll("tmp")
	})
}
