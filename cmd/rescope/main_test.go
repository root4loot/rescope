package main

import (
	"flag"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseCLI_WithIncludeExclude(t *testing.T) {
	os.Args = []string{"rescope", "--include-list", "include.txt", "--exclude-list", "exclude.txt"}
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	_, cli, err := parseCLI()
	assert.NoError(t, err, "Expected no error with valid flags")

	assert.Equal(t, "include.txt", cli.IncludeList, "Expected IncludeList to be set")
	assert.Equal(t, "exclude.txt", cli.ExcludeList, "Expected ExcludeList to be set")
}

func TestGetInputFileContents_ValidFiles(t *testing.T) {
	includeContent := "https://bugcrowd.com/program\nexample.com"
	excludeContent := "https://hackerone.com/program\noutofscope.com"

	os.WriteFile("include.txt", []byte(includeContent), 0644)
	os.WriteFile("exclude.txt", []byte(excludeContent), 0644)
	defer os.Remove("include.txt")
	defer os.Remove("exclude.txt")

	cli := CLI{
		IncludeList: "include.txt",
		ExcludeList: "exclude.txt",
	}

	includes, excludes, err := cli.getInputFileContents()
	assert.NoError(t, err, "Expected no error when reading valid files")
	assert.Equal(t, strings.Split(includeContent, "\n"), includes, "Expected includes to match file content")
	assert.Equal(t, strings.Split(excludeContent, "\n"), excludes, "Expected excludes to match file content")
}

func TestGetInputFileContents_NonExistentFiles(t *testing.T) {
	cli := CLI{
		IncludeList: "nonexistent_include.txt",
		ExcludeList: "nonexistent_exclude.txt",
	}

	_, _, err := cli.getInputFileContents()
	assert.Error(t, err, "Expected error when reading non-existent files")
}

func TestProcessStdin(t *testing.T) {
	input := "https://bugcrowd.com/program\nexample.com"
	r, w, _ := os.Pipe()
	os.Stdin = r
	w.WriteString(input)
	w.Close()

	targets := processStdin()
	expectedTargets := strings.Fields(input)

	assert.Equal(t, expectedTargets, targets, "Expected targets to match stdin input")
}

func TestInvalidFlags(t *testing.T) {
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	fs.String("include-list", "", "")
	fs.String("exclude-list", "", "")
	fs.String("auth-hackerone", "", "")
	fs.String("auth-intigriti", "", "")
	fs.String("auth-yeswehack", "", "")
	fs.String("auth-bugcrowd", "", "")
	fs.String("output-file", "", "")
	fs.Bool("output-text", false, "")
	fs.Bool("output-burp", false, "")
	fs.Bool("output-zap", false, "")
	fs.Bool("output-json", false, "")
	fs.Bool("output-json-lines", false, "")
	fs.Bool("expand-ip-ranges", false, "")
	fs.Int("concurrency", 5, "")
	fs.Bool("debug", false, "")
	fs.Bool("help", false, "")
	fs.Bool("version", false, "")

	err := fs.Parse([]string{"--invalid-flag"})
	assert.Error(t, err, "Expected error with invalid flag")
}
