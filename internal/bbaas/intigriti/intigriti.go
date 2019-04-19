//
// Written By : Daniel Antonsen (@root4loot)
//
// Distributed Under MIT License
// Copyrights (C) 2019 root4loot
//

package intigriti

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	errors "github.com/root4loot/rescope/internal/bbaas/pkg/errors"
	req "github.com/root4loot/rescope/internal/bbaas/pkg/request"
)

var scope []string

// Scrape tries to grab scope table for a given program on intigriti.com
func Scrape(url string) string {
	re := regexp.MustCompile(`([\w-]+)\/([\w-]+$)`)
	match := re.FindStringSubmatch(url)
	business := match[1]
	program := match[2]
	endpoint := "https://api-public.intigriti.com/api/project/" + business + "/" + program

	// clear global slice
	scope = nil

	// GET request to endpoint
	respJSON := req.GET(endpoint)

	// map interfaces
	m := map[string]interface{}{}

	// unmarshal
	err := json.Unmarshal([]byte(respJSON), &m)
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
		case map[string]interface{}:
			parseMap(val.(map[string]interface{}))
		case []interface{}:
			if key == "inScope" {
				scope = append(scope, "!INCLUDE")
				parseArray(val.([]interface{}))
			}
			if key == "outScope" {
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
		//fmt.Println(val)
		switch concreteVal := val.(type) {
		case map[string]interface{}:
			//fmt.Println("Index:", i)
			parseMap(val.(map[string]interface{}))
		case []interface{}:
			//fmt.Println("Index:", i)
			parseArray(val.([]interface{}))
		default:
			//fmt.Println("Index", i, ":", concreteVal)
			_ = fmt.Sprint(concreteVal)
		}
	}
}
