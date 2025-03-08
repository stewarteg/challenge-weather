package controller

import (
	"encoding/json"
	"net/http"

	"github.com/stewarteg/challenge-weather/clients"
)

type Controller struct {
	Client clients.Client
}

func (c *Controller) ConsultaTemperatura(w http.ResponseWriter, r *http.Request) {

	cep := r.URL.Query().Get("cep")

	localidade, err := c.Client.ConsultCep(cep)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	celcius, fahrenhait, _ := c.Client.ConsultTemperatura(*localidade)

	tempResponse := map[string]float64{
		"temp_C": *celcius,
		"temp_F": *fahrenhait,
		"temp_K": (*celcius + 273),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tempResponse)
}
