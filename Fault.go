package nmbrs

import "encoding/xml"

type Fault struct {
	XMLName xml.Name `xml:"Fault"`
	Code    string   `xml:"Code>Value"`
	Reason  string   `xml:"Reason>Text"`
	Detail  string   `xml:"Detail"`
}
