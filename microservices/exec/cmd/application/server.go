package application

import (
	"net/http"
	"webapi/pkg/http/middlewares"
)

func NewServer(serverURI string, app *Application) *http.Server {
	handler := middlewares.HttpTrace(app.Routes(), app.InfoLog)
	handler = middlewares.AllowCORS(handler)

	srv := &http.Server{
		Addr:     serverURI,
		ErrorLog: app.ErrorLog,
		Handler:  handler,
	}

	app.InfoLog.Printf("starting server on %v\n", serverURI)

	return srv
}
