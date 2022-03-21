package application_test

import (
	"fmt"
	"log"
	"net/http/httptest"
	"testing"
	"time"

	apigwApp "webapi/microservices/apigw/cmd/application"
	apigwEnv "webapi/microservices/apigw/env"
	execApp "webapi/microservices/exec/cmd/application"
	websiteApp "webapi/microservices/website/cmd/application"
	websiteEnv "webapi/microservices/website/env"
	pkgHttpURL "webapi/pkg/http/url"
)

var (
	app *websiteApp.Application
)

func serversSet() {
	app = websiteApp.New()

	// serve exec server.
	for _, url := range apigwEnv.New().ExecServers {
		port := pkgHttpURL.GetPortFromURL(url)
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
	for _, url := range websiteEnv.New().APIGateWayServers {
		port := pkgHttpURL.GetPortFromURL(url)
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
