package clients

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"

	dto "github.com/stewarteg/challenge-weather/dto"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

type Client interface {
	ConsultCep(ctx context.Context, cep string) (*string, error, int)
	ConsultTemperatura(ctx context.Context, localidade string) (*float64, *float64, error)
}

type RealClient struct{}

func (r *RealClient) ConsultCep(ctx context.Context, cep string) (*string, error, int) {
	tracer := otel.Tracer("client")
	ctx, span := tracer.Start(ctx, "ConsultCep")
	defer span.End()

	// Validação do CEP
	var validCep = regexp.MustCompile(`^\d{8}$`)
	if !validCep.MatchString(cep) {
		span.SetAttributes(attribute.String("error", "invalid zipcode"))
		return nil, errors.New("invalid zipcode"), http.StatusUnprocessableEntity
	}

	// Fazer a requisição HTTP
	//url := fmt.Sprintf("http://viacep.com.br/ws/%s/json/", cep)
	url := fmt.Sprintf("http://viacep.com.br/ws/%s/json/", cep)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		span.RecordError(err)
		return nil, err, http.StatusInternalServerError
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		span.RecordError(err)
		return nil, err, http.StatusInternalServerError
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		span.SetAttributes(attribute.Int("http.status_code", resp.StatusCode))
		return nil, errors.New("can not find zipcode"), http.StatusNotFound
	}

	// Processar a resposta
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		span.RecordError(err)
		return nil, err, http.StatusInternalServerError
	}

	var endereco dto.Endereco
	err = json.Unmarshal(body, &endereco)
	if err != nil {
		span.RecordError(err)
		return nil, err, http.StatusInternalServerError
	}

	if endereco.Localidade == "" {
		span.RecordError(errors.New("can not find zipcode"))
		return nil, errors.New("can not find zipcode"), 404
	}

	span.SetAttributes(attribute.String("response.localidade", endereco.Localidade))
	return &endereco.Localidade, nil, http.StatusOK
}

func (r *RealClient) ConsultTemperatura(ctx context.Context, localidade string) (*float64, *float64, error) {
	tracer := otel.Tracer("client")
	ctx, span := tracer.Start(ctx, "ConsultTemperatura")
	defer span.End()

	// Codificar a localidade para a URL
	encodedLocalidade := url.QueryEscape(localidade)
	url := fmt.Sprintf("http://api.weatherapi.com/v1/current.json?key=8bc0b8761d4a475fb59193932250803&q=%s&aqi=no", encodedLocalidade)

	// Fazer a requisição HTTP
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		span.RecordError(err)
		return nil, nil, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		span.RecordError(err)
		return nil, nil, err
	}
	defer resp.Body.Close()

	// Processar a resposta
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		span.RecordError(err)
		return nil, nil, err
	}

	var weatherapi dto.WeatherResponse
	err = json.Unmarshal(body, &weatherapi)
	if err != nil {
		span.RecordError(err)
		return nil, nil, err
	}

	span.SetAttributes(
		attribute.Float64("response.temp_C", weatherapi.Current.TempC),
		attribute.Float64("response.temp_F", weatherapi.Current.TempF),
	)
	return &weatherapi.Current.TempC, &weatherapi.Current.TempF, nil
}
