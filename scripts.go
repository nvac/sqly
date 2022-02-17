package sqly

import "encoding/xml"

type scripts struct {
	XMLName xml.Name `xml:"scripts"`
	Scripts []script `xml:"script"`
}

type script struct {
	XMLName  xml.Name `xml:"script"`
	Name     string   `xml:"name,attr"`
	Database string   `xml:"database,attr"`
	Content  string   `xml:",chardata"`
}
