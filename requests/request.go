package requests

import (
	"bytes"
	"net/http"

	"github.com/tcnksm/go-httpstat"
)

func DoGetRequest(url string) (*http.Response, httpstat.Result, error) {
	req, err := http.NewRequest("GET", url, nil)
	// Create go-httpstat powered context and pass it to http.Request
	var result httpstat.Result
	if err != nil {
		return nil, result, err
	}
	req.Header.Set("Connection", "close")

	ctx := httpstat.WithHTTPStat(req.Context(), &result)
	req = req.WithContext(ctx)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, result, err
	}
	return res, result, nil
}

func DoPostRequest(url string, jsonData []byte) (*http.Response, httpstat.Result, error) {
	// Create go-httpstat powered context and pass it to http.Request
	var result httpstat.Result
	req, err := http.NewRequest("POST", url, bytes.NewReader(jsonData))
	if err != nil {
		return nil, result, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Connection", "close")

	ctx := httpstat.WithHTTPStat(req.Context(), &result)
	req = req.WithContext(ctx)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, result, err
	}
	return res, result, nil
}
