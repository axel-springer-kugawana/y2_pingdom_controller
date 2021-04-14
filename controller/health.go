package controller

import "net/http"

// Health Simple health check
func Health(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte("Looking good here!\n"))
}
