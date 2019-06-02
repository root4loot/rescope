//
// Written By : Daniel Antonsen (@root4loot)
//
// Distributed Under MIT License
// Copyrights (C) 2019 root4loot
//

package federacy

import (
  "testing"
)

// TestFederacy scrape
func TestFederacy(t *testing.T) {
	Scrape("https://one.federacy.com/api/programs/federacy")
}