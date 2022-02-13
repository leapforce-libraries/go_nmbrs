package nmbrs

import (
	"encoding/xml"
	"io/ioutil"
	"net/http"

	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
)

type HourComponentFixed_GetResponse struct {
	XMLName  xml.Name `xml:"http://www.w3.org/2003/05/soap-envelope Envelope"`
	SoapBody struct {
		XMLName  xml.Name `xml:"Body"`
		Response *struct {
			XMLName xml.Name                     `xml:"HourComponentFixed_GetResponse"`
			Result  HourComponentFixed_GetResult `xml:"HourComponentFixed_GetResult"`
		}
		Fault *Fault
	}
}

type HourComponentFixed_GetResult struct {
	XMLName       xml.Name              `xml:"HourComponentFixed_GetResult"`
	HourComponent *[]HourComponentFixed `xml:"HourComponent"`
}

type HourComponentFixed struct {
	XMLName  xml.Name `xml:"HourComponent"`
	ID       int64    `xml:"Id"`
	HourCode int64    `xml:"HourCode"`
	Hours    float64  `xml:"Hours"`
}

type HourComponentFixed_Get struct {
	XMLName    xml.Name `xml:"HourComponentFixed_Get"`
	XMLNS      string   `xml:"xmlns,attr"`
	EmployeeID int64    `xml:"EmployeeId"`
	Period     int64    `xml:"Period"`
	Year       int64    `xml:"Year"`
}

func (service *Service) GetHourComponentFixeds(employeeID int64, period int64, year int64) (*[]HourComponentFixed, *errortools.Error) {
	xmlns := "https://api.nmbrs.nl/soap/v3/EmployeeService"
	bodyModel := service.GetSOAPEnvelope(xmlns, HourComponentFixed_Get{
		XMLNS:      xmlns,
		EmployeeID: employeeID,
		Period:     period,
		Year:       year,
	})

	requestConfig := go_http.RequestConfig{
		Method:    http.MethodPost,
		URL:       service.url("EmployeeService.asmx"),
		BodyModel: bodyModel,
	}

	_, response, e := service.httpRequest(&requestConfig)
	if e != nil {
		return nil, e
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, errortools.ErrorMessagef("ioutil.ReadAll error: %s ", err.Error())
	}
	defer response.Body.Close()

	r := HourComponentFixed_GetResponse{}

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

	return r.SoapBody.Response.Result.HourComponent, nil
}
