package post_test

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"testing"
	p "webapi/cli/pkg/post"
	app2 "webapi/microservices/exec/cmd/application"
	os2 "webapi/pkg/os"
)

var (
	currentDir string
	uploadFile string
)

func init() {
	c, err := os2.GetCurrentDir()
	if err != nil {
		log.Fatalln(err.Error())
	}
	currentDir = c
	uploadFile = "uploadfile"

	go func() {
		a := app2.New()
		srv := app2.NewServer("localhost:8881", a)
		if err := srv.ListenAndServe(); err != nil {
			panic(err.Error())
		}
	}()
}

func tearDown() {
	os.RemoveAll(filepath.Join(currentDir, "fileserver"))
	os.Remove(uploadFile)
}

func TestProcessFileOnServer(t *testing.T) {
	// create upload file
	err := os2.CreateSpecifiedFile(uploadFile, 200)
	if err != nil {
		t.Errorf("err from CreateSpecifiedFile: %v \n", err.Error())
	}

	// ファイル上で処理させる
	basename := filepath.Base(uploadFile)
	programName := "convertToJson"
	serverURL := "http://127.0.0.1:8881"
	parameta := ""

	// 処理させるためのurl
	url := fmt.Sprintf("%v/api/exec/%v", serverURL, programName)

	processor := p.New()

	res, err := processor.Post(programName, url, basename, parameta)
	if err != nil {
		t.Errorf("err from Post: %v \n", err.Error())
	}

	outBase := filepath.Base(res.OutURLs()[0])
	if outBase != "uploadfile.json1" {
		t.Errorf("output file is not %v \n", "uploadfile.json")
	}

	t.Cleanup(func() {
		tearDown()
	})
}
