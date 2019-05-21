//
// Written By : Daniel Antonsen (@root4loot)
//
// Distributed Under MIT License
// Copyrights (C) 2019 root4loot
//

package federacy

import (
	"log"
	"regexp"
	"strings"

	color "github.com/gookit/color"
	req "github.com/root4loot/rescope/internal/bbaas/pkg/request"
)

var include, exclude, scope []string

// Scrape tries to grab scope table for a given program on federacy.com
func Scrape(url string) string {
	re := regexp.MustCompile(`([\w-]+)\/([\w-]+$)`)
	match := re.FindStringSubmatch(url)
	program := match[2]

	// GET program ID
	resp := req.GET("https://one.federacy.com/api/programs/" + program)

	re = regexp.MustCompile(`"id":"([0-9a-z-]+)+`)
	match = re.FindStringSubmatch(resp)

	if match != nil {
		id := match[1]

		// GET scope
		resp = req.GET("https://one.federacy.com/api/program_scopes?program_id=" + id)
		re = regexp.MustCompile(`"identifier":"([^,]+)","in_scope":(true|false)`)
		matches := re.FindAllStringSubmatch(resp, -1)

		include = append(include, "!INCLUDE")
		exclude = append(exclude, "!EXCLUDE")

		// add to slice
		for _, match := range matches {
			if match[2] == "true" {
				include = append(include, match[1])
			} else {
				exclude = append(exclude, match[1])
			}
		}

		// concat slices
		scope = append(include, exclude...)

	} else {
		log.Fatalf("\n%s Failed to read scope from %s. Incorrect program?", color.FgRed.Text("[!]"), url)
	}
	return strings.Join(scope, "\n")
}
