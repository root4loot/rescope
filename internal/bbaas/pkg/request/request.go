//
// Written By : Daniel Antonsen (@root4loot)
//
// Distributed Under MIT License
// Copyrights (C) 2019 root4loot
//

package request

import (
	"bytes"
	"io/ioutil"
	"net/http"

	errors "github.com/root4loot/rescope/internal/bbaas/pkg/errors"
)

// GET returns response body of a given URL as string
func GET(url string) string {

	// request url
	resp, err := http.Get(url)

	// check response
	if err != nil {
		errors.NoResponse(url)
	}

	// close response
	defer resp.Body.Close()

	// response body
	respB, _ := ioutil.ReadAll(resp.Body)

	// response body string
	respBS := string(respB)

	return respBS
}

// POST returns response body of a given URL as bytes
func POST(url string, data []byte) []byte {

	// request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Connection", "close")
	client := &http.Client{}
	resp, err := client.Do(req)

	// check response
	if err != nil {
		errors.NoResponse(url)
	}

	// close response
	defer resp.Body.Close()

	// JSON response body
	respB, _ := ioutil.ReadAll(resp.Body)

	return respB
}
