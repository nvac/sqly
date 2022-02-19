package sqly

import "encoding/xml"

type scripts struct {
	XMLName xml.Name `xml:"scripts"`
	Scripts []script `xml:"script"`
}

type script struct {
	XMLName xml.Name `xml:"script"`
	Name    string   `xml:"name,attr"`
	Content string   `xml:",chardata"`
}
