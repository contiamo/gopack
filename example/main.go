package main

import "net/http"

func main() {
	http.Handle("/", packHandler)
	http.ListenAndServe(":8080", nil)
}
