package error

import (
	color "github.com/gookit/color"
	"log"
)

// BadJSON fatal error when JSON fails to unmarshal
func BadJSON() {
	log.Fatalf("\n%s Failed to parse JSON.", color.FgRed.Text("[!]"))
}

// NoScope fatal error when no scope tables was found
func NoScope(url string) {
	log.Fatalf("\n%s Unable to grab scope from %s. Is it public?", color.FgRed.Text("[!]"), url)
}

// NoResponse fatal error when http response failed
func NoResponse(url string) {
	log.Fatalf("\n%s No HTTP response from %s. Please make sure the URI is correct", color.FgRed.Text("[!]"), url)
}

// // BadResponse when http code is le 200 and ge 299
// func BadResponse(url string, code int, desc string) {
// 	log.Fatalf("\n%s %s %d (%s). Is the program public?", color.FgRed.Text("[!]"), url, code, desc)
// }
