package api

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

var router *mux.Router

func init() {
	router = mux.NewRouter()
	router.StrictSlash(true)

	// healthz
	router.Handle("/healthz", http.HandlerFunc(healthzHandler)).Methods("GET")
}

// Serve starts an HTTP server that handles all inbound requests. This function blocks while the
// server runs, so it should be run in its own goroutine.
func Serve(port int, errCh chan<- error) {
	logger.Infof("starting API server on port %d", port)
	host := fmt.Sprintf(":%d", port)
	if err := http.ListenAndServe(host, router); err != nil {
		logger.Errorf("api health check server (%s)", err)
		errCh <- err
	}
}
