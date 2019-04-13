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

	"github.com/gookit/color"

	burp "github.com/root4loot/rescope/internal/burp"
	cli "github.com/root4loot/rescope/internal/cli"
	scope "github.com/root4loot/rescope/internal/scope"
	url "github.com/root4loot/rescope/internal/url"
	zap "github.com/root4loot/rescope/internal/zap"
	file "github.com/root4loot/rescope/pkg/file"
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

			// check err
			if err, ok := err.(*os.PathError); ok {
				fmt.Println("\n%s Unable to read %s.", color.FgRed.Text("[!]"), f)
				log.Fatal(err)
			}

			// close
			defer fd.Close()

			// add to list of files
			fds = append(fds, fd)

		}

		// get infile(s) contents
		for _, fd := range fds {
			data, err := file.Read(fd)
			if err != nil {
				fmt.Println("\n%s Unable to read contents of %s", color.FgRed.Text("[!]"), fd.Name())
				log.Fatal(err)
			}

			// add to lists
			scopes = append(scopes, string(data[:]))
			source = append(source, fd.Name())
		}
	}

	// Deal with -u|--url
	scopes, source, bbaas = url.BBaas(a.URLs, scopes, source)

	// Identify scope targets
	m := scope.Match{}
	m = scope.Parse(m, scopes, source, a.Silent, a.IncTag, a.ExTag, bbaas)
	if m.Counter == 0 {
		log.Fatalf("%s Nothing to do here ¯\\_(ツ)_/¯", color.FgRed.Text("[!]"))
	}

	// Parse to burp/zap
	if a.Burp {
		fmt.Printf("%s Parsing to JSON (Burp Suite)", color.FgGray.Text("[-]"))
		buf = burp.Parse(m.L1, m.L2, m.L3, m.Excludes)
	} else if a.Zap {
		fmt.Printf("%s Parsing to XML (OWASP ZAP)", color.FgGray.Text("[-]"))
		buf = zap.Parse(m.L1, m.L2, m.L3, m.Excludes, a.Scopename)
	}

	// Attempt to create outfile
	outfile, err := file.Create(a.Outfile)
	if err != nil {
		log.Fatalf("\n%s Failed to create file at %s. Bad permisisons?", color.FgRed.Text("[!]"), outfile.Name())
	}

	// Write to outfile assuming we have permissions
	meta, err := file.Write(outfile, buf)

	if a.Burp {
		fmt.Printf("\n%s Done. Wrote %v bytes to %s\n", color.FgGreen.Text("[✓]"), meta, outfile.Name())
	} else if a.Zap {
		fmt.Printf("\n%s Done. Wrote %v bytes to %s\n", color.FgGreen.Text("[✓]"), meta, outfile.Name())
	}
	fmt.Println("")
}
