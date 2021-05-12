package nmbrs

import "encoding/xml"

type SOAPEnvelope struct {
	XMLName xml.Name `xml:"soap12:Envelope"`
	XMLNS   string   `xml:"xmlns:soap12,attr"`
	Header  SOAPHeader
	Body    SOAPBody
}

type SOAPHeader struct {
	XMLName              xml.Name `xml:"soap12:Header"`
	AuthHeaderWithDomain AuthHeaderWithDomain
}

type AuthHeaderWithDomain struct {
	XMLName  xml.Name `xml:"AuthHeaderWithDomain"`
	XMLNS    string   `xml:"xmlns,attr"`
	Username string   `xml:"Username"`
	Domain   string   `xml:"Domain"`
	Token    string   `xml:"Token"`
}

type SOAPBody struct {
	XMLName xml.Name `xml:"soap12:Body"`
	Body    interface{}
}

func (service *Service) GetSOAPEnvelope(authHeaderXMLNS string, body interface{}) *SOAPEnvelope {
	soapEnvelope := SOAPEnvelope{
		XMLNS: "http://www.w3.org/2003/05/soap-envelope",
	}
	soapEnvelope.Header.AuthHeaderWithDomain = AuthHeaderWithDomain{
		XMLNS:    authHeaderXMLNS,
		Username: service.username,
		Domain:   service.domain,
		Token:    service.apiToken,
	}

	soapEnvelope.Body.Body = body

	return &soapEnvelope
}
