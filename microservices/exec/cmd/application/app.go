package application

import (
	"log"
	"os"

	"webapi/microservices/exec/env"
)

type Application struct {
	Cfg      *env.Env
	InfoLog  *log.Logger
	ErrorLog *log.Logger
}

func New() *Application {

	cfg := env.New()

	// create logger for writing information and error messages.
	infoLog := log.New(os.Stdout, "INFO  ", log.LstdFlags)
	errorLog := log.New(os.Stderr, "ERROR ", log.LstdFlags|log.Lshortfile)

	a := &Application{
		Cfg:      cfg,
		InfoLog:  infoLog,
		ErrorLog: errorLog,
	}

	return a
}
