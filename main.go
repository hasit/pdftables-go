package pdftables

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

const (
	endpoint = "https://pdftables.com/api"
)

// PDFTables defines a new client for making API calls.
type PDFTables struct {
	APIKey string
	Host   string
}

// Error defines an error received when making a request to the API.
type Error struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

// Error returns a string representing the error, satisfying the error interface.
func (e Error) Error() string {
	return fmt.Sprintf("PDFTables: %s (%d)", e.Message, e.Code)
}

// NewClient returns a new PDFTables API client which can be used to make requests.
func NewClient(apikey string) *PDFTables {
	return &PDFTables{
		APIKey: apikey,
		Host:   endpoint,
	}
}

// GetBalance gets the number of remaining pages.
// Returns number of pages (integer) and error. Upon error (err != nil), balance returned will be -1 along with appropriate error message.
// Example: examples/getbalance.go
func (p *PDFTables) GetBalance() (int, error) {
	url := fmt.Sprintf("%v", p.Host+"/remaining?key="+p.APIKey)
	balance := -1

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return balance, Error{fmt.Sprintf("Could not create request: %s", err), -1}
	}

	client := &http.Client{}

	res, err := client.Do(req)
	if err != nil {
		return balance, Error{fmt.Sprintf("Failed to make request: %s", err), -1}
	}
	defer res.Body.Close()

	html, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return balance, Error{fmt.Sprintf("Could not read response %s", err), -1}
	}

	if http.StatusOK <= res.StatusCode && res.StatusCode < http.StatusMultipleChoices {
		stringbalance := strings.Trim(string(html), "\n")
		balance, err = strconv.Atoi(stringbalance)
		if err != nil {
			return -1, Error{fmt.Sprintf("Could not convert string into integer: %s", err), -1}
		}
		return balance, nil
	}

	return balance, Error{fmt.Sprintf("%s", html), -1}
}
