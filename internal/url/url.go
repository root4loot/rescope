package url

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	color "github.com/gookit/color"
	bugbountyjp "github.com/root4loot/rescope/internal/bbaas/bugbounty.jp"
	bugcrowd "github.com/root4loot/rescope/internal/bbaas/bugcrowd"
	federacy "github.com/root4loot/rescope/internal/bbaas/federacy"
	hackenproof "github.com/root4loot/rescope/internal/bbaas/hackenproof"
	hackerone "github.com/root4loot/rescope/internal/bbaas/hackerone"
	intigriti "github.com/root4loot/rescope/internal/bbaas/intigriti"
	openbugbounty "github.com/root4loot/rescope/internal/bbaas/openbugbounty"
	yeswehack "github.com/root4loot/rescope/internal/bbaas/yeswehack"
)

// BBaas identifies supported bounty programs and calls scrape functions to grab scopes
func BBaas(urls, scopes, source []string) ([]string, []string, bool) {
	var matches [][]string
	var foundInScope bool

	if len(urls) > 0 {
		for _, url := range urls {
			match := getBBmatch(url)
			if match != nil {
				matches = append(matches, match)
			} else {
				log.Fatalf("%s Invalid bug bounty URI: %s\n", color.FgRed.Text("[!]"), url)
			}
		}
	}

	for i, scope := range scopes {
		matched := getBBinScope(scope)
		if len(matched) != 0 {
			foundInScope = true
			for _, match := range matched {
				fmt.Printf("%s Found BBaaS URI (%s) in %s\n", color.FgYellow.Text("[-]"), match[0], source[i])
				scopes[i] = strings.Replace(scopes[i], match[0], "", -1)
				matches = append(matches, match)
			}
		}
	}

	m := map[string]func(string) string{
		"hackerone.com":     hackerone.Scrape,
		"bugcrowd.com":      bugcrowd.Scrape,
		"hackenproof.com":   hackenproof.Scrape,
		"intigriti.com":     intigriti.Scrape,
		"openbugbounty.org": openbugbounty.Scrape,
		"yeswehack.com":     yeswehack.Scrape,
		"bugbounty.jp":      bugbountyjp.Scrape,
		"federacy.com":      federacy.Scrape,
	}

	for _, match := range matches {
		scopes = append(scopes, m[match[4]](match[0]))
		source = append(source, match[0])
	}
	return scopes, source, foundInScope
}

// getBBinScope returns matched BB URIs from scope
func getBBinScope(s string) [][]string {
	lines := strings.FieldsFunc(s, split)
	var matches [][]string

	for _, line := range lines {
		match := getBBmatch(line)
		if match != nil {
			matches = append(matches, match)
		}
	}
	return matches
}

// getBBmatch returns slice containing submatches from expression that checks for valid URIs
func getBBmatch(s string) []string {
	re := regexp.MustCompile(`((https?:\/\/)?(www\.)?(hackerone\.com|bugcrowd\.com|hackenproof\.com|intigriti\.com\/([\w-_\/]+)|openbugbounty\.org|yeswehack\.com|bugbounty\.jp|(one\.)?federacy\.com)(\/[\w-_]+)?\/[\w-_]+)\/?`)
	var match []string

	if re.MatchString(s) {
		match = re.FindStringSubmatch(s)
	}
	return match
}

// Used with FieldsFunc to split on multiple delimiters
func split(r rune) bool {
	return r == ' ' || r == '\n'
}
