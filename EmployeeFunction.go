package nmbrs

import (
	"encoding/xml"
	"io/ioutil"

	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
)

type Function_GetAll_AllEmployeesByCompany_V2Response struct {
	XMLName  xml.Name `xml:"http://www.w3.org/2003/05/soap-envelope Envelope"`
	SoapBody struct {
		XMLName  xml.Name `xml:"Body"`
		Response *struct {
			XMLName xml.Name                                       `xml:"Function_GetAll_AllEmployeesByCompany_V2Response"`
			Result  Function_GetAll_AllEmployeesByCompany_V2Result `xml:"Function_GetAll_AllEmployeesByCompany_V2Result"`
		}
		Fault *Fault
	}
}

type Function_GetAll_AllEmployeesByCompany_V2Result struct {
	XMLName                 xml.Name                   `xml:"Function_GetAll_AllEmployeesByCompany_V2Result"`
	EmployeeFunctionItem_V2 *[]EmployeeFunctionItem_V2 `xml:"EmployeeFunctionItem_V2"`
}

type EmployeeFunctionItem_V2 struct {
	XMLName           xml.Name            `xml:"EmployeeFunctionItem_V2"`
	EmployeeID        int64               `xml:"EmployeeId"`
	EmployeeFunctions []EmployeeFunctions `xml:"EmployeeFunctions"`
}

type EmployeeFunctions struct {
	XMLName          xml.Name           `xml:"EmployeeFunctions"`
	EmployeeFunction []EmployeeFunction `xml:"EmployeeFunction"`
}

type EmployeeFunction struct {
	XMLName      xml.Name `xml:"EmployeeFunction"`
	RecordID     int64    `xml:"RecordId"`
	Function     Function `xml:"Function"`
	CreationDate string   `xml:"CreationDate"`
	StartPeriod  int64    `xml:"StartPeriod"`
	StartYear    int64    `xml:"StartYear"`
}

type Function struct {
	XMLName     xml.Name `xml:"Function"`
	ID          int64    `xml:"Id"`
	Code        int64    `xml:"Code"`
	Description string   `xml:"Description"`
}

type Function_GetAll_AllEmployeesByCompany_V2 struct {
	XMLName   xml.Name `xml:"Function_GetAll_AllEmployeesByCompany_V2"`
	XMLNS     string   `xml:"xmlns,attr"`
	CompanyID int64    `xml:"CompanyID"`
}

func (service *Service) GetEmployeeFunctions(companyID int64) (*[]EmployeeFunctionItem_V2, *errortools.Error) {
	xmlns := "https://api.nmbrs.nl/soap/v3/EmployeeService"
	bodyModel := service.GetSOAPEnvelope(xmlns, Function_GetAll_AllEmployeesByCompany_V2{
		XMLNS:     xmlns,
		CompanyID: companyID,
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

	r := Function_GetAll_AllEmployeesByCompany_V2Response{}

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

	return r.SoapBody.Response.Result.EmployeeFunctionItem_V2, nil
}
