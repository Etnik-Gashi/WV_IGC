package main

import (
	"fmt"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request){

	fmt.Fprintf(w, "<h1>Hi, here are igc info</h1>")
}

func main() {
	http.HandleFunc("/igcinfo",handler)
	http.ListenAndServe(":8080",nil)
}