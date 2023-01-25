package json

import (
	"encoding/json"
	"net/http"
)

func ParseJson(r *http.Request, body interface{}) error {
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return err
	}
	return nil
}

func SendJson(w http.ResponseWriter, body interface{}) error {
	w.Header().Set("Content-type", "application/json")
	if err := json.NewEncoder(w).Encode(body); err != nil {
		return err
	}
	return nil
}

func SendError(w http.ResponseWriter, incomeErr error) {
	err := SendJson(w, map[string]interface{}{"status": incomeErr.Error()})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
