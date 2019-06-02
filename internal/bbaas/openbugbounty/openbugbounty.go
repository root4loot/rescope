//
// Written By : Daniel Antonsen (@root4loot)
//
// Distributed Under MIT License
// Copyrights (C) 2019 root4loot
//

package openbugbounty

import (
	"regexp"
	"strings"

	"github.com/antchfx/htmlquery"

	errors "github.com/root4loot/rescope/internal/bbaas/pkg/errors"
	request "github.com/root4loot/rescope/internal/bbaas/pkg/request"
)

// Scrape tries to grab scope table for a given program on openbugbounty.com
func Scrape(url string) string {
	re := regexp.MustCompile(`([\w-]+)\/([\w-]+$)`)
	match := re.FindStringSubmatch(url)
	program := match[2]
	endpoint := "https://www.openbugbounty.org/bugbounty/" + program
	var scope []string

	// request endpoint
	respBody, status := request.GET(endpoint)

	// check bad status code
	if status != 200 {
		errors.BadStatusCode(url, status)
	}

	// parse response body to xQuery doc
	doc, _ := htmlquery.Parse(strings.NewReader(respBody))

	// xQuery to grab in-scope items
	inScope := htmlquery.Find(doc, "//h3[contains(text(), 'Bug Bounty Scope')]/following-sibling::table//td[text()]")

	// append to scope
	if inScope != nil {
		for _, item := range inScope {
			scope = append(scope, htmlquery.InnerText(item))
		}
	} else {
		errors.NoScope(url)
	}

	return strings.Join(scope, "\n")
}
