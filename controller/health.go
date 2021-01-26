package controller

import "net/http"

func Health(w http.ResponseWriter, r *http.Request){
	_, _ = w.Write([]byte("Looking good here!\n"))
}
