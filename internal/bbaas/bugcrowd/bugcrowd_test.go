//
// Written By : Daniel Antonsen (@root4loot)
//
// Distributed Under MIT License
// Copyrights (C) 2019 root4loot
//

package bugcrowd

import (
  "testing"
)

// TestBugcrowd scrape
func TestBugcrowd(t *testing.T) {
	Scrape("https://bugcrowd.com/bugcrowd")
}