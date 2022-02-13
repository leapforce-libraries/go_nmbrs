package nmbrs

import (
	"encoding/xml"
	"io/ioutil"
	"net/http"

	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
)

type CostCenter_GetAllEmployeesByCompanyResponse struct {
	XMLName  xml.Name `xml:"http://www.w3.org/2003/05/soap-envelope Envelope"`
	SoapBody struct {
		XMLName  xml.Name `xml:"Body"`
		Response *struct {
			XMLName xml.Name                                  `xml:"CostCenter_GetAllEmployeesByCompanyResponse"`
			Result  CostCenter_GetAllEmployeesByCompanyResult `xml:"CostCenter_GetAllEmployeesByCompanyResult"`
		}
		Fault *Fault
	}
}

type CostCenter_GetAllEmployeesByCompanyResult struct {
	XMLName             xml.Name               `xml:"CostCenter_GetAllEmployeesByCompanyResult"`
	EmployeeCostCenters *[]EmployeeCostCenters `xml:"EmployeeCostCenters"`
}

type EmployeeCostCenters struct {
	XMLName     xml.Name      `xml:"EmployeeCostCenters"`
	EmployeeID  int64         `xml:"EmployeeId"`
	CostCenters []CostCenters `xml:"CostCenters"`
}

type CostCenters struct {
	XMLName            xml.Name             `xml:"CostCenters"`
	EmployeeCostCenter []EmployeeCostCenter `xml:"EmployeeCostCenter"`
}

type EmployeeCostCenter struct {
	XMLName     xml.Name    `xml:"EmployeeCostCenter"`
	ID          int64       `xml:"Id"`
	CostCenter  CostCenter  `xml:"CostCenter"`
	Kostensoort Kostensoort `xml:"Kostensoort"`
	Percentage  float64     `xml:"Percentage"`
	Default     bool        `xml:"Default"`
}

type Kostensoort struct {
	XMLName     xml.Name `xml:"Kostensoort"`
	Code        int64    `xml:"Code"`
	Description string   `xml:"Description"`
}

type CostCenter_GetAllEmployeesByCompany struct {
	XMLName   xml.Name `xml:"CostCenter_GetAllEmployeesByCompany"`
	XMLNS     string   `xml:"xmlns,attr"`
	CompanyID int64    `xml:"CompanyId"`
	Period    int64    `xml:"Period"`
	Year      int64    `xml:"Year"`
}

func (service *Service) GetEmployeeCostCenters(companyID int64, period int64, year int64) (*[]EmployeeCostCenters, *errortools.Error) {
	xmlns := "https://api.nmbrs.nl/soap/v3/EmployeeService"
	bodyModel := service.GetSOAPEnvelope(xmlns, CostCenter_GetAllEmployeesByCompany{
		XMLNS:     xmlns,
		CompanyID: companyID,
		Period:    period,
		Year:      year,
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

	r := CostCenter_GetAllEmployeesByCompanyResponse{}

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

	return r.SoapBody.Response.Result.EmployeeCostCenters, nil
}
