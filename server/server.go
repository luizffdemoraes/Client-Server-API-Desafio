package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/luizffdemoraes/desafio/client-server-api/config"
	"github.com/luizffdemoraes/desafio/client-server-api/schemas"
)

func main() {
	// Initialize Configs
	err := config.Init()
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/cotacao", BuscaCambioHandler)
	http.ListenAndServe(":8080", nil)
}

func BuscaCambioHandler(w http.ResponseWriter, r *http.Request) {

	defer log.Println("Request finalizada")

	if r.URL.Path != "/cotacao" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	exchange, error := findExchange()
	if error != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	error = config.PersistDataBase(exchange)
	if error != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	log.Println("Request processada com sucesso")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(exchange.USDBRL.Bid)

}

func findExchange() (*schemas.UsdBrls, error) {
	log.Println("Request iniciada")

	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
	if err != nil {
		panic(err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	var usdbrl schemas.UsdBrls
	err = json.Unmarshal(body, &usdbrl)
	if err != nil {
		return nil, err
	}
	return &usdbrl, nil
}
