//
// Written By : Daniel Antonsen (@root4loot)
//
// Distributed Under MIT License
// Copyrights (C) 2019 root4loot
//

package files

import (
	"bufio"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
)

// IsExist takes a filename and determines whether or or not the file exists on the filesystem
// returns boolean
func IsExist(fn string) bool {
	_, err := os.Stat(fn)
	if os.IsNotExist(err) {
		return false
	}
	return true
}

// Create attempts to create a file for a given filepath
// returns file descriptor and error, if any
func Create(fn string) (*os.File, error) {
	f, err := os.Create(fn)
	return f, err
}

// Open attempts to open a file for a given filepath
// returns the file descriptor error, if any
func Open(file string) (*os.File, error) {
	f, err := os.Open(file)
	return f, err
}

// Read attempts to read data from a given file descriptor
// returns file data and error, if any
func Read(file *os.File) ([]byte, error) {
	reader := bufio.NewReader(file)
	data, err := ioutil.ReadAll(reader)
	return data, err
}

// Write attempts to write given data to an already opened file
func Write(f *os.File, data []byte) (int, error) {
	w, err := f.Write(data)
	return w, err
}

// ReadFromRoot attempts to read data from a file relative to the parent of this func and the project root
// the first argument is the file you want to read, relative to project root
// the second argument is the path between the file you want to read, and the callers parent
// Example: caller is 'project/pkg/files' and you want to read 'project/dir/somefile.xml'
// ReadFromRoot("configs/somefile.xml", "pkg")
// returns file data and error, if any
func ReadFromRoot(file, trim string) ([]byte, error) {
	// runtime caller
	_, caller, _, ok := runtime.Caller(0)
	if !ok {
		panic("No caller information")
	}

	// path to caller (project/dir/file)
	caller = path.Dir(caller)
	// path to caller parent dir (project/dir)
	callerDir := filepath.Dir(caller)
	// trim dir to get the project root (project)
	projectRoot := strings.TrimRight(callerDir, trim)
	// get file relative to project root (project/file)
	file = projectRoot + file

	// open and read
	fo, err := Open(file)
	fr, err := Read(fo)
	// close
	defer fo.Close()

	return fr, err
}
