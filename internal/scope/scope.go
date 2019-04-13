//
// Written By : Daniel Antonsen (@root4loot)
//
// Distributed Under MIT License
// Copyrights (C) 2019 root4loot
//

package scope

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"regexp"
	"strconv"
	"strings"

	"github.com/fatih/color"
)

// Match contains lists of regex matches
type Match struct {
	L1       [][]string // all except ip-range & CIDR
	L2       [][]string // ip-range
	L3       [][]string // ip/CIDR
	Excludes []string   // to be excluded
}

// Parse function takes a slice containing scope file data and
// applies regex to each line in order to extract targets from scope-
// matched targets are split into groups varying on type (host, url, iprange, etc)
// Returns a Match object containing all lists
func Parse(m Match, scopes []string, command string, files []string, silent bool, exTag string) Match {
	var exclude bool
	grey := color.New(color.Faint).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()

	// Set Tag used to indicate beginning of Excludes
	if len(exTag) == 0 {
		exTag = "!EXCLUDE"
	}

	r1 := regexp.MustCompile(`([a-z]+:\/\/)?(\*\.)?([a-z0-9-.]+(\.[a-z]+))(:\d+)?([A-Za-z0-9-._~:/?#@!$&'*+,;=]+)?`)
	// Groups: 1.  [ftp]://sub.example.com:25/d/foo.bar    // scheme
	//         2.   ftp://[*.]example.com:25/d/foo.bar     // wildcarded subdomain
	//	       3.   ftp://[sub.example.com]:25/d/foo.bar   // host
	//         4.   ftp://sub.example[.com]:25/d/foo.bar   // extension
	//         5.   ftp://sub.example.com[:25]/d/foo.bar   // port
	//         6.   ftp://sub.example.com:25[/d/foo.bar]   // path

	r2 := regexp.MustCompile(`(\d+\.\d+\.\d+\.)(\d+)-(\d+)`)
	// Matches IP-Range
	// Groups: 1.  (192.168.0).1-255    // IP minus last host portion
	//         2.   192.168.0.(1)-255   // start
	//         3.   192.168.0.1-(255)   // end

	r3 := regexp.MustCompile(`[\d\.]+\/\d{2}`)
	// Matches IP/CIDR

	for i, scope := range scopes {
		counter := 0
		scanner := bufio.NewScanner(strings.NewReader(scope))
		exclude = false // reset flag on each run

		fmt.Printf("%s Grabbing targets from [%s]\n", grey("[-]"), files[i])
		for scanner.Scan() {
			m1 := r1.FindAllStringSubmatch(scanner.Text(), -1)
			m2 := r2.FindAllStringSubmatch(scanner.Text(), -1)
			m3 := r3.FindAllString(scanner.Text(), -1)

			// Check exclude
			if strings.Contains(scanner.Text(), exTag) {
				exclude = true
			}

			if m3 != nil { // m3 ip/CIDR
				for _, arr := range m3 {
					// not interested in those ending with '.'
					if strings.HasSuffix(arr, ".") {
						continue
					}
					hosts, err := hostsFromCIDR(arr)
					if err != nil {
						fmt.Printf("\n%s Failed to parse IP/CIDR: %s", red("[!]"), m3)
						log.Fatal(err)
					} else {
						m.L3 = append(m.L3, hosts)
						counter++
						printFound(arr, exclude, silent)
					}
					if exclude == true {
						for _, host := range hosts {
							m.Excludes = append(m.Excludes, host)
						}
					}
				}

			} else if m2 != nil { // m2 ip-range
				for _, arr := range m2 {
					// not interested in those ending with '.'
					if strings.HasSuffix(arr[0], ".") {
						continue
					}

					hosts, err := hostsFromRange(arr)
					if err != nil {
						fmt.Printf("\n%s Failed to parse IP-range: %s", red("[!]"), m2[0])
						log.Fatal(err)
					} else {
						counter++
						m.L2 = append(m.L2, hosts)
						printFound(arr[0], exclude, silent)
						if exclude == true {
							for _, host := range hosts {
								m.Excludes = append(m.Excludes, host)
							}
						}
					}
				}

			} else if m1 != nil { // m1 all others
				// not interested in those ending with '.'
				for _, arr := range m1 {
					if strings.HasSuffix(arr[0], ".") {
						continue
					}
					m.L1 = append(m.L1, arr)
					counter++
					printFound(arr[0], exclude, silent)
					if exclude == true {
						m.Excludes = append(m.Excludes, arr[0])
					}

				}

			}
		}

		if counter == 0 {
			fmt.Printf("%s No targets found in %s. Wrong file?", red("[!]"), files[i])
		}
	}
	return m
}

// prints item in color depending on whether it is part of include or exclude
func printFound(item string, exclude bool, silent bool) {
	if exclude == true {
		if !silent {
			color.Red(" - " + item)
		}
	} else {
		if !silent {
			color.Green(" + " + item)
		}
	}
}

// hostsFromRange takes a m2 slice containing IP-range substrings
// converts range to a list of hosts and returns this
func hostsFromRange(m []string) ([]string, error) {
	ip := m[1] // (192.168.)0.1-255

	start, err := strconv.Atoi(m[2]) // 192.168.0.(1)-255
	end, err := strconv.Atoi(m[3])   // 192.168.(0).(1)-(255)
	var ips []string

	// loop range and append to list
	for i := start; i <= end; i++ {
		ip := ip + strconv.Itoa(i)
		ips = append(ips, ip)
	}
	return ips, err
}

// hostsFromCIDR takes a m3 slice containing IP/CIDR substrings
// converts CIDR to list of hosts and returns this
func hostsFromCIDR(cidr string) ([]string, error) {
	ip, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, err
	}
	var ips []string
	// we only want the IP
	for ip := ip.Mask(ipnet.Mask); ipnet.Contains(ip); inc(ip) {
		ips = append(ips, ip.String())
	}
	return ips[1 : len(ips)-1], nil
}

func inc(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}

