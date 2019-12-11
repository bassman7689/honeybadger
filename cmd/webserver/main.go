package main

import "net/http"

func download(w http.ResponseWriter, r *http.Request) {
}

func main() {
	http.Handle("/", http.FileServer(http.Dir("/static")))
	http.ListenAndServe("0.0.0.0:8080", nil)
}
