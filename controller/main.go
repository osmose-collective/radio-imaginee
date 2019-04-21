package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
)

func main() {
	fmt.Println("Starting server on :8042...")
	r := chi.NewRouter()
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})
	http.ListenAndServe(":8042", r)
}
