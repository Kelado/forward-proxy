package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

const (
	pokemonApi = "https://pokeapi.co"
	limit      = 500
)

func proxyHandler(w http.ResponseWriter, r *http.Request) {
	page := r.URL.Query().Get("page")
	if page == "" {
		http.Error(w, "Missing 'page' parameter", http.StatusBadRequest)
		return
	}

	// if config.ExpirationPeriod == 0 {
	// 	http.Error(w, "No valid config file found", http.StatusInternalServerError)
	// 	return
	// }

	var data [][]interface{}

	data, found := cache.Get(page)
	if !found {
		log.Printf("Page %s not found. Make request\n", page)
		var url = fmt.Sprintf("%s%s?limit=%d", pokemonApi, page, limit)
		for {

			resp, err := doRequest(url)
			if err != nil {
				handleErrorResponse(w, *err)
				return
			}
			data = append(data, resp.Results)

			if resp.Next == "" {
				break
			}

			url = resp.Next
		}
		cache.Set(page, data)
	} else {
		log.Printf("Page %s found. Reading from cache\n", page)
	}

	marshalledData, err := json.Marshal(data)
	if err != nil {
		handleErrorResponse(w, ErrorResponse{http.StatusInternalServerError, "Could not marshal response."})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(marshalledData)
}

func doRequest(url string) (PokemonAPIResponse, *ErrorResponse) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return PokemonAPIResponse{}, &ErrorResponse{http.StatusInternalServerError, "Could not create request."}
	}
	client := &http.Client{Timeout: 2 * time.Second}

	resp, err := client.Do(req)
	if err != nil {
		return PokemonAPIResponse{}, &ErrorResponse{http.StatusInternalServerError, "Could not make external api call."}
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return PokemonAPIResponse{}, &ErrorResponse{http.StatusInternalServerError, "Could not read request body."}
	}

	var apiResponse PokemonAPIResponse
	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		return PokemonAPIResponse{}, &ErrorResponse{http.StatusInternalServerError, "Could not unmarshal response."}
	}

	return apiResponse, nil
}

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func handleErrorResponse(w http.ResponseWriter, error ErrorResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(error.Code)

	jsonResponse, err := json.Marshal(error)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Write(jsonResponse)
}
