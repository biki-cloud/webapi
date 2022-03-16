package upload_test

import (
	"log"
	"os"
	"path/filepath"
	"testing"
	"webapi/microservices/exec/env"
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
	if err != nil {
		log.Fatalln(err)
	}
}

func tearDown() {
	os.RemoveAll(filepath.Join(currentDir, "fileserver"))
	os.Remove(uploadFile)
}

func TestUpload(t *testing.T) {
	err := pkgOs.CreateSpecifiedFile(uploadFile, 200)
	if err != nil {
		t.Errorf("err from CreateSpecifiedFile: %v \n", err.Error())
	}

	uploader := pkgHttpUpload.NewUploader()
	err = uploader.Upload(addrs[0]+"/api/upload", uploadFile)
	if err != nil {
		t.Errorf("err from uploadHelper: %v \n", err.Error())
	}

	uploadedFilePath := filepath.Join(currentDir, env.New().FileServer.Dir, "upload", uploadFile)
	if !pkgOs.FileExists(uploadedFilePath) {
		t.Errorf("uploadedPath(%v) is not exists. \n", uploadedFilePath)
	}

	t.Cleanup(func() {
		tearDown()
	})
}
