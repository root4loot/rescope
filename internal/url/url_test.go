//
// Written By : Daniel Antonsen (@root4loot)
//
// Distributed Under MIT License
// Copyrights (C) 2019 root4loot
//

package url

import (
	"testing"
)

// TestURIvariants runs various BB URI formats against expression
func TestURIvariants(t *testing.T) {
	leading := []string{"", "https://", "www.", "https://www."}
	var uris []string

	uris = append(uris, "bugbounty.jp/biz/program")
	uris = append(uris, "bugcrowd.com/program")

	uris = append(uris, "federacy.com/program")
	uris = append(uris, "federacy.com/programs/program")
	uris = append(uris, "federacy.com/api/programs/program")
	uris = append(uris, "one.federacy.com/api/programs/program")

	uris = append(uris, "hackenproof.com/biz/program")
	uris = append(uris, "hackerone.com/program")
	uris = append(uris, "intigriti.com/biz/program")
	uris = append(uris, "yeswehack.com/programs/program")
	uris = append(uris, "openbugbounty.org/program")

	for _, s := range leading {
		for _, uri := range uris {
			if len(getBBmatch(s+uri)) == 0 {
				t.Errorf("Failed to match %s ", uri)
			}
		}
	}
}

// TestBBinScope to determine if bugbounty urls in scope are to be found
func TestBBinScope(t *testing.T) {
	scope := `
	foo hackerone.com/program intigriti.com/biz/program\n
	&& yeswehack.com/programs/program\n
	google.com bar
	`
	m := getBBinScope(scope)

	if len(m) != 3 {
		t.Error("Failed to get all BB URIs from scope")
	}
}
