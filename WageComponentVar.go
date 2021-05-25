package nmbrs

import (
	"encoding/xml"
	"io/ioutil"

	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
)

type WageComponentVar_GetResponse struct {
	XMLName  xml.Name `xml:"http://www.w3.org/2003/05/soap-envelope Envelope"`
	SoapBody struct {
		XMLName  xml.Name `xml:"Body"`
		Response *struct {
			XMLName xml.Name                   `xml:"WageComponentVar_GetResponse"`
			Result  WageComponentVar_GetResult `xml:"WageComponentVar_GetResult"`
		}
		Fault *Fault
	}
}

type WageComponentVar_GetResult struct {
	XMLName       xml.Name            `xml:"WageComponentVar_GetResult"`
	WageComponent *[]WageComponentVar `xml:"WageComponent"`
}

type WageComponentVar struct {
	XMLName xml.Name `xml:"WageComponent"`
	ID      int64    `xml:"Id"`
	Code    int64    `xml:"Code"`
	Value   float64  `xml:"Value"`
}

type WageComponentVar_Get struct {
	XMLName    xml.Name `xml:"WageComponentVar_Get"`
	XMLNS      string   `xml:"xmlns,attr"`
	EmployeeID int64    `xml:"EmployeeId"`
	Period     int64    `xml:"Period"`
	Year       int64    `xml:"Year"`
}

func (service *Service) GetWageComponentVars(employeeID int64, period int64, year int64) (*[]WageComponentVar, *errortools.Error) {
	xmlns := "https://api.nmbrs.nl/soap/v3/EmployeeService"
	bodyModel := service.GetSOAPEnvelope(xmlns, WageComponentVar_Get{
		XMLNS:      xmlns,
		EmployeeID: employeeID,
		Period:     period,
		Year:       year,
	})

	requestConfig := go_http.RequestConfig{
		URL:       service.url("EmployeeService.asmx"),
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

	r := WageComponentVar_GetResponse{}

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

	return r.SoapBody.Response.Result.WageComponent, nil
}
