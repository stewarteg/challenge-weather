package main

import (
	"fmt"
	"net/http"

	"github.com/stewarteg/goexpert/challenge-weather/controller"
)

func main() {
	http.HandleFunc("/cep", controller.ConsultaTemperatura) // Define a rota e o handler

	port := ":8080"
	fmt.Printf("Starting HTTP server on port %s...\n", port)
	err := http.ListenAndServe(port, nil) // Inicia o servidor na porta 8080
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
