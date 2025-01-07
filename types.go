package main

type Pokemon struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type PokemonAPIResponse struct {
	Count    int           `json:"count"`
	Next     string        `json:"next"`
	Previous string        `json:"previous"`
	Results  []interface{} `json:"results"`
}
