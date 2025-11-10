package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type BrasilAPIResposta struct {
	Cep          string `json:"cep"`
	State        string `json:"state"`
	City         string `json:"city"`
	Neighborhood string `json:"neighborhood"`
	Street       string `json:"street"`
}

type ViaCEPResposta struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	Uf          string `json:"uf"`
	Erro        bool   `json:"erro"`
}

type Resultado struct {
	Dados interface{}
	API   string
	Erro  error
}

func buscarBrasilAPI(ctx context.Context, cep string, resultado chan Resultado) {
	url := fmt.Sprintf("https://brasilapi.com.br/api/cep/v1/%s", cep)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		resultado <- Resultado{Erro: err, API: "API BrasilAPI"}
		return
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		resultado <- Resultado{Erro: err, API: "API BrasilAPI"}
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		resultado <- Resultado{Erro: err, API: "API BrasilAPI"}
		return
	}

	var data BrasilAPIResposta
	err = json.Unmarshal(body, &data)
	if err != nil {
		resultado <- Resultado{Erro: err, API: "API BrasilAPI"}
		return
	}

	resultado <- Resultado{Dados: data, API: "API BrasilAPI"}
}

func buscarViaCEP(ctx context.Context, cep string, resultado chan Resultado) {
	url := fmt.Sprintf("http://viacep.com.br/ws/%s/json/", cep)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		resultado <- Resultado{Erro: err, API: "API ViaCEP"}
		return
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		resultado <- Resultado{Erro: err, API: "API ViaCEP"}
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		resultado <- Resultado{Erro: err, API: "API ViaCEP"}
		return
	}

	var data ViaCEPResposta
	err = json.Unmarshal(body, &data)
	if err != nil {
		resultado <- Resultado{Erro: err, API: "API ViaCEP"}
		return
	}

	if data.Erro {
		resultado <- Resultado{Erro: fmt.Errorf("CEP não encontrado"), API: "API ViaCEP"}
		return
	}

	resultado <- Resultado{Dados: data, API: "API ViaCEP"}
}

func main() {
	cep := "93037-220"

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	resultado := make(chan Resultado, 2)

	go buscarBrasilAPI(ctx, cep, resultado)
	go buscarViaCEP(ctx, cep, resultado)

	select {
	case result := <-resultado:
		if result.Erro != nil {
			fmt.Printf("Erro na requisição %s: %v\n", result.API, result.Erro)
			// Tenta a segunda resposta
			result = <-resultado
			if result.Erro != nil {
				fmt.Printf("Erro na requisição %s: %v\n", result.API, result.Erro)
			} else {
				imprimirResultado(result)
			}
		} else {
			imprimirResultado(result)
		}
	case <-ctx.Done():
		fmt.Println("Erro: Timeout de 1 segundo excedido")
	}
}

func imprimirResultado(result Resultado) {
	fmt.Printf("\n✓ Resposta recebida de: %s\n\n", result.API)

	switch data := result.Dados.(type) {
	case BrasilAPIResposta:
		fmt.Printf("CEP: %s\n", data.Cep)
		fmt.Printf("Rua: %s\n", data.Street)
		fmt.Printf("Bairro: %s\n", data.Neighborhood)
		fmt.Printf("Cidade: %s\n", data.City)
		fmt.Printf("Estado: %s\n", data.State)
	case ViaCEPResposta:
		fmt.Printf("CEP: %s\n", data.Cep)
		fmt.Printf("Rua: %s\n", data.Logradouro)
		fmt.Printf("Bairro: %s\n", data.Bairro)
		fmt.Printf("Cidade: %s\n", data.Localidade)
		fmt.Printf("Estado: %s\n", data.Uf)
	}
}
