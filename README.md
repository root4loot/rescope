<div align="center">
  <img src="logo.png" alt="Logo" width="970"/>
  <a href="https://img.shields.io/twitter/follow/danielantonsen"> </a>
  <p>Use this tool to fetch public/private scopes from bug bounty programs and output them in various formats.</p>
  <img src="https://img.shields.io/twitter/follow/danielantonsen" alt="Twitter Follow"/>
</div>

<hr>

<div align="center" style="padding: 20px; margin: 20px;">
  <a href="https://github.com/root4loot/rescope/actions/workflows/test-hackerone.yml">
    <img src="https://img.shields.io/badge/HackerOne-passing-brightgreen?style=flat-square" alt="HackerOne" style="margin: 10px;"/>
  </a>
  <a href="https://github.com/root4loot/rescope/actions/workflows/test-bugcrowd.yml">
    <img src="https://img.shields.io/badge/Bugcrowd-passing-brightgreen?style=flat-square" alt="Bugcrowd" style="margin: 10px;"/>
  </a>
  <a href="https://github.com/root4loot/rescope/actions/workflows/test-intigriti.yml">
    <img src="https://img.shields.io/badge/Intigriti-passing-brightgreen?style=flat-square" alt="Intigriti" style="margin: 10px;"/>
  </a>
  <a href="https://github.com/root4loot/rescope/actions/workflows/test-yeswehack.yml">
    <img src="https://img.shields.io/badge/YesWeHack-passing-brightgreen?style=flat-square" alt="YesWeHack" style="margin: 10px;"/>
  </a>
  <a href="https://github.com/root4loot/rescope/actions/workflows/test-hackenproof.yml">
    <img src="https://img.shields.io/badge/HackenProof-failing-red?style=flat-square" alt="HackenProof" style="margin: 10px;"/>
  </a>
</div>


## Installation

Requires Go 1.21 or later.

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
