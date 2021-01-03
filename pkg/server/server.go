package server

import (
	"net/http"
	"time"
)

func RunHTTPServer(port string, handler http.Handler) error {
	s := &http.Server{
		Addr:           ":" + port,
		Handler:        handler,
		MaxHeaderBytes: 1 << 20,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
	}

	return s.ListenAndServe()
}
