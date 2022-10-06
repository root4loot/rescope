//
// Author: Daniel Antonsen (@root4loot)
// Distributed Under MIT License
//

package bugbountyjp

import (
	"testing"
)

// TestBugbountyjp scrape
func TestBugbountyjp(t *testing.T) {
	Scrape("https://bugbounty.jp/program/57950e9f28e58a1c2cffd2d8")
}
