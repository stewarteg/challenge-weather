package clients

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"

	dto "github.com/stewarteg/challenge-weather/dto"
)

type Client interface {
	ConsultCep(cep string) (*string, error)
	ConsultTemperatura(localidade string) (*float64, *float64, error)
}

type RealClient struct{}

func (r *RealClient) ConsultCep(cep string) (*string, error) {

	var validCep = regexp.MustCompile(`^\d{8}$`)

	// Verifica se o CEP está no formato correto
	if !validCep.MatchString(cep) {
		return nil, errors.New("invalid zipcode")
	}

	url := fmt.Sprintf("https://viacep.com.br/ws/%s/json/", cep)

	// Criando um cliente HTTP
	client := &http.Client{}

	// Criando a requisição GET
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Erro ao criar a requisição:", err)
		return nil, err
	}

	// Enviando a requisição
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Erro ao fazer a requisição:", err)
		return nil, err
	}
	defer resp.Body.Close() // Certifique-se de fechar o corpo da resposta

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("can not find zipcode")
	}
	// Lendo o corpo da resposta
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Erro ao ler o corpo da resposta:", err)
		return nil, err
	}

	var endereco dto.Endereco
	err = json.Unmarshal(body, &endereco)
	if err != nil {
		log.Fatalf("Erro ao decodificar JSON1: %v", err)
	}

	return &endereco.Localidade, nil
}

func (r *RealClient) ConsultTemperatura(localidade string) (*float64, *float64, error) {

	fmt.Println("localidade:", localidade)
	url := fmt.Sprintf("https://api.weatherapi.com/v1/current.json?key=8bc0b8761d4a475fb59193932250803&q=%s&aqi=no", localidade)

	// Criando um cliente HTTP
	client := &http.Client{}

	// Criando a requisição GET
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Erro ao criar a requisição:", err)
		return nil, nil, err
	}

	// Enviando a requisição
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Erro ao fazer a requisição:", err)
		return nil, nil, err
	}
	defer resp.Body.Close() // Certifique-se de fechar o corpo da resposta

	// Lendo o corpo da resposta
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Erro ao ler o corpo da resposta:", err)
		return nil, nil, err
	}

	var weatherapi dto.WeatherResponse
	err = json.Unmarshal(body, &weatherapi)
	if err != nil {
		log.Fatalf("Erro ao decodificar JSON2: %v", err)
	}

	return &weatherapi.Current.TempC, &weatherapi.Current.TempF, nil
}
