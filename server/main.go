package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/mjishu/pokeDate/controller"
	"github.com/mjishu/pokeDate/database"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/cards", controller.CardsController)
	mux.HandleFunc("/animals", controller.AnimalController)
	mux.HandleFunc("/organizations/animals", controller.OrganizationController) //? change to /orgnaizations and make a new controller called organizations Cotnroller
	mux.HandleFunc("/shots", controller.ShotController)

	database.Database()

	port := ":8080"

	fmt.Println("listening on port " + port)
	err := http.ListenAndServe(port, mux)
	log.Fatal(err)
}
