package cards

import (
	"fmt"
	"net/http"
)

func CardsController(w http.ResponseWriter, r *http.Request) {
	fmt.Println("looking for cards to control")
}
