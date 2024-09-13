package main

import (
	"fmt"
	"log"

	"github.com/root4loot/rescope2/pkg/rescope"
)

func main() {
	opts := rescope.DefaultOptions()

	opts.AuthHackerOne = "your_hackerone_token" // Optional
	opts.AuthIntigriti = "your_intigriti_token" // Optional

	bugBountyURLs := []string{
		"https://hackerone.com/security",
		"https://bugcrowd.com/tesla",
	}

	for _, url := range bugBountyURLs {
		result, err := rescope.Run(url, opts)
		if err != nil {
			log.Printf("Failed to run rescope for URL %s: %v", url, err)
			continue
		}

		fmt.Printf("Results for %s:\n", url)
		fmt.Printf("In-Scope: %v\n", result.InScope)
		fmt.Printf("Out-Scope: %v\n", result.OutScope)
	}
}
