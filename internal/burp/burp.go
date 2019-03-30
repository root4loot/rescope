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
	"encoding/json"
	"fmt"
	"strings"
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
// compatible JSON. Regex matches are split into groups.
// Returns JSON data as byte
func Parse(L1, L2, L3 [][]string, Excludes []string) []byte {
	var host, protocol, port, file, wport string

	// L1 all matches except IP-range and IP/CIDR
	for _, submatch := range L1 {
		protocol = submatch[1]
		host = submatch[2]
		port = submatch[4]
		file = submatch[5]

		// parse regex for each group
		protocol, wport = parseProtocol(protocol)
		host = parseHost(host)
		port = parsePort(port, wport)
		file = parseFile(file)

		// check exclude
		isexclude := isExclude(Excludes, submatch[0])
		// add to list
		add(protocol, host, port, file, isexclude)
	}

	// L2 IP-range match
	for _, ipsets := range L2 {
		for _, ip := range ipsets {
			isexclude := isExclude(Excludes, ip)
			host := parseHost(ip)
			protocol = "Any"
			add(protocol, host, "^(80|443)$", file, isexclude)
		}
	}

	// L3 IP/CIDR match
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

// parseProtocol sets port depending on the protocol
func parseProtocol(protocol string) (string, string) {
	var port string
	protocol = strings.Replace(protocol, "://", "", -1)
	if isVar(protocol) {
		if protocol == "http" {
			port = "80"
		} else if protocol == "https" {
			port = "443"
		} else {
			protocol = "Any" // burp does not like anything but http(s)
		}
	} else {
		protocol = "Any"
	}
	return protocol, port
}

// parseHost to regex
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

// parsePort to regex
func parsePort(port string, wport string) string {
	if isVar(port) {
		port = strings.TrimPrefix(port, ":")
	}

	if isVar(port) && isVar(wport) {
		port = "^(" + port + "|" + wport + ")$"
	} else if isVar(port) {
		port = "^" + port + "$"
	} else if isVar(wport) {
		port = "^" + wport + "$"
	} else {
		port = "^(80|443)$"
	}
	return port
}

// parseFile to regex
func parseFile(file string) string {
	if isVar(file) {
		// replace wildcard
		file = strings.Replace(file, "*", `[\S]*`, -1)
		// escape '.'
		file = strings.Replace(file, ".", `\.`, -1)
		// add wildcard after dir suffix
		// note: the following statement is not really needed as blank files/dirs are
		// wildcarded by Burp. The reason for leaving this in is to
		// avoid being reliant on that fact.
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
