//
// Author: Daniel Antonsen (@root4loot)
// Distributed Under MIT License
//

package federacy

import (
	"testing"
)

// TestFederacy scrape
func TestFederacy(t *testing.T) {
	Scrape("https://one.federacy.com/api/programs/federacy")
}
