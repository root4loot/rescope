package url

import (
	"fmt"
	"log"
	"regexp"

	bugbountyjp "github.com/root4loot/rescope/internal/bbaas/bugbounty.jp"
	bugcrowd "github.com/root4loot/rescope/internal/bbaas/bugcrowd"
	federacy "github.com/root4loot/rescope/internal/bbaas/federacy"
	hackenproof "github.com/root4loot/rescope/internal/bbaas/hackenproof"
	hackerone "github.com/root4loot/rescope/internal/bbaas/hackerone"
	intigriti "github.com/root4loot/rescope/internal/bbaas/intigriti"
	openbugbounty "github.com/root4loot/rescope/internal/bbaas/openbugbounty"
	yeswehack "github.com/root4loot/rescope/internal/bbaas/yeswehack"

	"github.com/gookit/color"
)

// BBaas identifies supported bounty programs and calls Scrape functions to get scopes.
func BBaas(urls, scopes, source []string) ([]string, []string, bool) {
	// indicates whether BBaaS url was found
	var bbaas bool
	// Move bounty URLs from infile to a.URLs
	for i, scope := range scopes {
		re := regexp.MustCompile(`((https?:\/\/)?(www\.)?(hackerone\.com|bugcrowd\.com|hackenproof\.com|intigriti\.com\/([\w-_\/]+)|openbugbounty\.org|yeswehack\.com|bugbounty\.jp|(one\.)?federacy\.com)(\/[\w-_]+)?\/[\w-_]+)\/?`)

		// get all bb URLs from scope
		bountyuris := re.FindAllString(scope, -1)

		// add them to the list of bb URLs
		for _, v := range bountyuris {
			bbaas = true
			fmt.Printf("%s Identified BBaaS program (%s) in %s\n", color.FgYellow.Text("[-]"), v, source[i])
			urls = append(urls, v)
		}

		// remove from infile
		scopes[i] = re.ReplaceAllString(scope, "")
	}

	// Get scope from bugbounty URL(s)
	if urls != nil {
		re := regexp.MustCompile(`^(https?:\/\/)?(www\.)?(([a-z]+\.)?[a-zA-Z0-9-]+\.[a-z]+)\/([a-zA-Z0-9-_]+)(\/[a-zA-Z0-9-_\/]+)?`)
		// relevant groups
		// 1. [www.example.com/biz/program]
		// 3. [www.[example.com]/biz/program]

		for _, v := range urls {
			r := re.FindStringSubmatch(v)
			var url, host string

			if r != nil {
				url = r[0]
				host = r[3]
				// program = r[5]
			} else {
				log.Fatalf("%s Invalid bug bounty URL: %s\n", color.FgRed.Text("[!]"), v)
			}

			// Scrape scopes from BB program tables
			if host == "hackerone.com" {
				scopes = append(scopes, hackerone.Scrape(url))
				source = append(source, url)
			} else if host == "bugcrowd.com" {
				scopes = append(scopes, bugcrowd.Scrape(url))
				source = append(source, url)
			} else if host == "hackenproof.com" {
				scopes = append(scopes, hackenproof.Scrape(url))
				source = append(source, url)
			} else if host == "intigriti.com" {
				scopes = append(scopes, intigriti.Scrape(url))
				source = append(source, url)
			} else if host == "openbugbounty.org" {
				scopes = append(scopes, openbugbounty.Scrape(url))
				source = append(source, url)
			} else if host == "yeswehack.com" {
				scopes = append(scopes, yeswehack.Scrape(url))
				source = append(source, url)
			} else if host == "bugbounty.jp" {
				scopes = append(scopes, bugbountyjp.Scrape(url))
				source = append(source, url)
			} else if host == "federacy.com" {
				scopes = append(scopes, federacy.Scrape(url))
				source = append(source, url)
			} else {
				log.Fatalf("%s Unsupported bug bounty program: %s\n", color.FgRed.Text("[!]"), host)
			}
		}
	}
	return scopes, source, bbaas
}
