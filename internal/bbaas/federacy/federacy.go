//
// Author: Daniel Antonsen (@root4loot)
// Distributed Under MIT License
//

package federacy

import (
	"regexp"
	"strings"

	errors "github.com/root4loot/rescope/internal/bbaas/pkg/errors"
	req "github.com/root4loot/rescope/internal/bbaas/pkg/request"
)

var hasExclude bool = false

// Scrape attempts to grab scope table for a given program on federacy.com
func Scrape(url string) string {
	var scope []string
	re := regexp.MustCompile(`([\w-]+)\/([\w-]+$)`)
	match := re.FindStringSubmatch(url)
	program := strings.ToLower(match[2])

	// GET program ID
	resp, status := req.GET("www.federacy.com/" + program)
	if status != 200 {
		errors.BadStatusCode(url, status)
	}

	re = regexp.MustCompile(`identifier:"(.*?)",in_scope:(a|d)`)
	matches := re.FindAllStringSubmatch(resp, -1)

	if matches != nil {
		scope = append(scope, "!INCLUDE")
		for _, match := range matches {
			if match[2] == "d" {
				scope = append(scope, match[1])
			} else if match[2] == "a" {
				hasExclude = true
			}
		}
		if hasExclude {
			scope = append(scope, "!EXCLUDE")
			for _, match := range matches {
				if match[2] == "a" {
					scope = append(scope, match[1])
				}
			}
		}
	}
	return strings.Join(scope, "\n")
}
