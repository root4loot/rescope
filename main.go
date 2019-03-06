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

	burp "github.com/root4loot/rescope/burp"
	cli "github.com/root4loot/rescope/cli"
	io "github.com/root4loot/rescope/io"
	scope "github.com/root4loot/rescope/scope"
	zap "github.com/root4loot/rescope/zap"
)

func main() {
	// data to be written to outfile
	var buf []byte
	// struct containing various args
	c := cli.Parse()

	// fancy colors
	grey := color.New(color.Faint).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()

	// determine if infiles exists
	for _, f := range c.Infiles {
		if io.IsFileExist(f) != true {
			fmt.Printf("\n%s Couldn't find file %s. Does it exist?", red("[!]"), f)
		}
	}

	// file descriptors
	var fds []*os.File

	// attempt to open infiles
	for _, f := range c.Infiles {
		fd, err := io.OpenFile(f)
		// remember to close file
		defer fd.Close()
		// add to fds
		fds = append(fds, fd)

		if err, ok := err.(*os.PathError); ok {
			fmt.Printf("\n%s Failed to read file at location %s. Bad permissions?", red("[!]"), f)
			log.Fatal(err)
		}
	}

	// file data
	var scopes []string

	// attempt to read infiles contents
	for _, fd := range fds {
		data, err := io.ReadFile(fd)
		if err != nil {
			fmt.Printf("\n%s Failed to read contents of file %s", red("[!]"), fd.Name())
			log.Fatal(err)
		}
		// append to scopes
		scopes = append(scopes, string(data[:]))
	}

	// apply regex matching to scopes
	m := scope.Match{}
	m = scope.Parse(m, scopes, c.Command, c.Infiles, c.Silent, c.ExTag)

	// parse to burp/zap
	if cli.IsCommand(c, "burp") {
		fmt.Printf("%s Parsing to JSON (Burp Suite)", grey("[-]"))
		buf = burp.Parse(m.L1, m.L2, m.L3, m.Excludes)
		fmt.Printf("\n%s Done", green("[✓]"))
	} else if cli.IsCommand(c, "zap") {
		fmt.Printf("%s Parsing to XML (OWASP ZAP)", grey("[-]"))
		buf = zap.Parse(m.L1, m.L2, m.L3, m.Excludes, c.Scopename)
		fmt.Printf("\n%s Done", green("[✓]"))
	}

	// attempt to create outfile
	outfile, err := io.CreateFile(c.Outfile)
	if err != nil {
		fmt.Printf("\n%s Failed to create file at location %s. Bad permisisons?", red("[!]"), outfile.Name())
		log.Fatal(err)
	}

	// write to outfile assuming we have permissions as
	// file was created
	meta, err := io.WriteFile(outfile, buf)

	if cli.IsCommand(c, "burp") {
		fmt.Printf("\n%s Wrote %v bytes to %s\n\n", green("[✓]"), meta, outfile.Name())
	} else if cli.IsCommand(c, "zap") {
		fmt.Printf("\n%s Wrote %v bytes to %s\n\n", green("[✓]"), meta, outfile.Name())
	}
}
