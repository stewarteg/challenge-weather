package controller

import (
	"encoding/json"
	"net/http"

	"github.com/stewarteg/goexpert/challenge-weather/clients"
)

func ConsultaTemperatura(w http.ResponseWriter, r *http.Request) {

	cep := r.URL.Query().Get("cep")

	localidade, err := clients.ConsultCep(cep)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	celcius, fahrenhait, _ := clients.ConsultTemperatura(*localidade)

	tempResponse := map[string]float64{
		"temp_C": *celcius,
		"temp_F": *fahrenhait,
		"temp_K": (*celcius + 273),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tempResponse)
}
