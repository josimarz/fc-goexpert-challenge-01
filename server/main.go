package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/josimarz/fc-goexpert-challenge-01/common"
	_ "github.com/mattn/go-sqlite3"
)

type QuotationDetails struct {
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

type Quotation struct {
	Details QuotationDetails `json:"USDBRL"`
}

var db *sql.DB

func main() {
	setup()
	db.Close()
}

func setup() {
	setupDatabase()
	setupServer()
}

func setupDatabase() {
	_, err := os.Stat("db.sqlite3")
	if err != nil {
		log.Fatal(err.Error())
	}
	db, err = sql.Open("sqlite3", "db.sqlite3")
	if err != nil {
		log.Fatal(err.Error())
	}
	_, err = db.Exec(`
		create table if not exists "quotation" (
			code string,
			codein string,
			name string,
			high string,
			low string,
			varBid string,
			pctChange string,
			bid string,
			ask string,
			timestamp string,
			create_date string
		)
	`)
	if err != nil {
		log.Fatal(err.Error())
	}
}

func setupServer() {
	mux := http.NewServeMux()
	mux.HandleFunc("/cotacao", getQuotation)
	fmt.Println("\033[32mServer listening on port 8080")
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatal(err.Error())
	}
}

func getQuotation(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Millisecond*200)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "GET", "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
	if err != nil {
		log.Fatal(err.Error())
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer res.Body.Close()
	var quotation Quotation
	json.NewDecoder(res.Body).Decode(&quotation)
	select {
	case <-ctx.Done():
		log.Fatal("Unable to query the quotation. Timeout reached.")
	default:
		err = saveQuotation(&quotation.Details)
		if err != nil {
			log.Fatal(err)
		}
		result := &common.ClientQuotation{
			Bid: quotation.Details.Bid,
		}
		w.Header().Add("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)
	}
}

func saveQuotation(details *QuotationDetails) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*10)
	defer cancel()
	stmt, err := db.Prepare(`
		insert into quotation (
			code,
			codein,
			name,
			high,
			low,
			varBid,
			pctChange,
			bid,
			ask,
			timestamp,
			create_date
		) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(
		details.Code,
		details.Codein,
		details.Name,
		details.High,
		details.Low,
		details.VarBid,
		details.PctChange,
		details.Bid,
		details.Ask,
		details.Timestamp,
		details.CreateDate)
	if err != nil {
		return err
	}
	select {
	case <-ctx.Done():
		log.Fatal("Unable to save quotation on database. Timetout reached.")
	default:
		return nil
	}
	return nil
}
