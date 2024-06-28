package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Cotacao struct {
	USDBRL USDBRL `json:"USDBRL"`
}

type USDBRL struct {
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

const dataBasePath = "./cotacao.db"

func main() {
	// DB SQLite
	db, err := sql.Open("sqlite3", dataBasePath)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS cotacao (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		code TEXT,
		codein TEXT,
		name TEXT,
		high TEXT,
		low TEXT,
		var_bid TEXT,
		pct_change TEXT,
		bid TEXT,
		ask TEXT,
		timestamp TEXT,
		create_date TEXT
	)`)
	if err != nil {
		panic(err)
	}

	// Server
	http.HandleFunc("/cotacao", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		ctx := r.Context()
		ctx, cancel := context.WithTimeout(ctx, 200*time.Millisecond)
		defer cancel()

		cotacao, error := BuscarCotacao(ctx)

		if error != nil {
			log.Println(error)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		err = insertCotacao(ctx, db, cotacao)
		if err != nil {
			log.Println(error)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		json.NewEncoder(w).Encode(cotacao.Bid)
	})

	http.ListenAndServe(":8080", nil)
}

func BuscarCotacao(ctx context.Context) (*USDBRL, error) {
	req, error := http.NewRequestWithContext(ctx, "GET", "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
	if error != nil {
		return nil, error
	}

	resp, error := http.DefaultClient.Do(req)
	if error != nil {
		return nil, error
	}
	defer resp.Body.Close()

	body, error := io.ReadAll(resp.Body)
	if error != nil {
		return nil, error
	}

	var c Cotacao
	error = json.Unmarshal(body, &c)
	if error != nil {
		return nil, error
	}

	return &c.USDBRL, nil
}

func insertCotacao(ctx context.Context, db *sql.DB, cotacao *USDBRL) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Millisecond)
	defer cancel()

	stmt, error := db.Prepare("insert into cotacao(code, codein, name, high, low, var_bid, pct_change, bid, ask, timestamp, create_date) values (?, ?, ?,  ?, ?, ?, ?, ?, ?, ?, ?)")
	if error != nil {
		return error
	}
	defer stmt.Close()

	_, error = stmt.ExecContext(ctx, cotacao.Code, cotacao.Codein, cotacao.Name, cotacao.High, cotacao.Low, cotacao.VarBid, cotacao.PctChange, cotacao.Bid, cotacao.Ask, cotacao.Timestamp, cotacao.CreateDate)
	if error != nil {
		return error
	}

	return nil
}
