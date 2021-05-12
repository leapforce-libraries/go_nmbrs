package nmbrs

import (
	"bytes"
	"encoding/xml"
	"io/ioutil"

	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
)

type Response struct {
	XMLName  xml.Name `xml:"http://schemas.xmlsoap.org/soap/envelope/ Envelope"`
	SoapBody *SOAPBodyResponse
}

type SOAPBodyResponse struct {
	XMLName  xml.Name `xml:"Body"`
	Response *CostCenter_GetListResponse
	Fault    *Fault
}

type Fault struct {
	XMLName xml.Name `xml:"Fault"`
	Code    string   `xml:"Code>Value"`
	Reason  string   `xml:"Reason>Text"`
	Detail  string   `xml:"Detail"`
}

type CostCenter_GetListResponse struct {
	XMLName xml.Name                 `xml:"CostCenter_GetListResponse"`
	Result  CostCenter_GetListResult `xml:"CostCenter_GetListResult"`
}

type CostCenter_GetListResult struct {
	XMLName    xml.Name     `xml:"CostCenter_GetListResult"`
	CostCenter []CostCenter `xml:"CostCenter"`
}

type CostCenter struct {
	XMLName     xml.Name `xml:"CostCenter"`
	Code        int64    `xml:"Code"`
	Description string   `xml:"Description"`
	ID          int64    `xml:"Id"`
}

func (service *Service) GetCostCenters(companyID int64) (*[]CostCenter, *errortools.Error) {

	data := struct {
		APIToken  string
		Username  string
		Domain    string
		CompanyID int64
	}{
		service.apiToken,
		service.username,
		service.domain,
		companyID,
	}

	doc := &bytes.Buffer{}
	// Replacing the doc from template with actual req values
	err := service.templates.ExecuteTemplate(doc, "", &data)
	if err != nil {
		return nil, errortools.ErrorMessagef("template.Execute error. %s ", err.Error())
	}

	buffer := &bytes.Buffer{}
	encoder := xml.NewEncoder(buffer)
	err = encoder.Encode(doc.String())
	if err != nil {
		return nil, errortools.ErrorMessagef("encoder.Encode error. %s ", err.Error())
	}

	requestConfig := go_http.RequestConfig{
		URL:       service.url("CompanyService.asmx"),
		BodyModel: bytes.NewBuffer([]byte(doc.String())),
	}
	_, response, e := service.post(&requestConfig)

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, errortools.ErrorMessagef("ioutil.ReadAll error. %s ", err.Error())
	}
	defer response.Body.Close()

	r := &Response{}
	err = xml.Unmarshal(body, &r)
	if err != nil {
		return nil, errortools.ErrorMessagef("xml.Unmarshal error. %s ", err.Error())
	}

	if e != nil {
		if r.SoapBody.Fault != nil {
			e.SetMessage(r.SoapBody.Fault.Reason)
		}

		return nil, e
	}

	return &r.SoapBody.Response.Result.CostCenter, nil
}
