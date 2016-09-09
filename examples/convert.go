package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/hasit/pdftables-go"
)

func main() {
	apikey := os.Getenv("PDFTABLES_APIKEY")
	p := pdftables.NewClient(apikey)
	path, err := filepath.Abs("examples/test.pdf")
	if err != nil {
		log.Fatal(err)
	}
	err = p.Convert(path, pdftables.FormatCSV)
	if err != nil {
		log.Fatal(err)
	}
}
