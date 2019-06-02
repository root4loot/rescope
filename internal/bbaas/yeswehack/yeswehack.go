//
// Written By : Daniel Antonsen (@root4loot)
//
// Distributed Under MIT License
// Copyrights (C) 2019 root4loot
//

package yeswehack

import (
	"encoding/json"
	"regexp"
	"strings"
	"fmt"

	errors "github.com/root4loot/rescope/internal/bbaas/pkg/errors"
	req "github.com/root4loot/rescope/internal/bbaas/pkg/request"
)

var scope []string

// Scrape tries to grab scope table for a given program on yeswehack.com
func Scrape(url string) string {
	re := regexp.MustCompile(`([\w-]+)\/([\w-]+$)`)
	match := re.FindStringSubmatch(url)
	program := match[2]
	endpoint := "https://api.yeswehack.com/programs/" + program

	// clear global slice
	scope = nil

	// GET request to endpoint
	resp, status := req.GET(endpoint)

	// check bad status code
	if status != 200 {
		errors.BadStatusCode(url, status)
	}

	// map interfaces
	m := map[string]interface{}{}

	// unmarshal
	err := json.Unmarshal([]byte(resp), &m)
	_ = err

	// check for errors
	if err != nil {
		errors.BadJSON()
	}

	// parse map
	parseMap(m)

	// throw error if scope array is empty
	if scope == nil {
		errors.NoScope(url)
	}

	return strings.Join(scope, "\n")
}

// parseMap iterates interface maps recursively and calls parseArray to get keys for each map
func parseMap(aMap map[string]interface{}) {
	for key, val := range aMap {
		switch concreteVal := val.(type) {
		case []interface{}:
			//fmt.Println(key)
			if key == "scopes" {
				scope = append(scope, "!INCLUDE")
				scope = append(scope, fmt.Sprint(concreteVal))
				parseArray(val.([]interface{}))
			} else if key == "out_of_scope" {
				scope = append(scope, "!EXCLUDE")
				scope = append(scope, fmt.Sprint(concreteVal))
				parseArray(val.([]interface{}))
			}
		default:
			// fmt.Println(key, ":", concreteVal)
			if key == "content" {
				scope = append(scope, fmt.Sprint(concreteVal))
			}
		}
	}
}

// parseArray does the same thing as parseMap though it iterates the keys instead.
func parseArray(anArray []interface{}) {
	for _, val := range anArray {
		switch concreteVal := val.(type) {
		case map[string]interface{}:
			//fmt.Println("Index:", i)
			parseMap(val.(map[string]interface{}))
		default:
			_ = fmt.Sprint(concreteVal)

		}
	}
}
