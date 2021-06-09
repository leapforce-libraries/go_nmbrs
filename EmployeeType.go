package nmbrs

import (
	"encoding/xml"
	"io/ioutil"

	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
)

type EmployeeType_GetListResponse struct {
	XMLName  xml.Name `xml:"http://www.w3.org/2003/05/soap-envelope Envelope"`
	SoapBody struct {
		XMLName  xml.Name `xml:"Body"`
		Response *struct {
			XMLName xml.Name                   `xml:"EmployeeType_GetListResponse"`
			Result  EmployeeType_GetListResult `xml:"EmployeeType_GetListResult"`
		}
		Fault *Fault
	}
}

type EmployeeType_GetListResult struct {
	XMLName      xml.Name        `xml:"EmployeeType_GetListResult"`
	EmployeeType *[]EmployeeType `xml:"EmployeeType"`
}

type EmployeeType struct {
	XMLName     xml.Name `xml:"EmployeeType"`
	ID          int64    `xml:"Id"`
	Description string   `xml:"Description"`
}

type EmployeeType_GetList struct {
	XMLName xml.Name `xml:"EmployeeType_GetList"`
	XMLNS   string   `xml:"xmlns,attr"`
}

func (service *Service) GetEmployeeTypes() (*[]EmployeeType, *errortools.Error) {
	xmlns := "https://api.nmbrs.nl/soap/v3/EmployeeTypeService"
	bodyModel := service.GetSOAPEnvelope(xmlns, EmployeeType_GetList{
		XMLNS: xmlns,
	})

	requestConfig := go_http.RequestConfig{
		URL:       service.url("EmployeeTypeService.asmx"),
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

	r := EmployeeType_GetListResponse{}

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

	return r.SoapBody.Response.Result.EmployeeType, nil
}
