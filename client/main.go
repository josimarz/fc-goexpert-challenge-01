package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"text/template"
	"time"

	"github.com/josimarz/fc-goexpert-challenge-01/common"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*300)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/cotacao", nil)
	if err != nil {
		log.Fatal(err.Error())
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer res.Body.Close()
	var clientQuotation common.ClientQuotation
	json.NewDecoder(res.Body).Decode(&clientQuotation)
	select {
	case <-ctx.Done():
		log.Fatal("Unable to query the quotation. Timeout rechead.")
	default:
		err := saveQuotation(&clientQuotation)
		if err != nil {
			log.Fatal(err.Error())
		}
	}
}

func saveQuotation(clientQuotation *common.ClientQuotation) error {
	file, err := os.Create("cotacao.txt")
	if err != nil {
		return err
	}
	defer file.Close()
	tmpl, err := template.New("output").Parse("Dólar: {{.Bid}}")
	if err != nil {
		return err
	}
	err = tmpl.Execute(file, clientQuotation)
	if err != nil {
		return err
	}
	return nil
}
