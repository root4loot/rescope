//
// Written By : Daniel Antonsen (@root4loot)
//
// Distributed Under MIT License
// Copyrights (C) 2019 root4loot
//

package zap

import (
	"bufio"
	"strings"

	file "github.com/root4loot/rescope/pkg/file"
)

var includes []string
var excludes []string

// Parse takes slices containing regex matches and turns them into Zap
// compatible XML (Context)
// Returns xml data as bytes
func Parse(L1, L2, L3 [][]string, Excludes []string, scopeName string) []byte {
	var oldxml []string
	var newxml []string

	// read default scope template
	fr, err := file.ReadFromRoot("configs/default.context", "pkg")

	// loop template and append each line to var
	scanner := bufio.NewScanner(strings.NewReader(string(fr[:])))
	for scanner.Scan() {
		oldxml = append(oldxml, scanner.Text())
	}

	// L1 all matches except IP-range & IP/CIDR
	for _, submatch := range L1 {
		match := submatch[0]
		protocol := submatch[1]
		port := submatch[4]
		target := parse(match, protocol, port) // [0] fullmatch

		if !isExclude(Excludes, submatch[0]) {
			item := "<incregexes>" + target + "</incregexes>"
			includes = append(includes, item)
		}
	}

	// L2 ip-range
	for _, ipsets := range L2 {
		for _, set := range ipsets {
			ip := parse(set, "", "")
			if !isExclude(Excludes, ip) {
				item := "<incregexes>" + ip + "</incregexes>"
				includes = append(includes, item)
			}
		}
	}

	// l3 ip/CIDR
	for _, ipsets := range L3 {
		for _, set := range ipsets {
			ip := parse(set, "", "")
			if !isExclude(Excludes, ip) {
				item := "<incregexes>" + ip + "</incregexes>"
				includes = append(includes, item)
			}
		}
	}

	// add to excludes
	for _, item := range Excludes {
		item := parse(item, "", "")
		item = "<excregexes>" + item + "</excregexes>"
		excludes = append(excludes, item)
	}

	// replace line 3 in template with scope name
	oldxml[3] = "<name>" + scopeName + "</name>"

	// append each line of template (oldxml) to newxml.
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
	// convert string to byte, separated with newline
	xml := []byte(strings.Join(newxml, "\n"))
	return xml
}

// parse host, url, etc to regex
func parse(target, protocol, port string) string {
	line := target

	// if no protocol &  no port // example.com
	if len(protocol) == 0 && len(port) == 0 {
		// scope only http/https
		line = `http(s)?://` + line

		// if port, but no protocol // example.com:8080
	} else if len(protocol) == 0 && len(port) != 0 {
		// scope any protocol
		line = `\w+://` + line

		// if protocol != http(s) // ftp://example.com
	} else if len(protocol) != 0 && protocol != "http://" && protocol != "https://" {
		line = protocol + `://` + line
	}

	// escape '.'
	line = strings.Replace(line, ".", `\.`, -1)
	// escape '/'
	line = strings.Replace(line, "/", `\/`, -1)
	// replace wildcard
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
