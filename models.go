package sqly

import (
	"encoding/xml"
)

type databases struct {
	XMLName   xml.Name   `xml:"databases"`
	Databases []database `xml:"database"`
}

type database struct {
	XMLName         xml.Name `xml:"database"`
	Name            string   `xml:"name,attr"`
	Driver          string   `xml:"driver,attr"`
	Source          string   `xml:"source,attr"`
	Environment     string   `xml:"environment,attr"`
	MaxIdleConns    *int     `xml:"maxIdleConns,attr"`
	MaxOpenConns    *int     `xml:"maxOpenConns,attr"`
	ConnMaxLifetime *int     `xml:"connMaxLifetime,attr"`
	ConnMaxIdleTime *int     `xml:"connMaxIdleTime,attr"`
}

type scripts struct {
	XMLName xml.Name `xml:"scripts"`
	Scripts []script `xml:"script"`
}

type script struct {
	XMLName xml.Name `xml:"script"`
	Name    string   `xml:"name,attr"`
	Content string   `xml:",chardata"`
}
