package nmbrs

import (
	"encoding/xml"
	"io/ioutil"
	"net/http"

	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
)

type CostCenter_GetListResponse struct {
	XMLName  xml.Name `xml:"http://www.w3.org/2003/05/soap-envelope Envelope"`
	SoapBody struct {
		XMLName  xml.Name `xml:"Body"`
		Response *struct {
			XMLName xml.Name                 `xml:"CostCenter_GetListResponse"`
			Result  CostCenter_GetListResult `xml:"CostCenter_GetListResult"`
		}
		Fault *Fault
	}
}

type CostCenter_GetListResult struct {
	XMLName    xml.Name      `xml:"CostCenter_GetListResult"`
	CostCenter *[]CostCenter `xml:"CostCenter"`
}

type CostCenter struct {
	XMLName     xml.Name `xml:"CostCenter"`
	Code        int64    `xml:"Code"`
	Description string   `xml:"Description"`
	ID          int64    `xml:"Id"`
}

type CostCenter_GetList struct {
	XMLName   xml.Name `xml:"CostCenter_GetList"`
	XMLNS     string   `xml:"xmlns,attr"`
	CompanyID int64    `xml:"CompanyId"`
}

func (service *Service) GetCostCenters(companyID int64) (*[]CostCenter, *errortools.Error) {
	xmlns := "https://api.nmbrs.nl/soap/v3/CompanyService"
	bodyModel := service.GetSOAPEnvelope(xmlns, CostCenter_GetList{
		XMLNS:     xmlns,
		CompanyID: companyID,
	})

	requestConfig := go_http.RequestConfig{
		Method:    http.MethodPost,
		URL:       service.url("CompanyService.asmx"),
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

	r := CostCenter_GetListResponse{}

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

	return r.SoapBody.Response.Result.CostCenter, nil
}
