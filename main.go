package main

import (
	"fmt"
	"net/http"

	"github.com/stewarteg/challenge-weather/clients"
	"github.com/stewarteg/challenge-weather/controller"
)

func main() {

	// Crie uma instância do cliente real
	client := &clients.RealClient{}

	// Crie uma instância do controller com o cliente
	ctrl := &controller.Controller{Client: client}
	http.HandleFunc("/cep", ctrl.ConsultaTemperatura) // Define a rota e o handler

	port := ":8080"
	fmt.Printf("Starting HTTP server on port %s...\n", port)
	err := http.ListenAndServe(port, nil) // Inicia o servidor na porta 8080
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
