package config

import "encoding/xml"

type ZapConfig struct {
	XMLName xml.Name `xml:"configuration"`
	Context struct {
		Name       string   `xml:"name"`
		Desc       string   `xml:"desc"`
		Inscope    string   `xml:"inscope"`
		Incregexes []string `xml:"incregexes,omitempty"`
		Excregexes []string `xml:"excregexes,omitempty"`
		Tech       struct {
			Include []string `xml:"include"`
		} `xml:"tech"`
		Urlparser struct {
			Class  string `xml:"class"`
			Config string `xml:"config"`
		} `xml:"urlparser"`
		Postparser struct {
			Class  string `xml:"class"`
			Config string `xml:"config"`
		} `xml:"postparser"`
		Authentication struct {
			Type      int    `xml:"type"`
			Strategy  string `xml:"strategy"`
			Pollurl   string `xml:"pollurl"`
			Polldata  string `xml:"polldata"`
			Pollfreq  int    `xml:"pollfreq"`
			Pollunits string `xml:"pollunits"`
		} `xml:"authentication"`
		Forceduser string `xml:"forceduser"`
		Session    struct {
			Type int `xml:"type"`
		} `xml:"session"`
		Authorization struct {
			Type  int `xml:"type"`
			Basic struct {
				Header string `xml:"header"`
				Body   string `xml:"body"`
				Logic  string `xml:"logic"`
				Code   int    `xml:"code"`
			} `xml:"basic"`
		} `xml:"authorization"`
	} `xml:"context"`
}
