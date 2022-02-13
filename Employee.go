package nmbrs

import (
	"encoding/xml"
	"io/ioutil"
	"net/http"

	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
)

type List_GetByCompanyResponse struct {
	XMLName  xml.Name `xml:"http://www.w3.org/2003/05/soap-envelope Envelope"`
	SoapBody struct {
		XMLName  xml.Name `xml:"Body"`
		Response *struct {
			XMLName xml.Name                `xml:"List_GetByCompanyResponse"`
			Result  List_GetByCompanyResult `xml:"List_GetByCompanyResult"`
		}
		Fault *Fault
	}
}

type List_GetByCompanyResult struct {
	XMLName  xml.Name    `xml:"List_GetByCompanyResult"`
	Employee *[]Employee `xml:"Employee"`
}

type Employee struct {
	XMLName     xml.Name `xml:"Employee"`
	ID          int64    `xml:"Id"`
	Number      int64    `xml:"Number"`
	DisplayName string   `xml:"DisplayName"`
}

type List_GetByCompany struct {
	XMLName      xml.Name `xml:"List_GetByCompany"`
	XMLNS        string   `xml:"xmlns,attr"`
	CompanyID    int64    `xml:"CompanyId"`
	EmployeeType int64    `xml:"EmployeeType"`
}

func (service *Service) GetEmployees(companyID int64, employeeType int64) (*[]Employee, *errortools.Error) {
	xmlns := "https://api.nmbrs.nl/soap/v3/EmployeeService"
	bodyModel := service.GetSOAPEnvelope(xmlns, List_GetByCompany{
		XMLNS:        xmlns,
		CompanyID:    companyID,
		EmployeeType: employeeType,
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

	r := List_GetByCompanyResponse{}

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

	return r.SoapBody.Response.Result.Employee, nil
}
