package api

import (
	"net/http"
	"encoding/json"

)

func errorResponse(w http.ResponseWriter, statusCode int, errorRes string, err error) {
	w.Header().Set("Content-Type", "application/json")
	log.Printf("\nError: %v\n", err)

	type errBody struct {
		Error string `json:"error"`
	}

	res := errBody{
		Error: errorRes,
	}
	jsonRes, err := json.Marshal(res)
	if err != nil {
		log.Printf("\nError: error while encoding response body - %v\n", err)
		errorResponse(w, http.InternalServerError, "Something went wrong", err)
		return
	}

	w.WriteHeader(statusCode)
	w.Write(jsonRes)
}

func successResponse(w http.ResponseWriter, statusCode int, res any) {
	w.Header().Set("Content-Type", "application/json")

	jsonRes, err := json.Marshal(res)
	if err != nil {
		log.Printf("\nError: error while encoding response body - %v\n", err)
		errorResponse(w, http.InternalServerError, "Something went wrong", err)
		return
	}

	w.WriteHeader(statusCode)
	w.Write(jsonRes)
}
