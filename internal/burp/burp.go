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
// compatible JSON. Regex matches are split into groups. See internal scope package.
// Returns JSON data as byte
func Parse(L1, L2, L3 [][]string, Excludes []string) []byte {
	var host, scheme, port, path string
	fr, err := File.ReadFromRoot("configs/services", "internal")

	// L1 (all matches except IP-range and IP/CIDR)
	for _, submatch := range L1 {
		scheme = strings.TrimRight(submatch[1], "://")
		host = submatch[3]
		port = strings.TrimLeft(submatch[5], ":")
		path = submatch[6]

		//fmt.Println("S:" + scheme + "H:" + host + "PO:" + port + "PA:" + path)
		scheme, port = parseSchemeAndPort(fr, scheme, port)

		// parse regex for each group
		host = parseHost(host)
		path = parseFile(path)

		// check exclude
		isexclude := isExclude(Excludes, submatch[0])
		// add to list
		add(scheme, host, port, path, isexclude)
	}

	// L2 (IP-range match)
	for _, ipsets := range L2 {
		for _, ip := range ipsets {
			isexclude := isExclude(Excludes, ip)
			host := parseHost(ip)
			scheme = "Any"
			add(scheme, host, "^(80|443)$", path, isexclude)
		}
	}

	// L3 (IP/CIDR match)
	for _, ipsets := range L3 {
		for _, ip := range ipsets {
			isexclude := isExclude(Excludes, ip)
			host := parseHost(ip)
			scheme = "Any"
			add(scheme, host, "^(80|443)$", path, isexclude)
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
func add(scheme, host, port, path string, exclude bool) {
	if !exclude {
		incslice.Include = append(incslice.Include, Include{Enabled: true, File: path, Host: host, Port: port, Protocol: scheme})
	} else {
		exslice.Exclude = append(exslice.Exclude, Exclude{Enabled: true, File: path, Host: host, Port: port, Protocol: scheme})
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

// parseSchemeAndPort sets scheme and ports accordingly
// parseHost parse/set scheme & ports accordingly
// returns scheme, port (string) expressions
func parseSchemeAndPort(services []byte, scheme, port string) (string, string) {
	re := regexp.MustCompile(`([a-zA-Z0-9-]+)\s+(\d+)`) // for configs/services
	// re groups:     0. full match   -   [ftp 21]
	//                1. service      -   [ftp] 21
	//                2. port         -   ftp [21]

	if isVar(scheme) && !isVar(port) {
		// set corresponding port from configs/services
		scanner := bufio.NewScanner(strings.NewReader(string(services[:])))
		for scanner.Scan() {
			match := re.FindStringSubmatch(scanner.Text())
			if scheme == match[1] {
				port = "^" + match[2] + "$"
			}
		}
	} else if !isVar(scheme) && !isVar(port) {
		// set port to 80, 443
		port = "^(80|443)$"
	} else if isVar(scheme) && isVar(port) {
		// set whatever port + service port
		if scheme == "http" {
			port = "^(80|" + port + ")$"
		} else if scheme == "https" {
			port = "^(443|" + port + ")$"
		} else {
			port = "^" + port + "$"
		}
	} else if isVar(port) {
		port = "^" + port + "$"
	}

	// set "Any" when not http(s)
	if scheme != "http" && scheme != "https" {
		scheme = "Any"
	}

	return scheme, port
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

// parseFile parse path portion
// returns path (string) expression
func parseFile(path string) string {
	if isVar(path) {
		// replace wildcard
		path = strings.Replace(path, "*", `[\S]*`, -1)
		// escape '.'
		path = strings.Replace(path, ".", `\.`, -1)
		// add wildcard after dir suffix
		// note: this is not really needed as
		// burp will treat blank files as wildcards
		if strings.HasSuffix(path, "/") {
			path = path + `[\S]*`
		}
		path = "^" + path + "$"
	} else {
		path = `^[\S]*$`
	}
	return path
}

// isVar returns bool depening on len of var
func isVar(s string) bool {
	if len(s) > 0 {
		return true
	}
	return false
}
