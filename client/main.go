package main

import (
	"context"
	"io"
	"net/http"
	"os"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/cotacao", nil)
	if err != nil {
		panic(err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	fileName := "cotacao.txt"
	f, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	_, err = io.WriteString(f, "Dólar: "+string(data))
	if err != nil {
		panic(err)
	}
}
