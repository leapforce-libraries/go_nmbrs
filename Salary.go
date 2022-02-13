package nmbrs

import (
	"encoding/xml"
	"io/ioutil"
	"net/http"

	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
)

type Salary_GetListResponse struct {
	XMLName  xml.Name `xml:"http://www.w3.org/2003/05/soap-envelope Envelope"`
	SoapBody struct {
		XMLName  xml.Name `xml:"Body"`
		Response *struct {
			XMLName xml.Name             `xml:"Salary_GetListResponse"`
			Result  Salary_GetListResult `xml:"Salary_GetListResult"`
		}
		Fault *Fault
	}
}

type Salary_GetListResult struct {
	XMLName xml.Name  `xml:"Salary_GetListResult"`
	Salary  *[]Salary `xml:"Salary"`
}

type Salary struct {
	XMLName     xml.Name     `xml:"Salary"`
	Value       float64      `xml:"Value"`
	Type        string       `xml:"Type"`
	SalaryTable *SalaryTable `xml:"SalaryTable"`
	StartDate   string       `xml:"StartDate"`
}

type SalaryTable struct {
	XMLName xml.Name `xml:"SalaryTable"`
	Schaal  *Schaal  `xml:"Schaal"`
	Trede   *Trede   `xml:"Trede"`
}

type Schaal struct {
	XMLName            xml.Name `xml:"Schaal"`
	Scale              string   `xml:"Scale"`
	ScaleDescription   string   `xml:"SchaalDescription"`
	ScaleValue         float64  `xml:"ScaleValue"`
	ScalePercentageMax float64  `xml:"ScalePercentageMax"`
	ScalePercentageMin float64  `xml:"ScalePercentageMin"`
}

type Trede struct {
	XMLName         xml.Name `xml:"Trede"`
	Step            int64    `xml:"Step"`
	StepDescription string   `xml:"StepDescription"`
	StepValue       float64  `xml:"StepValue"`
}

type Salary_GetList struct {
	XMLName    xml.Name `xml:"Salary_GetList"`
	XMLNS      string   `xml:"xmlns,attr"`
	EmployeeID int64    `xml:"EmployeeId"`
	Period     int64    `xml:"Period"`
	Year       int64    `xml:"Year"`
}

func (service *Service) GetSalaries(employeeID int64, period int64, year int64) (*[]Salary, *errortools.Error) {
	xmlns := "https://api.nmbrs.nl/soap/v3/EmployeeService"
	bodyModel := service.GetSOAPEnvelope(xmlns, Salary_GetList{
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

	r := Salary_GetListResponse{}

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

	return r.SoapBody.Response.Result.Salary, nil
}
