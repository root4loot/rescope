//
// Author: Daniel Antonsen (@root4loot)
// Distributed Under MIT License
//

package request

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	doerror "github.com/root4loot/rescope/internal/bbaas/pkg/errors"
)

// GET returns response body and status code for a given URL
func GET(url string) (string, int) {

	// prepend https:// unless prefixed
	if !strings.HasPrefix(url, "https://") {
		url = "https://" + url
	}

	client := http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}

	// headers
	req.Header = http.Header{
		"User-Agent": []string{"Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:97.0) Gecko/20100101 Firefox/97.0"},
		"Accept":     []string{"*/*"},
	}

	// do req
	resp, err := client.Do(req)
	if err != nil {
		doerror.NoResponse(url)
	}

	// status code
	respS := resp.StatusCode

	// body
	respB, _ := ioutil.ReadAll(resp.Body)

	// body string
	respBS := string(respB)

	// close response
	defer resp.Body.Close()

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
