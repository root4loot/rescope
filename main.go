//
// Written By : Daniel Antonsen (@root4loot)
//
// Distributed Under MIT License
// Copyrights (C) 2019 root4loot
//

package main

import (
	"fmt"
	"log"
	"os"

	"github.com/fatih/color"

	burp "github.com/root4loot/rescope/internal/burp"
	cli "github.com/root4loot/rescope/internal/cli"
	io "github.com/root4loot/rescope/internal/io"
	scope "github.com/root4loot/rescope/internal/scope"
	zap "github.com/root4loot/rescope/internal/zap"
)

func main() {
	// file descriptors
	var fds []*os.File
	// data to be written to outfile
	var buf []byte
	// slice containing scope definitions
	var source []string
	// slice containing scope origins (filename or bugbounty url)
	var scopes []string
	// indicates whether BBaaS url was found
	var bbaas bool
	// struct containing various args
	a := cli.ArgParse()

	// Read infiles and add contents to scope list
	if a.Infiles != nil {
		for _, f := range a.Infiles {
			if file.IsExist(f) != true {
				log.Fatalf("\n%s File %s not found", color.FgRed.Text("[!]"), f)
		}
	}

	// attempt to open infiles
		for _, f := range a.Infiles {
			fd, err := file.Open(f)

		if err, ok := err.(*os.PathError); ok {
			fmt.Printf("\n%s Failed to read %s.", red("[!]"), f)
			log.Fatal(err)
		}
	}

	// file data
	var scopes []string

	// attempt to read infiles contents
	for _, fd := range fds {
		data, err := io.ReadFile(fd)
		if err != nil {
			fmt.Printf("\n%s Failed to read contents of %s", red("[!]"), fd.Name())
			log.Fatal(err)
		}
		// append to scopes
		scopes = append(scopes, string(data[:]))
	}

	// Identify scope targets
	m := scope.Match{}
	m = scope.Parse(m, scopes, source, a.Silent, a.IncTag, a.ExTag, bbaas)

	// parse to burp/zap
	if c.Burp {
		fmt.Printf("%s Parsing to JSON (Burp Suite)", grey("[-]"))
		buf = burp.Parse(m.L1, m.L2, m.L3, m.Excludes)
		fmt.Printf("\n%s Done", green("[✓]"))
	} else if c.Zap {
		fmt.Printf("%s Parsing to XML (OWASP ZAP)", grey("[-]"))
		buf = zap.Parse(m.L1, m.L2, m.L3, m.Excludes, c.Scopename)
		fmt.Printf("\n%s Done", green("[✓]"))
	}

	// attempt to create outfile
	outfile, err := io.CreateFile(c.Outfile)
	if err != nil {
		fmt.Printf("\n%s Failed to create file at %s. Bad permisisons?", red("[!]"), outfile.Name())
		log.Fatal(err)
	}

	// write to outfile assuming we have permissions
	meta, err := io.WriteFile(outfile, buf)

	if c.Burp {
		fmt.Printf("\n%s Wrote %v bytes to %s\n\n", green("[✓]"), meta, outfile.Name())
	} else if c.Zap {
		fmt.Printf("\n%s Wrote %v bytes to %s\n\n", green("[✓]"), meta, outfile.Name())
	}
}
