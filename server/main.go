package main

import (
	"fmt"
	"net/http"

	"github.com/mjishu/pokeDate/cards"
)

func headers(w http.ResponseWriter, req *http.Request) {
}

func main() {
	fmt.Println("connected")
	http.HandleFunc("/cards", cards.CardsController)

	http.ListenAndServe(":8080", nil)
}
