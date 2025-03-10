package controller

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/stewarteg/challenge-weather/clients"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

type Controller struct {
	Client clients.Client
}

func (c *Controller) ConsultaTemperatura(w http.ResponseWriter, r *http.Request) {
	// Iniciar o span para o método ConsultaTemperatura

	log.Println("Iniciando ConsultaTemperatura - Criando Span")
	defer log.Println("Finalizando ConsultaTemperatura - Span fechado")

	tracer := otel.Tracer("controller")
	ctx, span := tracer.Start(r.Context(), "ConsultaController")
	defer span.End()

	// Obter o CEP da query string
	cep := r.URL.Query().Get("cep")
	span.SetAttributes(attribute.String("input.cep", cep))

	// Consultar o CEP
	localidade, err, statusCode := c.Client.ConsultCep(ctx, cep)
	if err != nil {
		span.RecordError(err)
		http.Error(w, err.Error(), statusCode)
		return
	}
	span.SetAttributes(attribute.String("response.localidade", *localidade))

	// Consultar a temperatura
	celcius, fahrenhait, err := c.Client.ConsultTemperatura(ctx, *localidade)
	if err != nil {
		span.RecordError(err)
		http.Error(w, "error fetching temperature", http.StatusInternalServerError)
		return
	}
	span.SetAttributes(
		attribute.Float64("response.temp_C", *celcius),
		attribute.Float64("response.temp_F", *fahrenhait),
		attribute.Float64("response.temp_K", *celcius+273),
	)

	// Criar a resposta
	tempResponse := map[string]interface{}{
		"city":   *localidade,
		"temp_C": *celcius,
		"temp_F": *fahrenhait,
		"temp_K": (*celcius + 273),
	}

	// Configurar o cabeçalho e enviar a resposta
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tempResponse)
}
