//
// Written By : Daniel Antonsen (@root4loot)
//
// Distributed Under MIT License
// Copyrights (C) 2021 root4loot
//

package cli

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/akamensky/argparse"
	"github.com/gookit/color"
)

// Args flags
type Args struct {
	Command    string
	Infiles    []string
	URLs       []string
	Outfile    string
	Burp       bool
	Zap        bool
	Raw        bool
	Silent     bool
	Scopename  string
	verbose    string
	version    bool
	IncTag     string
	ExTag      string
	ResolveAll bool
	Avoid3P    bool
}

func ArgParse() Args {
	version := "2.1"
	banner := `
  _ __ ___  ___  ___ ___  _ __   ___ 
 | '__/ _ \/ __|/ __/ _ \| '_ \ / _ \
 | | |  __/\__ \ (_| (_) | |_) |  __/
 |_|  \___||___/\___\___/| .__/ \___|
  @ r o o t 4 l o o t    |_|     ` + version + `
https://github.com/root4loot/rescope 
     
Example Usage:
  rescope -u hackerone.com/security -o burpscope.json  
  rescope -u hackerone.com/security --zap -o zapscope.context 
  rescope --zap  -i scope.txt -o zap.context --name CoolScope
`
	parser := argparse.NewParser("rescope", banner)

	//usage := parser.Usage
	a := Args{}
	z := parser.Flag("z", "zap", &argparse.Options{Required: false, Help: "Export scope to ZAP-compatible XML instead of default (Burp JSON)"})
	u := parser.List("u", "url", &argparse.Options{Required: false, Help: "Public bug bounty program URL"})
	i := parser.List("i", "infile", &argparse.Options{Required: false, Help: "File (scope) to be parsed"})
	n := parser.String("n", "name", &argparse.Options{Required: false, Help: "Name of ZAP context"})
	o := parser.String("o", "outfile", &argparse.Options{Required: false, Help: "Save results to given filename"})
	s := parser.Flag("s", "silent", &argparse.Options{Required: false, Help: "Do not print identified targets"})
	r := parser.Flag("r", "raw", &argparse.Options{Required: false, Help: "Export raw scope-definitions to list of text"})
	ex := parser.String("", "itag", &argparse.Options{Required: false, Help: "Custom include tag (default: !INCLUDE)"})
	in := parser.String("", "etag", &argparse.Options{Required: false, Help: "Custom exclude tag (default: !EXCLUDE)"})
	res := parser.Flag("", "resolveConflicts", &argparse.Options{Required: false, Help: "Resolve all exclude conflicts"})
	avoid3P := parser.Flag("", "avoid3P", &argparse.Options{Required: false, Help: "Avoid all third party resources"})
	ver := parser.Flag("", "version", &argparse.Options{Required: false, Help: "Display version"})
	_ = parser.Parse(os.Args)

	a.Zap = *z
	a.Raw = *r
	a.Infiles = *i
	a.URLs = *u
	a.Outfile = *o
	a.Scopename = *n
	a.Silent = *s
	a.IncTag = *in
	a.ExTag = *ex
	a.ResolveAll = *res
	a.Avoid3P = *avoid3P
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

	if !isList(a.Infiles) && !isList(a.URLs) {
		argErr = append(argErr, "Missing (-i <file>) or bugbounty (-u <url>)")
	}

	if !isVar(a.Outfile) {
		if a.Zap && a.Raw {
			argErr = append(argErr, "You cannot have both (-z|--zap) and (-r|--raw) in one command (mutually exclusive)")
		} else if a.Zap {
			if !isVar(a.Scopename) {
				a.Scopename = setScopeName()
				re_nonalpha, _ := regexp.Compile("[^A-Za-z0-9]+")
				a.Scopename = re_nonalpha.ReplaceAllString(a.Scopename, "_")
			}
			a.Outfile = "./scope_" + a.Scopename + "_zap" + ".context"
		} else if a.Raw {
			a.Outfile = "./scope_raw.txt"
		} else {
			a.Outfile = "./scope_burp.json"
		}
	} else if len(strings.Split(a.Outfile, ".")) < 2 {
		argErr = append(argErr, "Outfile must have an extension (E.g: -o coolScope.burp)")
	}

	// print arg errors from list
	if len(argErr) > 0 {
		for i := 1; i <= len(argErr); i++ {
			fmt.Printf("%s %s\n", color.FgRed.Text("[!]"), argErr[i-1])
		}
		os.Exit(1)
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

// btoi bool to int
func btoi(b bool) int {
	if b {
		return 1
	}
	return 0
}
