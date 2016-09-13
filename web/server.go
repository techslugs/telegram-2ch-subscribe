package web

import (
	"fmt"
	"log"
	"net/http"
	"github.com/techslugs/telegram-2ch-subscribe/version"
)

func StartServer(ipAddress string, port int) error {
	http.HandleFunc("/", versionHandler)
	fullAddress := fmt.Sprintf("%s:%d", ipAddress, port)
	log.Printf("[web] Starting web server on http://%s", fullAddress)
	return http.ListenAndServe(fullAddress, nil)
}

func versionHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("[web] Serving request to %s", r.URL)
	fmt.Fprintf(w, "Version: %s", version.Version)
}
