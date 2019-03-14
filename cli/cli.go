//
// Written By : Daniel Antonsen (@root4loot)
//
// Distributed Under MIT License
// Copyrights (C) 2019 root4loot
//

package cli

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/akamensky/argparse"
	"github.com/fatih/color"
)

// CLI contains arguments
type CLI struct {
	Command   string
	Infiles   []string
	Outfile   string
	Burp      bool
	Zap       bool
	Silent    bool
	Scopename string
	verbose   string
	version   bool
	ExTag     string
}

// Parse cli arguments
func Parse() CLI {
	banner := `
  _ __ ___  ___  ___ ___  _ __   ___ 
 | '__/ _ \/ __|/ __/ _ \| '_ \ / _ \
 | | |  __/\__ \ (_| (_) | |_) |  __/
 |_|  \___||___/\___\___/| .__/ \___|
  ~ r o o t 4 l o o t    |_|    v.0.2
	
Setting Excludes (optional):
  specify !EXCLUDE in -i <file> prior to targets you wish to exclude.         

Example Usage: 
  rescope --burp -i scope.txt -o burp.json
  rescope --zap  -i scope1.txt -i scope2.txt -o zap.context --name CoolScope              

Upgrading:
  go get -u github.com/root4loot/rescope 

Documentation:
  https://github.com/root4loot/rescope
`
	version := "0.2"
	parser := argparse.NewParser("rescope", banner)
	red := color.New(color.FgRed).SprintFunc()
	//usage := parser.Usage
	c := CLI{}
	b := parser.Flag("b", "burp", &argparse.Options{Help: "Parse to Burp Suite JSON (required option)"})
	z := parser.Flag("z", "zap", &argparse.Options{Help: "Parse to OWASP ZAP XML (required option)"})
	i := parser.List("i", "infile", &argparse.Options{Help: "File (scope) to be parsed (required)\n\t\t Infile can be set multiple times"})
	o := parser.String("o", "outfile", &argparse.Options{Help: "File to write parsed results (required)"})
	n := parser.String("n", "name", &argparse.Options{Help: "Name of ZAP context"})
	e := parser.String("e", "extag", &argparse.Options{Help: "Custom exclude tag (default: !EXCLUDE)"})
	s := parser.Flag("s", "silent", &argparse.Options{Help: "Do not print identified targets"})
	v := parser.Flag("", "version", &argparse.Options{Help: "Print version"})

	_ = n
	_ = s
	_ = v

	_ = parser.Parse(os.Args)

	c.Burp = *b
	c.Zap = *z
	c.Infiles = *i
	c.Outfile = *o
	c.Scopename = *n
	c.Silent = *s
	c.ExTag = *e
	c.version = *v

	// remove timestamp from exits
	log.SetFlags(0)

	// slice of error strings
	var argErr []string

	// print version
	if c.version {
		fmt.Println("rescope version " + version)
		os.Exit(0)
	}

	// check for args and add to list
	if !c.Burp && !c.Zap {
		argErr = append(argErr, "Missing program identifier [-b|--burp] [-z|--zap]")
	}
	if !isList(c.Infiles) {
		argErr = append(argErr, "Missing infile (-i <file>)")
	}
	if !isVar(c.Outfile) {
		argErr = append(argErr, "Missing outfile (-o <file>)")
	}

	// print arg errors from list
	if len(argErr) > 0 {
		for i := 1; i <= len(argErr); i++ {
			fmt.Printf("%s %s\n", red("[!]"), argErr[i-1])
		}
		os.Exit(1)
	}

	// check/set scopename
	if c.Command == "zap" {
		if !isVar(c.Scopename) {
			c.Scopename = setScopeName()
		}
	}
	return c
}

// setScopeName for Zap Context
// returns scopename
func setScopeName() string {
	c := color.New(color.Faint).SprintFunc()
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("%s Enter name of Scope (required for Zap): ", c("[>]"))
	name, _ := reader.ReadString('\n')
	name = strings.TrimSuffix(name, "\n")
	return name
}

// GetScopeName for Zap Context
func GetScopeName(c CLI) string {
	return c.Scopename
}

// isVar check if var is empty or not
// returns bool
func isVar(v string) bool {
	if len(v) > 0 {
		return true
	}
	return false
}

// check if list is empty or not
// returns bool
func isList(l []string) bool {
	if len(l) > 0 {
		return true
	}
	return false
}
