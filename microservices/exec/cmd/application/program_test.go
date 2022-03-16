package application_test

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	execApp "webapi/microservices/exec/cmd/application"
	"webapi/microservices/exec/pkg/execution/outputManager"
	"webapi/microservices/exec/pkg/msgs"
	pkgHttpRequest "webapi/pkg/http/request"
	pkgHttpUpload "webapi/pkg/http/upload"
	pkgHttpURL "webapi/pkg/http/url"
	pkgOs "webapi/pkg/os"
)

var (
	uploadFile string
	a          *execApp.Application
	srv        *http.Server
)

func set() {
	uploadFile = "uploadfile"

	addr, err := pkgHttpURL.GetLoopBackURL()
	if err != nil {
		log.Fatalln(err.Error())
	}

	go func() {
		a = execApp.New()
		srv = execApp.NewServer(addr, a)
		if err := srv.ListenAndServe(); err != nil {
			a.ErrorLog.Fatalln(err.Error())
		}
	}()

	err = pkgOs.CreateSpecifiedFile(uploadFile, 200)
	if err != nil {
		panic(err.Error())
	}

	uploader := pkgHttpUpload.NewUploader()
	err = uploader.Upload(addr+"/upload", uploadFile)
	if err != nil {
		panic(err.Error())
	}
}

func tearDown() {
	os.RemoveAll("fileserver")
	os.Remove(uploadFile)
}

func TestProgramHandler(t *testing.T) {
	set()
	// 保持しているプログラムの場合
	programName := "convertToJson"
	t.Run("Success test", func(t *testing.T) {
		rr, out := testProgramHandler(t, programName)
		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusOK)
		}

		if len(out.OutputURLs) != 2 {
			t.Errorf("len(out.OutputURLs): %v, want: 1", len(out.OutputURLs))
		}

		if out.Stdout == "" {
			t.Errorf("out.Stdout is empty.")
		}

		if out.Stderr != "" {
			t.Errorf("out.Stderr is not empty.")
		}

		if out.StaTus != msgs.OK {
			t.Errorf("out.StaTus: %v, want : %v", out.StaTus, msgs.SERVERERROR)
		}

		if out.Errormsg != "" {
			t.Errorf("out.Errormsg: %v, want: %v", out.Errormsg, "")
		}
	})

	// 保持していないプログラム名の場合
	programName = "nothingProgram"
	t.Run("fail test", func(t *testing.T) {
		rr, out := testProgramHandler(t, programName)
		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusOK)
		}

		if len(out.OutputURLs) != 0 {
			t.Errorf("len(out.OutputURLs): %v, want: 0", len(out.OutputURLs))
		}

		if out.Stdout != "" {
			t.Errorf("out.Stdout is not empty.")
		}

		if out.Stderr != "" {
			t.Errorf("out.Stderr is not empty.")
		}

		if out.StaTus != msgs.SERVERERROR {
			t.Errorf("out.StaTus: %v, want : %v", out.StaTus, msgs.SERVERERROR)
		}

		expected := programName + " is not found."
		if out.Errormsg != expected {
			t.Errorf("out.Errormsg: %v, want: %v", out.Errormsg, expected)
		}
	})

	t.Cleanup(func() {
		tearDown()
	})
}

func testProgramHandler(t *testing.T, proName string) (*httptest.ResponseRecorder, *outputManager.OutputInfo) {
	uploadFile = "uploadfile"
	err := pkgOs.CreateSpecifiedFile(uploadFile, 2)
	if err != nil {
		panic(err.Error())
	}

	fields := map[string]string{
		"proName":  proName,
		"parameta": "dummyParameta",
	}
	poster := pkgHttpRequest.NewPostGetter()
	r, err := poster.GetPostRequest("/exec/"+proName, uploadFile, fields)
	if err != nil {
		panic(err.Error())
	}
	w := httptest.NewRecorder()

	a.APIExec(w, r)

	var out *outputManager.OutputInfo

	if err := json.Unmarshal(w.Body.Bytes(), &out); err != nil {
		t.Errorf("jsonUnmarshal fail(msg: %v). body is %v", err.Error(), w.Body.String())
	}

	return w, out
}

func TestProgramAllHandler(t *testing.T) {
	set()
	r := httptest.NewRequest("GET", "/json/program/all", nil)
	w := httptest.NewRecorder()

	a.AllHandler(w, r)

	if status := w.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got: %v want: %v",
			status, http.StatusOK)
	}

	var m map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &m)
	if err != nil {
		t.Errorf("%v is not json format. err msg: %v", w.Body.String(), err.Error())
	}

	for _, v := range m {
		m2, ok := v.(map[string]interface{})
		if !ok {
			t.Errorf("got: %v, want: true", ok)
		}
		if _, ok = m2["help"]; !ok {
			t.Errorf("got: %v, want: true.", ok)
		}
	}

	expectedProgramNames := []string{"convertToJson", "err", "sleep"}
	for _, name := range expectedProgramNames {
		if !strings.Contains(w.Body.String(), name) {
			t.Errorf("%v is not contaned of %v", name, w.Body.String())
		}
	}

	t.Cleanup(func() {
		tearDown()
	})
}
