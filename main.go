package pdftables

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// endpoint defines the host address for making calls to PDFTables API.
const endpoint = "https://pdftables.com/api"

const (
	// FormatCSV sets the document to be returned in Comma Separated Values, blank row between pages.
	FormatCSV = "csv"
	// FormatXML sets the document to be returned in HTML <table> tags; <td> tags may have colspan= attributes.
	FormatXML = "xml"
	// FormatXLSXSingle sets the document to be returned in Excel, all PDF pages on one sheet, blank row between pages.
	FormatXLSXSingle = "xlsx-single"
	// FormatXLSXMultiple sets the document to be returned in Excel, one sheet per page of the PDF.
	FormatXLSXMultiple = "xlsx-multiple"
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
	req.Header.Set("Content-Type", "text-plain")

	client := &http.Client{}

	res, err := client.Do(req)
	if err != nil {
		return balance, Error{fmt.Sprintf("Failed to make request: %s", err), -1}
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return balance, Error{fmt.Sprintf("Could not read response %s", err), -1}
	}

	if http.StatusOK <= res.StatusCode && res.StatusCode < http.StatusMultipleChoices {
		stringbalance := strings.Trim(string(body), "\n")
		balance, err = strconv.Atoi(stringbalance)
		if err != nil {
			return -1, Error{fmt.Sprintf("Could not convert string into integer: %s", err), -1}
		}
		return balance, nil
	}

	return balance, Error{fmt.Sprintf("%s", body), -1}
}

// Convert extracts data from PDF file by calling PDFTables API into supported formats.
// Supported formats are CSV, XML, XLSX.
// Returns nil error and creates a file in specified format in the same directory as the PDF file. Upon error (err != nil), no file will be created and appropriate error message will be returned.
// Example: examples/convert.go
// Note: `file` parameter only accepts abosulte file path.
func (p *PDFTables) Convert(file, format string) error {
	url := fmt.Sprintf("%v", p.Host+"?key="+p.APIKey+"&format="+format)
	b := bytes.Buffer{}

	f, err := os.Open(file)
	if err != nil {
		return Error{fmt.Sprintf("Could not open file: %s", err), -1}
	}
	defer f.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("f", filepath.Base(file))
	if err != nil {
		return Error{fmt.Sprintf("Failed to create new form-data header: %s", err), -1}
	}
	_, err = io.Copy(part, f)
	if err != nil {
		return Error{fmt.Sprintf("Unsuccessful copy: %s", err), -1}
	}
	err = writer.Close()
	if err != nil {
		return Error{fmt.Sprintf("Failed to close writer: %s", err), -1}
	}

	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return Error{fmt.Sprintf("Could not create request: %s", err), -1}
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := http.Client{}

	res, err := client.Do(req)
	if err != nil {
		return Error{fmt.Sprintf("Failed to make request: %s", err), -1}
	}
	defer res.Body.Close()

	_, err = b.ReadFrom(res.Body)
	if err != nil {
		return Error{fmt.Sprintf("Failed to read response: %s", err), -1}
	}

	if http.StatusOK <= res.StatusCode && res.StatusCode < http.StatusMultipleChoices {
		err = makeFile(b, file, format)
		if err != nil {
			return Error{fmt.Sprintf("Failed to create file: %s", err), -1}
		}
		return nil
	}

	return Error{fmt.Sprintf("%s", b.String()), -1}
}

// makeFile makes a file in the specified format in the same directory as the PDF file.
// Returns appropriate rror.
func makeFile(b bytes.Buffer, file, format string) error {
	dir := filepath.Dir(file)
	name := filepath.Join(dir, strings.Trim(filepath.Base(file), ".pdf"))
	newfile := ""

	switch format {
	case FormatCSV:
		newfile = name + ".csv"
	case FormatXML:
		newfile = name + ".xml"
	case FormatXLSXSingle, FormatXLSXMultiple:
		newfile = name + ".xlsx"
	default:
		return Error{fmt.Sprintf("Unsupported format"), -1}
	}

	f, err := os.Create(newfile)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(b.Bytes())
	if err != nil {
		return err
	}

	f.Sync()

	return nil
}
