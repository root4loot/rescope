//
// Author: Daniel Antonsen (@root4loot)
// Distributed Under MIT License
//

package hackerone

import (
	"testing"
)

// TestHackerone scrape
func TestHackerone(t *testing.T) {
	Scrape("hackerone.com/security")
}
