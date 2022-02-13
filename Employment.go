package nmbrs

import (
	"encoding/xml"
	"io/ioutil"
	"net/http"

	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
)

type Employment_GetAll_AllEmployeesByCompanyResponse struct {
	XMLName  xml.Name `xml:"http://www.w3.org/2003/05/soap-envelope Envelope"`
	SoapBody struct {
		XMLName  xml.Name `xml:"Body"`
		Response *struct {
			XMLName xml.Name                                      `xml:"Employment_GetAll_AllEmployeesByCompanyResponse"`
			Result  Employment_GetAll_AllEmployeesByCompanyResult `xml:"Employment_GetAll_AllEmployeesByCompanyResult"`
		}
		Fault *Fault
	}
}

type Employment_GetAll_AllEmployeesByCompanyResult struct {
	XMLName                xml.Name                  `xml:"Employment_GetAll_AllEmployeesByCompanyResult"`
	EmployeeEmploymentItem *[]EmployeeEmploymentItem `xml:"EmployeeEmploymentItem"`
}

type EmployeeEmploymentItem struct {
	XMLName             xml.Name            `xml:"EmployeeEmploymentItem"`
	EmployeeID          int64               `xml:"EmployeeId"`
	EmployeeEmployments EmployeeEmployments `xml:"EmployeeEmployments"`
}
type EmployeeEmployments struct {
	XMLName    xml.Name     `xml:"EmployeeEmployments"`
	Employment []Employment `xml:"Employment"`
}

type Employment struct {
	XMLName          xml.Name `xml:"Employment"`
	EmploymentID     int64    `xml:"EmploymentId"`
	CreationDate     string   `xml:"CreationDate"`
	StartDate        string   `xml:"StartDate"`
	EndDate          *string  `xml:"EndDate"`
	InitialStartDate string   `xml:"InitialStartDate"`
}

type Employment_GetAll_AllEmployeesByCompany struct {
	XMLName   xml.Name `xml:"Employment_GetAll_AllEmployeesByCompany"`
	XMLNS     string   `xml:"xmlns,attr"`
	CompanyID int64    `xml:"CompanyID"`
}

func (service *Service) GetEmployments(companyID int64) (*[]EmployeeEmploymentItem, *errortools.Error) {
	xmlns := "https://api.nmbrs.nl/soap/v3/EmployeeService"
	bodyModel := service.GetSOAPEnvelope(xmlns, Employment_GetAll_AllEmployeesByCompany{
		XMLNS:     xmlns,
		CompanyID: companyID,
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

	r := Employment_GetAll_AllEmployeesByCompanyResponse{}

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

	return r.SoapBody.Response.Result.EmployeeEmploymentItem, nil
}
