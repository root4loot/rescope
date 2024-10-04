<img src="logo.png" alt="Logo" width="900"/>

[![HackerOne](https://github.com/root4loot/rescope/actions/workflows/test-hackerone.yml/badge.svg?branch=main)](https://github.com/root4loot/rescope/actions/workflows/test-hackerone.yml)
[![Bugcrowd](https://github.com/root4loot/rescope/actions/workflows/test-bugcrowd.yml/badge.svg?branch=main)](https://github.com/root4loot/rescope/actions/workflows/test-bugcrowd.yml)
[![Intigriti](https://github.com/root4loot/rescope/actions/workflows/test-intigriti.yml/badge.svg?branch=main)](https://github.com/root4loot/rescope/actions/workflows/test-intigriti.yml)
[![YesWeHack](https://github.com/root4loot/rescope/actions/workflows/test-yeswehack.yml/badge.svg?branch=main)](https://github.com/root4loot/rescope/actions/workflows/test-yeswehack.yml)
[![HackenProof](https://github.com/root4loot/rescope/actions/workflows/test-hackenproof.yml/badge.svg?branch=main)](https://github.com/root4loot/rescope/actions/workflows/test-hackenproof.yml)
[![Test CLI](https://github.com/root4loot/rescope/actions/workflows/test-cli.yml/badge.svg?branch=main)](https://github.com/root4loot/rescope/actions/workflows/test-cli.yml)
![Twitter Follow](https://img.shields.io/twitter/follow/danielantonsen.svg?style=dark)


Use this tool to fetch public/private scopes from bugbounty programs and output them in various formats. 

## Supported platforms

- [HackerOne](https://hackerone.com)
- [Bugcrowd](https://bugcrowd.com)
- [Intigriti](https://www.intigriti.com)
- [YesWeHack](https://yeswehack.com)
- [HackenProof](https://hackenproof.com)

## Installation

Requires Go 1.23 or later.

```bash
go install github.com/root4loot/rescope2/cmd/rescope@latest
```

## Docker

```bash
docker build -t rescope .
docker run --rm -it rescope [options] [<BugBountyURL>...]
```

## Usage

```
Usage:
  rescope [options] [<BugBountyURL>...] [-iL <file>] [-eL <file>]

INPUT:
  -iL, --include-list         file containing list of URLs or custom in-scope definitions (newline separated)
  -eL, --exclude-list         file containing list of URLs or custom out-of-scope definitions (newline separated)

OUTPUT:
  -oF, --output-file          output to given file (default: stdout)

OUTPUT FORMAT:
  -oT, --output-text          output simple text (default)
  -oB, --output-burp          output Burp Suite Scope (JSON)
  -oZ, --output-zap           output ZAP Scope (XML)
  -oJ, --output-json          output JSON
  -oJL, --output-json-lines   output JSON lines

OUTPUT FILTER:
  --filter-expand-ip-ranges   output individual IPs instead of IP ranges / CIDRs

AUTHORIZATION:
  --auth-bugcrowd             bugcrowd secret    (_bugcrowd_session=cookie.value) [Optional]
  --auth-hackenproof          hackenproof secret (_hackenproof_session=cookie.value) [Optional]
  --auth-hackerone            hackerone secret   (Authorization bearer token) [Optional]
  --auth-yeswehack            yeswehack secret   (Authorization bearer token) [Optional]
  --auth-intigriti            intigriti secret   (see https://app.intigriti.com/researcher/personal-access-tokens) [Optional]

GENERAL:
  -c, --concurrency           maximum number of concurrent requests (default: 5)
      --debug                 enable debug mode
      --version               display version
```

### Examples

#### Basic Usage

```bash
rescope https://hackerone.com/security https://bugcrowd.com/tesla
```

```bash
rescope --output-file burp_scope.json --output-burp https://hackerone.com/security https://bugcrowd.com/tesla
```

#### Custom includes / excludes

Note that `--include-list` file may also contain bug bounty URLs.

```bash
rescope -iL include.txt -eL exclude.txt
```

#### Piping to rescope

You may also pipe a list of bug bounty URLs directly to rescope:

```bash
cat urls.txt | rescope
```

## As a library

```go
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
```

## Importing to Burp Suite and OWASP ZAP

### Burp Suite

1. Select Settings -> Project -> Scope
2. Click the ⚙︎ icon below the "Target Scope" title and choose "Load settings"
3. Select Burp JSON file exported from rescope

### OWASP ZAP

1. Select File -> Import Context
2. Select the ZAP XML file exported from rescope

## Contributing

Contributions are welcome. To contribute, fork the repository, create a new branch, make your changes, and send a pull request.
