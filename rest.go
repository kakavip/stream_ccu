package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type HttpResponse struct {
}

func (HttpResponse) success(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	// json.NewEncoder((*w)).Encode(data)

	jsonResp, err := json.Marshal(&data)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}
	w.Write(jsonResp)
}

var response HttpResponse = HttpResponse{}
