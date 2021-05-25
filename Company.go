package nmbrs

import (
	"encoding/xml"
	"io/ioutil"

	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
)

type List_GetAllResponse struct {
	XMLName  xml.Name `xml:"http://www.w3.org/2003/05/soap-envelope Envelope"`
	SoapBody struct {
		XMLName  xml.Name `xml:"Body"`
		Response *struct {
			XMLName xml.Name          `xml:"List_GetAllResponse"`
			Result  List_GetAllResult `xml:"List_GetAllResult"`
		}
		Fault *Fault
	}
}

type List_GetAllResult struct {
	XMLName xml.Name   `xml:"List_GetAllResult"`
	Company *[]Company `xml:"Company"`
}

type Company struct {
	XMLName             xml.Name `xml:"Company"`
	ID                  int64    `xml:"ID"`
	Number              int64    `xml:"Number"`
	Name                string   `xml:"Name"`
	LoonaangifteTijdvak string   `xml:"LoonaangifteTijdvak"`
	KvkNr               string   `xml:"KvkNr"`
}

type List_GetAll struct {
	XMLName xml.Name `xml:"List_GetAll"`
	XMLNS   string   `xml:"xmlns,attr"`
}

func (service *Service) GetCompanies() (*[]Company, *errortools.Error) {
	xmlns := "https://api.nmbrs.nl/soap/v3/CompanyService"
	bodyModel := service.GetSOAPEnvelope(xmlns, List_GetAll{
		XMLNS: xmlns,
	})

	requestConfig := go_http.RequestConfig{
		URL:       service.url("CompanyService.asmx"),
		BodyModel: bodyModel,
	}
	_, response, e := service.post(&requestConfig)
	if e != nil {
		return nil, e
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, errortools.ErrorMessagef("ioutil.ReadAll error: %s ", err.Error())
	}
	defer response.Body.Close()

	r := List_GetAllResponse{}

	err = xml.Unmarshal(body, &r)
	if err != nil {
		return nil, errortools.ErrorMessagef("xml.Unmarshal error: %s ", err.Error())
	}

	if e != nil {
		if r.SoapBody.Fault != nil {
			e.SetMessage(r.SoapBody.Fault.Reason)
		}

		return nil, e
	}

	if r.SoapBody.Response.Result.Company == nil {
		return nil, nil
	}

	return r.SoapBody.Response.Result.Company, nil
}
