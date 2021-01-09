package handler

import (
	"encoding/json"
	"net/http"
)

type Status struct {
	Status string `json:"status"`
}

type Error struct {
	Err string `json:"error"`
}

func HTTPError(w http.ResponseWriter, err string, code int) {
	e := Error{
		Err: err,
	}

	w.WriteHeader(code)
	json.NewEncoder(w).Encode(e)
}
