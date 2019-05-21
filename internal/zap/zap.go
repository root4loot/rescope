//
// Written By : Daniel Antonsen (@root4loot)
//
// Distributed Under MIT License
// Copyrights (C) 2019 root4loot
//

package zap

import (
	"bufio"
	"regexp"
	"strings"

	file "github.com/root4loot/rescope/pkg/file"
)

var includes []string
var excludes []string

// Parse takes slices containing regex matches and turns them into Zap compatible XML (Context)
// Returns xml data as bytes
func Parse(Includes, Excludes [][]string, scopeName string) []byte {
	var oldxml []string
	var newxml []string
	var cludes [][][]string

	cludes = append(cludes, Includes)
	cludes = append(cludes, Excludes)

	// read default scope template
	fr := file.ReadFromRoot("configs/default.context", "pkg")

	// Loop template and append each line to var
	scanner := bufio.NewScanner(strings.NewReader(string(fr[:])))
	for scanner.Scan() {
		oldxml = append(oldxml, scanner.Text())
	}

	for i, clude := range cludes {
		for _, item := range clude {
			ip := regexp.MustCompile(`\d+\.\d+\.\d+\.\d+`)

			if ip.MatchString(item[0]) {
				for _, ip := range item {
					ip = parse(ip, "", "")
					if i == 0 {
						item := "<incregexes>" + ip + "</incregexes>"
						includes = append(includes, item)
					} else {
						item := "<excregexes>" + ip + "</excregexes>"
						excludes = append(excludes, item)
					}
				}
			} else {
				full := item[0]
				scheme := string(item[1])
				port := string(item[6])
				target := parse(full, scheme, port) // [0] fullmatch

				if i == 0 {
					item := "<incregexes>" + target + "</incregexes>"
					includes = append(includes, item)
				} else {
					item := "<excregexes>" + target + "</excregexes>"
					excludes = append(excludes, item)
				}
			}
		}
	}

	// Replace line 3 in template with scope name
	oldxml[3] = "<name>" + scopeName + "</name>"

	// Append each line of template (oldxml) to newxml.
	// at line 5, begin appending []includes and []excludes
	for i, v := range oldxml {
		newxml = append(newxml, v)
		if i == 5 {
			for _, v := range includes {
				newxml = append(newxml, v)
			}
			for _, v := range excludes {
				newxml = append(newxml, v)
			}
		}
	}

	// Convert string to byte, separated with newline
	xml := []byte(strings.Join(newxml, "\n"))
	return xml
}

// parse host, url, etc to regex
func parse(target, scheme, port string) string {
	line := target

	// if no scheme, no port // example.com
	if len(scheme) == 0 && len(port) == 0 {
		// scope only http/https
		line = `http(s)?://` + line

		// if port, but no scheme // example.com:8080
	} else if len(scheme) == 0 && len(port) != 0 {
		line = `[a-z]+://` + line + port

		// if port and scheme
	} else if len(scheme) != 0 && len(port) != 0 {
		line = scheme + `://` + line + port
	}

	// escape '.'
	line = strings.Replace(line, ".", `\.`, -1)
	// escape '/'
	line = strings.Replace(line, "/", `\/`, -1)
	// replace wildcards
	line = strings.Replace(line, "*", `[\S]*`, -1)
	// Zap needs this to scope URL params
	line = `^` + line + `[\S]*$`

	return line
}

// isExclude takes a 2d slice and a string
// checks whether string was found in list
// returns bool
func isExclude(Excludes []string, target string) bool {
	for _, exclude := range Excludes {
		if target == exclude {
			return true
		}
	}
	return false
}
