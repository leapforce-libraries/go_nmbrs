package nmbrs

import (
	"encoding/xml"
	"io/ioutil"
	"net/http"

	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
)

type Contract_GetAll_AllEmployeesByCompanyResponse struct {
	XMLName  xml.Name `xml:"http://www.w3.org/2003/05/soap-envelope Envelope"`
	SoapBody struct {
		XMLName  xml.Name `xml:"Body"`
		Response *struct {
			XMLName                                     xml.Name                                    `xml:"Contract_GetAll_AllEmployeesByCompanyResponse"`
			Contract_GetAll_AllEmployeesByCompanyResult Contract_GetAll_AllEmployeesByCompanyResult `xml:"Contract_GetAll_AllEmployeesByCompanyResult"`
		}
		Fault *Fault
	}
}

type Contract_GetAll_AllEmployeesByCompanyResult struct {
	XMLName                    xml.Name               `xml:"Contract_GetAll_AllEmployeesByCompanyResult"`
	EmployeeContractItemGlobal []EmployeeContractItem `xml:"EmployeeContractItemGlobal"`
}

type EmployeeContractItem struct {
	XMLName           xml.Name          `xml:"EmployeeContractItemGlobal"`
	EmployeeID        int64             `xml:"EmployeeId"`
	EmployeeContracts EmployeeContracts `xml:"EmployeeContracts"`
}

type EmployeeContracts struct {
	XMLName          xml.Name           `xml:"EmployeeContracts"`
	EmployeeContract []EmployeeContract `xml:"EmployeeContract"`
}

type EmployeeContract struct {
	XMLName                 xml.Name `xml:"EmployeeContract"`
	ContractID              int64    `xml:"ContractID"`
	CreationDate            string   `xml:"CreationDate"`
	StartDate               string   `xml:"StartDate"`
	TrialPeriod             string   `xml:"TrialPeriod"`
	EndDate                 string   `xml:"EndDate"`
	EmploymentType          int64    `xml:"EmployementType"`
	EmploymentSequenceTaxID int64    `xml:"EmploymentSequenceTaxId"`
	Indefinite              bool     `xml:"Indefinite"`
	PhaseClassification     int64    `xml:"PhaseClassification"`
	WrittenContract         bool     `xml:"WrittenContract"`
	HoursPerWeek            int64    `xml:"HoursPerWeek"`
}

type Contract_GetAll_AllEmployeesByCompany struct {
	XMLName   xml.Name `xml:"Contract_GetAll_AllEmployeesByCompany"`
	XMLNS     string   `xml:"xmlns,attr"`
	CompanyID int64    `xml:"CompanyID"`
}

func (service *Service) GetContracts(companyID int64) (*[]EmployeeContractItem, *errortools.Error) {
	xmlns := "https://api.nmbrs.nl/soap/v3/EmployeeService"
	bodyModel := service.GetSOAPEnvelope(xmlns, Contract_GetAll_AllEmployeesByCompany{
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

	r := Contract_GetAll_AllEmployeesByCompanyResponse{}

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

	return &r.SoapBody.Response.Contract_GetAll_AllEmployeesByCompanyResult.EmployeeContractItemGlobal, nil
}
