package controller_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stewarteg/challenge-weather/controller"
	"github.com/stretchr/testify/assert"
)

type MockClient struct{}

func (m *MockClient) ConsultCep(cep string) (*string, error, int) {
	oi := "oi"
	return &oi, nil, 200
}

func (m *MockClient) ConsultTemperatura(localidade string) (*float64, *float64, error) {
	celsius := 25.0
	fahrenheit := 77.0
	return &celsius, &fahrenheit, nil
}

func TestConsultaTemperatura(t *testing.T) {
	mockClient := &MockClient{}
	controller := &controller.Controller{Client: mockClient}

	// Criar uma requisição HTTP simulada
	req, err := http.NewRequest("GET", "/consulta-temperatura?cep=12345678", nil)
	assert.NoError(t, err)

	// Criar um ResponseRecorder para capturar a resposta
	rr := httptest.NewRecorder()

	// Chamar a função que será testada
	handler := http.HandlerFunc(controller.ConsultaTemperatura)
	handler.ServeHTTP(rr, req)

	// Verificar o status da resposta
	assert.Equal(t, http.StatusOK, rr.Code)

	// Verificar o corpo da resposta
	expectedResponse := map[string]float64{
		"temp_C": 25.0,
		"temp_F": 77.0,
		"temp_K": 298.0,
	}
	var actualResponse map[string]float64
	err = json.Unmarshal(rr.Body.Bytes(), &actualResponse)
	assert.NoError(t, err)
	assert.Equal(t, expectedResponse, actualResponse)
}
