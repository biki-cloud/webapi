package application

import (
	"log"
	"os"

	"webapi/microservices/apigw/env"
)

type Application struct {
	Env      *env.Env
	InfoLog  *log.Logger
	ErrorLog *log.Logger
}

func New() *Application {

	e := env.New()

	// create logger for writing information and error messages.
	infoLog := log.New(os.Stdout, "INFO  ", log.LstdFlags)
	errorLog := log.New(os.Stderr, "ERROR ", log.LstdFlags|log.Lshortfile)

	a := &Application{
		Env:      e,
		InfoLog:  infoLog,
		ErrorLog: errorLog,
	}

	return a
}
