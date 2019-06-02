//
// Written By : Daniel Antonsen (@root4loot)
//
// Distributed Under MIT License
// Copyrights (C) 2019 root4loot
//

package hackerone

import (
  "testing"
)

// TestHackerone scrape
func TestHackerone(t *testing.T) {
	Scrape("hackerone.com/security")
}