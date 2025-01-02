package controller

import (
	"encoding/json"
	"io"
	"net/http"
)

func checkForBodyItem(key string, w http.ResponseWriter, r *http.Request) (bool, any) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "unable to read body", http.StatusInternalServerError)
		return false, ""
	}
	defer r.Body.Close()

	var elements map[string]interface{}
	err = json.Unmarshal(body, &elements)
	if err != nil {
		http.Error(w, "unable to unmarshal json", http.StatusInternalServerError)
		return false, ""
	}

	if _, exists := elements[key]; exists {
		return true, elements[key]
	}
	return false, ""
}
