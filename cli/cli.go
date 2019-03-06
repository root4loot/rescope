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
	 ~ r o o t 4 l o o t    |_|
	
	
	Setting Excludes (optional):
	  specify !EXCLUDE in -i <file>, followed by targets you wish to exclude. Anything succeeding this tag will be explicity excluded from scope.
	
	
	Example usage: 
	 rescope burp -i scope.txt -o burp.json
	 rescope zap  -i scope1.txt -i scope2.txt -o zap.context --name CoolScope
	 `
	version := 0.1
	parser := argparse.NewParser("rescope [burp|zap]", banner)
	//usage := parser.Usage
	c := CLI{}

	i := parser.List("i", "infile", &argparse.Options{Help: "File (scope) to be parsed (required)\n\t\t Can be set multiple times"})
	o := parser.String("o", "outfile", &argparse.Options{Help: "File to write parsed results (required)"})
	s := parser.Flag("s", "silent", &argparse.Options{Help: "Do not print identified targets"})
	n := parser.String("n", "name", &argparse.Options{Help: "Name of ZAP context"})
	e := parser.String("e", "extag", &argparse.Options{Help: "Custom exclude tag (default: !EXCLUDE)"})
	v := parser.Flag("", "version", &argparse.Options{Help: "Print version"})
	_ = n
	_ = s
	_ = v

	_ = parser.Parse(os.Args)

	c.Infiles = *i
	c.Outfile = *o
	c.Scopename = *n
	c.Silent = *s
	c.ExTag = *e

	red := color.New(color.FgRed).SprintFunc()

	// remove timestamp from exits
	log.SetFlags(0)

	// slice of error strings
	var argErr []string

	// set CLI values
	for _, arg := range os.Args {
		if strings.ToLower(arg) == "burp" {
			c.Command = "burp"
		} else if strings.ToLower(arg) == "zap" {
			c.Command = "zap"
		}
	}

	// check for args and add to list
	if !isVar(c.Command) {
		argErr = append(argErr, "Missing command argument (burp|zap)")
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

	// Version
	if c.version {
		fmt.Println("rescope", version)
		os.Exit(0)
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

// GetCommand (valid "commands" [burp, zap])
// returns command string
func GetCommand(c CLI) string {
	return c.Command
}

// IsCommand compare command name
// returns bool
func IsCommand(c CLI, command string) bool {
	if c.Command == command {
		return true
	}
	return false
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
