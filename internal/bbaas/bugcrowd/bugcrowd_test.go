//
// Author: Daniel Antonsen (@root4loot)
// Distributed Under MIT License
//

package bugcrowd

import (
	"testing"
)

// TestBugcrowd scrape
func TestBugcrowd(t *testing.T) {
	Scrape("https://bugcrowd.com/bugcrowd")
}
