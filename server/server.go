package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type UsdBrls struct {
	USDBRL Usdbrl
}

type Usdbrl struct {
	ID         int    `gorm:"primaryKey"`
	Code       string `json:"code"`
	Codein     string `json:"codein"`
	Name       string `json:"name"`
	High       string `json:"high"`
	Low        string `json:"low"`
	VarBid     string `json:"varBid"`
	PctChange  string `json:"pctChange"`
	Bid        string `json:"bid"`
	Ask        string `json:"ask"`
	Timestamp  string `json:"timestamp"`
	CreateDate string `json:"create_date"`
}

func main() {
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

	error = persistDataBase(exchange)
	if error != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	log.Println("Request processada com sucesso")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(exchange.USDBRL.Bid)

}

func findExchange() (*UsdBrls, error) {
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

	var usdbrl UsdBrls
	err = json.Unmarshal(body, &usdbrl)
	if err != nil {
		return nil, err
	}
	return &usdbrl, nil
}

func persistDataBase(exchange *UsdBrls) error {
	dsn := "root:root@tcp(localhost:3306)/goexpert?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	db.AutoMigrate(&Usdbrl{})

	// CREATE
	db.WithContext(ctx).Create(&Usdbrl{
		Code:       exchange.USDBRL.Code,
		Codein:     exchange.USDBRL.Codein,
		Name:       exchange.USDBRL.Name,
		High:       exchange.USDBRL.High,
		Low:        exchange.USDBRL.Low,
		VarBid:     exchange.USDBRL.VarBid,
		PctChange:  exchange.USDBRL.PctChange,
		Bid:        exchange.USDBRL.Bid,
		Ask:        exchange.USDBRL.Ask,
		Timestamp:  exchange.USDBRL.Timestamp,
		CreateDate: exchange.USDBRL.CreateDate,
	})

	log.Println("Persistencia concluida.")

	return nil
}
