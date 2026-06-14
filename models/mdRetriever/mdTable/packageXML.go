package mdTable

import "encoding/xml"

type Package struct {
	XMLName xml.Name `xml:"Package"`
	Types   []Types  `xml:"types"`
	Version string   `xml:"version"`
}

// Types represents a <types> element
type Types struct {
	Members []string `xml:"members"`
	Name    string   `xml:"name"`
}
