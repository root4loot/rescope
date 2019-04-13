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
	"github.com/gookit/color"
)

// Args flags
type Args struct {
	Command   string
	Infiles   []string
	URLs      []string
	Outfile   string
	Burp      bool
	Zap       bool
	Silent    bool
	Scopename string
	verbose   string
	version   bool
	IncTag    string
	ExTag     string
}

// ArgParse and check arguments etc
func ArgParse() Args {
	banner := `
  _ __ ___  ___  ___ ___  _ __   ___ 
 | '__/ _ \/ __|/ __/ _ \| '_ \ / _ \
 | | |  __/\__ \ (_| (_) | |_) |  __/
 |_|  \___||___/\___\___/| .__/ \___|
  ~ r o o t 4 l o o t    |_|     v0.3
	
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
	version := "0.3"
	parser := argparse.NewParser("rescope", banner)

	//usage := parser.Usage
	a := Args{}
	b := parser.Flag("b", "burp", &argparse.Options{Help: "Parse to Burp Suite JSON (required option A)"})
	z := parser.Flag("z", "zap", &argparse.Options{Help: "Parse to OWASP ZAP XML (required option A)"})
	u := parser.List("u", "url", &argparse.Options{Help: "Public bug bounty program URL (required option B)\n\t\t URL can be set multiple times"})
	i := parser.List("i", "infile", &argparse.Options{Help: "File (scope) to be parsed (required option B)\n\t\t Infile can be set multiple times"})
	o := parser.String("o", "outfile", &argparse.Options{Help: "File to write parsed results (required)"})
	n := parser.String("n", "name", &argparse.Options{Help: "Name of ZAP context"})
	ex := parser.String("", "itag", &argparse.Options{Help: "Custom include tag (default: !INCLUDE)"})
	in := parser.String("", "etag", &argparse.Options{Help: "Custom exclude tag (default: !EXCLUDE)"})
	s := parser.Flag("s", "silent", &argparse.Options{Help: "Do not print identified targets"})
	ver := parser.Flag("", "version", &argparse.Options{Help: "Display version"})

	_ = parser.Parse(os.Args)

	a.Burp = *b
	a.Zap = *z
	a.Infiles = *i
	a.URLs = *u
	a.Outfile = *o
	a.Scopename = *n
	a.Silent = *s
	a.IncTag = *in
	a.ExTag = *ex
	a.version = *ver

	// remove timestamp from exits
	log.SetFlags(0)

	// slice of error strings
	var argErr []string

	// print version
	if a.version {
		fmt.Println("rescope v" + version)
		os.Exit(0)
	}

	// check for args and add to list
	if !a.Burp && !a.Zap {
		argErr = append(argErr, "Missing program identifier: (-b|--burp) | (-z|--zap)")
	}
	if !isList(a.Infiles) && !isList(a.URLs) {
		argErr = append(argErr, "Missing infile/URL: (-i|--infile <file>) | (-u|--url <url>)")
	}
	if !isVar(a.Outfile) {
		argErr = append(argErr, "Missing outfile [-o|--outfile <file>]")
	}
	if len(strings.Split(a.Outfile, ".")) < 2 {
		argErr = append(argErr, "Outfile must have an extension (-o filename.ext)")
	}

	// print arg errors from list
	if len(argErr) > 0 {
		for i := 1; i <= len(argErr); i++ {
			fmt.Printf("%s %s\n", color.FgRed.Text("[!]"), argErr[i-1])
		}
		os.Exit(1)
	}

	// check/set scopename
	if a.Command == "zap" {
		if !isVar(a.Scopename) {
			a.Scopename = setScopeName()
		}
	}
	return a
}

// setScopeName for Zap Context
// returns scopename
func setScopeName() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("%s Enter name of Scope (required for ZAP): ", color.FgGray.Text("[>]"))
	name, _ := reader.ReadString('\n')
	name = strings.TrimSuffix(name, "\n")
	return name
}

// GetScopeName for Zap Context
func GetScopeName(a Args) string {
	return a.Scopename
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
