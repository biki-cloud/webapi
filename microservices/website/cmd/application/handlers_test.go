package application_test

import (
	"fmt"
	"log"
	"net/http/httptest"
	"testing"
	"time"
	apigwApp "webapi/microservices/apigw/cmd/application"
	apigwConfig "webapi/microservices/apigw/env"
	execApp "webapi/microservices/exec/cmd/application"
	"webapi/microservices/website/cmd/application"
	"webapi/microservices/website/env"
	http2 "webapi/pkg/http/url"
)

var (
	app *application.Application
)

func serversSet() {
	app = application.New()

	// serve exec server.
	for _, url := range apigwConfig.New().ExecServers {
		port := http2.GetPortFromURL(url)
		var done chan error
		go func() {
			srv := execApp.NewServer(":"+port, execApp.New())
			done <- srv.ListenAndServe()
		}()
		select {
		case err := <-done:
			if err != nil {
				log.Fatalf("%v", err.Error())
			}
		case <-time.After(1 * time.Second):
			fmt.Println("serve success")
		}
	}

	// serve apigw server.
	for _, url := range env.New().APIGateWayServers {
		port := http2.GetPortFromURL(url)
		var done chan error
		go func() {
			srv := apigwApp.NewServer(":"+port, apigwApp.New())
			done <- srv.ListenAndServe()
		}()
		select {
		case err := <-done:
			if err != nil {
				log.Fatalf("%v", err.Error())
			}
		case <-time.After(1 * time.Second):
			fmt.Println("serve success")
		}
	}
}

func TestApplication_Top(t *testing.T) {
	serversSet()
	r := httptest.NewRequest("GET", "/user/top", nil)
	w := httptest.NewRecorder()

	app.Top(w, r)

	fmt.Println(w.Body.String())

	//func TestUserTopHandler(t *testing.T) {
	//	r, _ := http.NewRequest(http.MethodGet, "/user/top", nil)
	//	w := httptest.NewRecorder()
	//
	//	app.UserTopHandler(w, r)
	//
	//	if w.Code != http.StatusOK {
	//		t.Errorf("got %v, want %v", w.Code, http.StatusOK)
	//	}
	//
	//	expectedHtmlTags := []string{
	//		"<td>convertToJson</td>",
	//		"<td>err</td>",
	//		"<td>move</td>",
	//		"<td>sleep</td>",
	//	}
	//	for _, tag := range expectedHtmlTags {
	//		if !strings.Contains(w.Body.String(), tag) {
	//			t.Errorf("expect containing this html tag: %v, but html doesn't contain: %v \n", tag, w.Body.String())
	//		}
	//	}
	//
	//	t.Cleanup(func() {
	//		tearDown()
	//	})
	//}
}
