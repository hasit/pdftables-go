# pdftables-go

[![Go Report Card](https://goreportcard.com/badge/github.com/hasit/pdftables-go)](https://goreportcard.com/report/github.com/hasit/pdftables-go)
[![GoDoc Reference](https://godoc.org/github.com/hasit/pdftables-go?status.svg)](https://godoc.org/github.com/hasit/pdftables-go)

pdftables-go is an SDK for using [PDFTables' API](https://pdftables.com/pdf-to-excel-api) in your Go application. [PDFTables](https://pdftables.com) automates extraction of data from PDF files.

## Install

``` fish
go get github.com/hasit/pdftables-go
```

## Usage

In order to get started with making your Go application that uses the PDFTables' API, you will need to [register](https://pdftables.com/join) first. Upon successful registration, you will be provided with an API key.

Go ahead and save it somewhere. I propose that you create an environment variable for easy retrieval. Since I use [fish](https://fishshell.com), here is how you can create a global environment variable in fish shell. Finding a way to do the same in other shells is simply a matter of searching a little.

``` fish
set -xg PDFTABLES_APIKEY <your-api-key>
``` 

Make sure you restart your terminal or reload your shell before you go forward.

Don't worry about losing it; you can easily find your API key on the [API Documentation](https://pdftables.com/pdf-to-excel-api) page.

Next step is to create a new Go program and import this library.

``` go
package main

import "github.com/hasit/pdftables-go"
...
```

Now we will create a new client with the recently saved API key.

``` go
...
func main()  {
  apikey := os.Getenv("PDFTABLES_APIKEY")
  p := pdftables.NewClient(apikey)
}
...
```

### Convert PDF File

``` go
...
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
```

### Supported Formats

| Format | pdftables-go const name |
| :----- | :---------------------- |
| CSV | pdftables.FormatCSV |
| XML | pdftables.FormatXML |
| XLSX | pdftables.FormatXLSXSingle, pdftables.FormatXLSXMultiple |

Read the [Formats](https://pdftables.com/pdf-to-excel-api#formats) section for more information on formats and how xlsx-single and xlsx-multiple differ.

### Get Remaining Balance

``` go
...
func main()  {
  apikey := os.Getenv("PDFTABLES_APIKEY")
	p := pdftables.NewClient(apikey)
	balance, err := p.GetBalance()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(balance)
}
```

## Contribute 

Feel free to ask questions, post issues and open pull requests. My only requirement is that you run `gofmt` on your code before you send in a PR.