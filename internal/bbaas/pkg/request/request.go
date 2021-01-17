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

	doerror "github.com/root4loot/rescope/internal/bbaas/pkg/errors"
)

// GET returns response body and status code for a given URL
func GET(url string) (string, int) {

	// request url
	resp, err := http.Get(url)
	respS := resp.StatusCode

	// check response
	if err != nil {
		doerror.NoResponse(url)
	}

	// close response
	defer resp.Body.Close()

	// response body
	respB, _ := ioutil.ReadAll(resp.Body)

	// response body string
	respBS := string(respB)

	return respBS, respS
}

// POST returns response body and status code for a given URL
func POST(url string, data []byte) ([]byte, int) {

	// request
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Connection", "close")
	client := &http.Client{}
	resp, err := client.Do(req)
	respS := resp.StatusCode

	// check response
	if err != nil {
		doerror.NoResponse(url)
	}

	// close response
	defer resp.Body.Close()

	// JSON response body
	respB, _ := ioutil.ReadAll(resp.Body)

	return respB, respS
}
