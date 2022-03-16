package middlewares

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

// HttpTrace はhttpリクエストをロギングする。
func HttpTrace(h http.Handler, logger *log.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rAddr := r.RemoteAddr
		method := r.Method
		path := r.URL.Path
		err := r.ParseForm()
		if err != nil {
			logger.Printf("r.ParseForm() err: %v\n", err.Error())
			return
		}
		logger.SetFlags(log.Ldate | log.Ltime)

		// jsやcssのGETはいらないログなので避ける。
		avoid := false
		logAvoidExts := []string{".css", ".js", ".png", "ico"}
		for _, ext := range logAvoidExts {
			if strings.Contains(path, ext) {
				avoid = true
			}
		}
		if !avoid {
			logger.Printf("%s %s%s", method, rAddr, path)
		}

		h.ServeHTTP(w, r)
	})
}

// AllowCORS はCORSを全部許容する
func AllowCORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Methods", "*")
		w.Header().Set("Access-Control-Allow-Headers", "*")

		if r.Method == http.MethodOptions {
			fmt.Println("options. return")
			w.WriteHeader(http.StatusOK)
			return
		}

		h.ServeHTTP(w, r)
	})
}
