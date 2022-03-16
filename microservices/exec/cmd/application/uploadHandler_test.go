package application_test

import (
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"webapi/microservices/exec/env"
	http2 "webapi/pkg/http/request"
	uf "webapi/pkg/os"
)

func TestUploadHandler(t *testing.T) {
	tests := []struct {
		testName           string
		fileName           string
		fileSize           int64
		cfgMaxUploadSizeMB int64
		uploadIsSuccess    bool
	}{
		{
			testName:           "success test",
			fileName:           "uploadFile",
			fileSize:           200,
			cfgMaxUploadSizeMB: 300,
			uploadIsSuccess:    true,
		},
		{
			testName:           "fail test",
			fileName:           "uploadFile2",
			fileSize:           300,
			cfgMaxUploadSizeMB: 100,
			uploadIsSuccess:    false,
		},
	}

	for _, tt := range tests {
		err := uf.CreateSpecifiedFile(tt.fileName, tt.fileSize)
		if err != nil {
			t.Fatalf(err.Error())
		}
		t.Run(tt.testName, func(j *testing.T) {
			cfg := env.New()
			cfg.MaxUploadSizeMB = tt.cfgMaxUploadSizeMB
			testUpload(t, tt.uploadIsSuccess, tt.fileName, cfg)
		})
	}
}

func testUpload(t *testing.T, uploadIsSuccess bool, uploadFile string, cfg *env.Env) {
	postGetter := http2.NewPostGetter()
	r, err := postGetter.GetPostRequest("/upload", uploadFile, map[string]string{"dummy": "x"})
	if err != nil {
		t.Fatalf(err.Error())
	}

	w := httptest.NewRecorder()

	a.APIUpload(w, r)

	// アップロードされているか
	uploadedPath := filepath.Join(cfg.FileServer.Dir, "upload", uploadFile)
	if uf.FileExists(uploadedPath) != uploadIsSuccess {
		t.Errorf("got: %v, want: %v", uf.FileExists(uploadedPath), uploadIsSuccess)
	}

	t.Cleanup(func() {
		err := os.RemoveAll("fileserver")
		if err != nil {
			panic(err)
		}
		err = os.Remove(uploadFile)
		if err != nil {
			panic(err)
		}
	})
}