package nmbrs

import (
	"encoding/xml"
	"io/ioutil"
	"net/http"

	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
	go_types "github.com/leapforce-libraries/go_types"
)

type reports_BackgroundTaskResponse struct {
	XMLName  xml.Name `xml:"http://www.w3.org/2003/05/soap-envelope Envelope"`
	SoapBody struct {
		XMLName  xml.Name `xml:"Body"`
		Response *struct {
			//XMLName xml.Name       `xml:",any"`
			Result *go_types.GUID `xml:",any"`
		} `xml:",any"`
		Fault *Fault `xml:",any"`
	}
}

func (service *Service) getReportsBackgroundTask(body interface{}) (*go_types.GUID, *errortools.Error) {
	xmlns := "https://api.nmbrs.nl/soap/v3/ReportService"
	bodyModel := service.GetSOAPEnvelope(xmlns, body)

	requestConfig := go_http.RequestConfig{
		Method:    http.MethodPost,
		URL:       service.url("ReportService.asmx"),
		BodyModel: bodyModel,
	}
	_, response, e := service.httpRequest(&requestConfig)
	if e != nil {
		return nil, e
	}

	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, errortools.ErrorMessagef("ioutil.ReadAll error: %s ", err.Error())
	}
	defer response.Body.Close()

	r := reports_BackgroundTaskResponse{}

	err = xml.Unmarshal(responseBody, &r)
	if err != nil {
		return nil, errortools.ErrorMessagef("xml.Unmarshal error: %s ", err.Error())
	}

	if e != nil {
		if r.SoapBody.Fault != nil {
			e.SetMessage(r.SoapBody.Fault.Reason)
		}

		return nil, e
	}

	return r.SoapBody.Response.Result, nil
}
