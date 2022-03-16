package contextManager_test

import (
	"fmt"
	"log"
	"net/http/httptest"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"webapi/microservices/exec/env"
	ec "webapi/microservices/exec/pkg/execution/contextManager"
	"webapi/pkg/http/request"
	pkgHttpUpload "webapi/pkg/http/upload"
	os2 "webapi/pkg/os"
	"webapi/test"
)

var (
	currentDir    string
	uploadFile    string
	dummyParameta string
	programName   string
	ctx           ec.ContextManager
	addrs         []string
	ports         []string
)

var cfg *env.Env

func init() {
	c, err := os2.GetCurrentDir()
	if err != nil {
		log.Fatalln(err.Error())
	}
	currentDir = c
	uploadFile = "uploadfile"
	dummyParameta = "dummyParameta"
	programName = "convertToJson"

	addrs, ports, err = test.GetStartedServers(1)
	fmt.Printf("addrs: %v, ports: %v \n", addrs, ports)
	if err != nil {
		panic(err)
	}

	err = os2.CreateSpecifiedFile(uploadFile, 200)
	if err != nil {
		panic(err.Error())
	}

	uploader := pkgHttpUpload.NewUploader()
	err = uploader.Upload(addrs[0]+"/upload", uploadFile)
	if err != nil {
		panic(err.Error())
	}

	cfg = env.New()
	ctx, err = GetDummyContextManager(cfg)
	if err != nil {
		panic(err)
	}
}

func tearDown() {
	os.RemoveAll(filepath.Join(currentDir, "fileserver"))
	os.Remove(uploadFile)
}

func GetDummyContextManager(cfg *env.Env) (ec.ContextManager, error) {
	uploadFile = "uploadfile"
	err := os2.CreateSpecifiedFile(uploadFile, 2)
	if err != nil {
		panic(err.Error())
	}

	fields := map[string]string{
		"proName":  "convertToJson",
		"parameta": "dummyParameta",
	}

	poster := request.NewPostGetter()
	r, err := poster.GetPostRequest("/exec/convertToJson", uploadFile, fields)
	if err != nil {
		panic(err.Error())
	}
	w := httptest.NewRecorder()

	var ctx ec.ContextManager
	ctx, err = ec.New(w, r, cfg)
	if err != nil {
		return nil, fmt.Errorf("GetDummyContextManager: %v", err)
	}

	return ctx, nil
}

func TestNewContextManager(t *testing.T) {
	// ctx.SetProgramOutDir, SetUploadedFilePathAndParametaはここで同時に試験できる。
	if ctx.Parameta() != dummyParameta {
		t.Errorf("ctx.Parameta(): %v, want: %v \n", ctx.Parameta(), dummyParameta)
	}
	if filepath.Base(ctx.UploadedFilePath()) != uploadFile {
		t.Errorf("ctx.UploadFilePath(): %v, want: %v \n", filepath.Base(ctx.UploadedFilePath()), uploadFile)
	}
	if ctx.ProgramName() != programName {
		t.Errorf("ctx.ProgramName(): %v, want: %v \n", ctx.ProgramName(), programName)
	}
	if filepath.Base(ctx.InputFilePath()) != uploadFile {
		t.Errorf("ctx.InputFilePath(): %v , want: %v \n", filepath.Base(ctx.InputFilePath()), uploadFile)
	}
	if !os2.FileExists(ctx.InputFilePath()) {
		t.Errorf("ctx.InputFilePath(%v) is not found \n", ctx.InputFilePath())
	}

	if !os2.FileExists(ctx.OutputDir()) {
		t.Errorf("ctx.OutputDir() is not found.")
	}

	if !reflect.DeepEqual(ctx.Env(), cfg) {
		t.Errorf("ctx.Env(%v) is not equal env(%v) \n", ctx.Env(), cfg)
	}

	if !os2.FileExists(ctx.ProgramTempDir()) {
		t.Errorf("ctx.ProgramTempDir is not found")
	}

	t.Cleanup(func() {
		tearDown()
	})
}
