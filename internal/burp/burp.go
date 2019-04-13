//
// Written By : Daniel Antonsen (@root4loot)
//
// Distributed Under MIT License
// Copyrights (C) 2019 root4loot
//

// Package burp involves parsing list of scope targets to Burp
// compatible JSON (Regex)
package burp

import (
	"bufio"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	File "github.com/root4loot/rescope/pkg/file"
)

// Scope is the JSON structure that burp wants
type Scope struct {
	Target struct {
		Scope struct {
			AdvancedMode bool      `json:"advanced_mode"`
			Exclude      []Exclude `json:"exclude"`
			Include      []Include `json:"include"`
		} `json:"scope"`
	} `json:"target"`
}

// Include host details
type Include struct {
	Enabled  bool   `json:"enabled"`
	File     string `json:"file"`
	Host     string `json:"host"`
	Port     string `json:"port"`
	Protocol string `json:"protocol"`
}

// Exclude host details
type Exclude struct {
	Enabled  bool   `json:"enabled"`
	File     string `json:"file"`
	Host     string `json:"host"`
	Port     string `json:"port"`
	Protocol string `json:"protocol"`
}

// IncludeSlice contains Include structs
type IncludeSlice struct {
	Include []Include
}

// ExcludeSlice contains Enclude structs
type ExcludeSlice struct {
	Exclude []Exclude
}

// var containing IncludeSlice
var incslice IncludeSlice

// var containing ExcludeSlice
var exslice ExcludeSlice

// Parse takes slices containing regex matches and turns them into Burp
// compatible JSON. Regex matches are split into groups. See Scope package.
// Returns JSON data as byte
func Parse(L1, L2, L3 [][]string, Excludes []string) []byte {
	fr, err := File.ReadFromRoot("configs/services", "internal")

	// L1 (all matches except IP-range and IP/CIDR)
	for _, submatch := range L1 {
		scheme = strings.TrimRight(submatch[1], "://")
		host = submatch[3]
		port = strings.TrimLeft(submatch[5], ":")
		path = submatch[6]

		protocol, port = parseProtocolAndPort(fr, protocol, port)

		// parse regex for each group
		host = parseHost(host)
		//port = parsePort(port, wport, protocol)
		file = parseFile(file)

		// check exclude
		isexclude := isExclude(Excludes, submatch[0])
		// add to list
		add(protocol, host, port, file, isexclude)
	}

	// L2 (IP-range match)
	for _, ipsets := range L2 {
		for _, ip := range ipsets {
			isexclude := isExclude(Excludes, ip)
			host := parseHost(ip)
			protocol = "Any"
			add(protocol, host, "^(80|443)$", file, isexclude)
		}
	}

	// L3 (IP/CIDR match)
	for _, ipsets := range L3 {
		for _, ip := range ipsets {
			isexclude := isExclude(Excludes, ip)
			host := parseHost(ip)
			protocol = "Any"
			add(protocol, host, "^(80|443)$", file, isexclude)
		}
	}

	// scope object
	scope := Scope{}
	scope.Target.Scope.AdvancedMode = true
	// add include/exclude slices
	scope.Target.Scope.Include = incslice.Include
	scope.Target.Scope.Exclude = exslice.Exclude

	// parse pretty json
	json, err := json.MarshalIndent(scope, "", "  ")
	if err != nil {
		fmt.Println("json err:", err)
	}
	return json
}

// add match to appropriate list
func add(protocol, host, port, file string, exclude bool) {
	if !exclude {
		incslice.Include = append(incslice.Include, Include{Enabled: true, File: file, Host: host, Port: port, Protocol: protocol})
	} else {
		exslice.Exclude = append(exslice.Exclude, Exclude{Enabled: true, File: file, Host: host, Port: port, Protocol: protocol})
	}
}

// isExclude takes a 2d slice and a string
// returns bool depending on whether the string was found in slice
func isExclude(Excludes []string, item string) bool {
	for _, exclude := range Excludes {
		if item == exclude {
			return true
		}
	}
	return false
}

// parseProtocolAndPort sets protocol and ports accordingly
// parseHost parse/set protocol & ports accordingly
// returns protocol, port (string) expressions
func parseProtocolAndPort(services []byte, protocol, port string) (string, string) {
	re := regexp.MustCompile(`([a-zA-Z0-9-]+)\s+(\d+)`) // for configs/services
	// re groups:     0. full match   - (ftp 21)
	//                1. service      - (ftp) 21
	//                2. port         - ftp (21)

	if isVar(protocol) && !isVar(port) {
		// set corresponding port from configs/services
		scanner := bufio.NewScanner(strings.NewReader(string(services[:])))
		for scanner.Scan() {
			match := re.FindStringSubmatch(scanner.Text())
			if protocol == match[1] {
				port = "^" + match[2] + "$"
			}
		}
	} else if !isVar(protocol) && !isVar(port) {
		// set port to 80, 443
		port = "^(80|443)$"
	} else if isVar(protocol) && isVar(port) {
		// set whatever port + service port
		if protocol == "http" {
			port = "^(80|" + port + ")$"
		} else if protocol == "https" {
			port = "^(443|" + port + ")$"
		} else {
			port = "^" + port + "$"
		}
	} else if isVar(port) {
		port = "^" + port + "$"
	}

	// set "Any" when not http(s)
	if protocol != "http" && protocol != "https" {
		protocol = "Any"
	}

	return protocol, port
}

// parseHost parse host portion
// returns host (string) expression
func parseHost(host string) string {
	if isVar(host) {
		if strings.Contains(host, "*") {
			host = strings.Replace(host, "*", `[\S]*`, -1)
		}
		host = strings.Replace(host, ".", `\.`, -1)
		host = "^" + host + "$"
	}
	return host
}

// parseFile parse file portion
// returns file (string) expression
func parseFile(file string) string {
	if isVar(file) {
		// replace wildcard
		file = strings.Replace(file, "*", `[\S]*`, -1)
		// escape '.'
		file = strings.Replace(file, ".", `\.`, -1)
		// add wildcard after dir suffix
		// note: this is not really needed as
		// burp will treat blank files as wildcards
		if strings.HasSuffix(file, "/") {
			file = file + `[\S]*`
		}
		file = "^" + file + "$"
	} else {
		file = `^[\S]*$`
	}
	return file
}

// isVar returns bool depening on len of var
func isVar(s string) bool {
	if len(s) > 0 {
		return true
	}
	return false
}
