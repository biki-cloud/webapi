package application

import (
	"net/http"
	pkgHttpMiddlewares "webapi/pkg/http/middlewares"
)

func NewServer(serverURI string, app *Application) *http.Server {
	handler := pkgHttpMiddlewares.HttpTrace(app.Routes(), app.InfoLog)

	srv := &http.Server{
		Addr:     serverURI,
		ErrorLog: app.ErrorLog,
		Handler:  handler,
	}

	app.InfoLog.Printf("starting server on %v\n", serverURI)

	return srv
}
