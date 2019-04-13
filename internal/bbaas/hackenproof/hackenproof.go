//
// Written By : Daniel Antonsen (@root4loot)
//
// Distributed Under MIT License
// Copyrights (C) 2019 root4loot
//

package hackenproof

import (
	"regexp"
	"strings"

	"github.com/antchfx/htmlquery"

	errors "github.com/root4loot/rescope/internal/bbaas/pkg/errors"
	req "github.com/root4loot/rescope/internal/bbaas/pkg/request"
)

// Scrape returns a string containing scope that was scraped from the given program on hackenproof.com
func Scrape(url string) string {
	var includes, excludes, scope []string

	re := regexp.MustCompile(`([\w-]+)\/([\w-]+$)`)
	match := re.FindStringSubmatch(url)
	business := match[1]
	program := match[2]

	endpoint := "https://hackenproof.com/" + business + "/" + program

	// GET request to endpoint
	respBody := req.GET(endpoint)

	// parse response body to xQuery doc
	doc, _ := htmlquery.Parse(strings.NewReader(respBody))

	// xQuery to grab in-scope and out-of-scope tables
	inScope := htmlquery.Find(doc, "//h4[contains(text(), 'In Scope')]//following::table//h5")
	outScope := htmlquery.Find(doc, "//h4[contains(text(), 'Out of scope')]//following::table//h5")

	// get in-scope / out-scope content
	if inScope != nil {
		includes = append(scope, "!INCLUDE")
		for _, item := range inScope {
			includes = append(includes, htmlquery.InnerText(item))
		}
		if outScope != nil {
			excludes = append(scope, "!EXCLUDE")
			for _, item := range outScope {
				excludes = append(excludes, htmlquery.InnerText(item))
			}
		}
	} else {
		errors.NoScope(url)
	}

	// remove duplicates as inScope contains outScope as well
	for i, v := range includes {
		for _, v2 := range excludes {
			if v == v2 {
				includes[i] = ""
			}
		}
	}

	for _, v := range includes {
		scope = append(scope, v)
	}

	for _, v := range excludes {
		scope = append(scope, v)
	}

	return strings.Join(scope, "\n")
}
