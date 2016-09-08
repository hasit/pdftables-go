package main

import (
	"fmt"
	"log"
	"os"

	"github.com/hasit/pdftables-go"
)

func main() {
	apikey := os.Getenv("PDFTABLES_APIKEY")
	p := pdftables.NewClient(apikey)
	balance, err := p.GetBalance()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(balance)
}
