package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/mjishu/pokeDate/cards"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", index)
	mux.HandleFunc("/cards", cards.CardsController)

	port := ":8080"

	fmt.Println("listening on port " + port)
	err := http.ListenAndServe(port, mux)
	log.Fatal(err)
}

func index(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	switch r.Method {
	case http.MethodGet:
		fmt.Println("Get was called")
	case http.MethodPost:
		fmt.Println(("Post was called"))
	case http.MethodOptions:
		w.Header().Set("Allow", "GET, POST, OPTIONS")
		w.WriteHeader((http.StatusNoContent))
	default:
		w.Header().Set("Allow", "GET,POST,OPTIONS")
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}
