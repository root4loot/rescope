//
// Written By : Daniel Antonsen (@root4loot)
//
// Distributed Under MIT License
// Copyrights (C) 2019 root4loot
//

package hackerone

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	doerror "github.com/root4loot/rescope/internal/bbaas/pkg/errors"
	req "github.com/root4loot/rescope/internal/bbaas/pkg/request"
)

var scope []string
var isAssetURL bool

// Scrape tries to grab scope table for a given program on hackerone.com
func Scrape(url string) string {
	re := regexp.MustCompile(`([\w-]+)\/([\w-]+$)`)
	match := re.FindStringSubmatch(url)
	program := match[2]
	endpoint := "https://hackerone.com/graphql?"

	var resp []byte

	// clear global slice
	scope = nil

	// JSON POST data
	var data = []byte(`{  
		"query":"query Team_assets($first_0:Int!) {query {id,...F0}} fragment F0 on Query {_teamAgUhl:team(handle:\"` + program + `\") {handle,_structured_scope_versions2ZWKHQ:structured_scope_versions(archived:false) {max_updated_at},_structured_scopeszxYtW:structured_scopes(first:$first_0,archived:false,eligible_for_submission:true) {edges {node {asset_type, asset_identifier}},pageInfo {hasNextPage,hasPreviousPage}},_structured_scopes3FF98f:structured_scopes(first:$first_0,archived:false,eligible_for_submission:false) {edges {node {asset_type,asset_identifier,},},},},}",
		"variables":{  
		   "first_0":1337
		}
	 }`)

	// POST request to endpoint
	resp, _ = (req.POST(endpoint, data))

	// map interfaces
	m := map[string]interface{}{}

	// unmarshal
	err := json.Unmarshal([]byte(resp), &m)

	// check for errors
	if err != nil {
		doerror.BadJSON()
	}

	// parse map
	parseMap(m)



	return strings.Join(scope, "\n")
}

// parseMap iterates interface maps recursively and calls parseArray to get keys for each map
func parseMap(aMap map[string]interface{}) {
	for key, val := range aMap {
		switch concreteVal := val.(type) {
		case map[string]interface{}:
			// fmt.Println(key)
			if key == "_structured_scopeszxYtW" {
				scope = append(scope, "!INCLUDE")
			} else if key == "_structured_scopes3FF98f" {
				scope = append(scope, "!EXCLUDE")
			}
			parseMap(val.(map[string]interface{}))
		case []interface{}:
			// fmt.Println(key)
			parseArray(val.([]interface{}))
		default:
			if key == "asset_type" && concreteVal == "URL" {
				//fmt.Println("yes")
				isAssetURL = true
			}

			if isAssetURL {
				scope = append(scope, fmt.Sprint(concreteVal))
			}
			// fmt.Println(key, ":", concreteVal)
		}
	}
}

// parseArray does the same thing as parseMap though it iterates the keys instead.
func parseArray(anArray []interface{}) {
	for _, val := range anArray {
		isAssetURL = false
		switch concreteVal := val.(type) {
		case map[string]interface{}:
			// fmt.Println("Index:", i)
			parseMap(val.(map[string]interface{}))
		case []interface{}:
			// fmt.Println("Index:", i)
			parseArray(val.([]interface{}))
		default:
			// fmt.Println("Index", i, ":", concreteVal)
			_ = append(scope, fmt.Sprint(concreteVal))

		}
	}
}
