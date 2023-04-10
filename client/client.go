package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

func main() {
	// Criação de Arquivo
	file, err := os.Create("cotacao.txt")
	if err != nil {
		panic(err)
	}

	defer file.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080", nil)
	if err != nil {
		panic(err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)

	if err != nil {
		panic(err)
	}

	fmt.Println(b)

	io.Copy(os.Stdout, resp.Body)

	// Escrita
	_, err = file.WriteString(fmt.Sprintf("Dólar: %s", b))
	if err != nil {
		panic(err)
	}

	file.Close()
}
