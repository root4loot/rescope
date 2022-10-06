//
// Author: Daniel Antonsen (@root4loot)
// Distributed Under MIT License
//

package yeswehack

import (
	"testing"
)

// TestHackerone scrape
func TestYeswehack(t *testing.T) {
	Scrape("https://yeswehack.com/programs/yes-we-hack")
}
