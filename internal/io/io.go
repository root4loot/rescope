//
// Written By : Daniel Antonsen (@root4loot)
//
// Distributed Under MIT License
// Copyrights (C) 2019 root4loot
//

package io

import (
	"bufio"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
)

// IsFileExist returns bool
func IsFileExist(fn string) bool {
	_, err := os.Stat(fn)
	if os.IsNotExist(err) {
		return false
	}
	return true
}

// OpenFile returns file, error
func OpenFile(file string) (*os.File, error) {
	f, err := os.Open(file)
	return f, err
}

// ReadFile returns file data, error
func ReadFile(file *os.File) ([]byte, error) {
	reader := bufio.NewReader(file)
	data, err := ioutil.ReadAll(reader)
	return data, err
}

// CreateFile returns created file, error
func CreateFile(fn string) (*os.File, error) {
	f, err := os.Create(fn)
	return f, err
}

// WriteFile takes file, data and writes to it
// returns bytes written, error
func WriteFile(f *os.File, data []byte) (int, error) {
	w, err := f.Write(data)
	return w, err
}

// ReadFileFromProjectRoot opens, read and returns a file from project root dir
func ReadFileFromProjectRoot(file, trim string) (fr []byte) {
	// runtime caller
	_, caller, _, ok := runtime.Caller(0)
	if !ok {
		panic("No caller information")
	}

	// path to caller (project/dir/package/package.go)
	caller = path.Dir(caller)
	// path to caller parent dir (project/dir)
	callerDir := filepath.Dir(caller)
	// trim dir to get the project root (project)
	projectRoot := strings.TrimRight(callerDir, trim)
	// get file from project root (project/file)
	file = projectRoot + file

	// open and read
	fo, _ := OpenFile(file)
	fr, _ = ReadFile(fo)
	// close
	defer fo.Close()

	return fr
}
