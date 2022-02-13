package nmbrs

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
	go_types "github.com/leapforce-libraries/go_types"
)

type Reports_BackgroundTask_ResultResultResponse struct {
	XMLName  xml.Name `xml:"http://www.w3.org/2003/05/soap-envelope Envelope"`
	SoapBody struct {
		XMLName  xml.Name `xml:"Body"`
		Response *struct {
			Result Reports_BackgroundTask_ResultResult `xml:",any"`
		} `xml:",any"`
		Fault *Fault
	}
}

type Reports_BackgroundTask_ResultResult struct {
	XMLName xml.Name      `xml:"Reports_BackgroundTask_ResultResult"`
	XMLNS   string        `xml:"xmlns,attr"`
	TaskID  go_types.GUID `xml:"TaskId"`
	Status  string        `xml:"Status"`
	Content *string       `xml:"Content"`
}

type Reports_BackgroundTask_Result struct {
	XMLName xml.Name      `xml:"Reports_BackgroundTask_Result"`
	XMLNS   string        `xml:"xmlns,attr"`
	TaskID  go_types.GUID `xml:"TaskId"`
}

type ContentReports struct {
	XMLName xml.Name         `xml:"reports"`
	Report  *[]ContentReport `xml:"report"`
}

type ContentReport struct {
	XMLName xml.Name `xml:"report"`
	Content interface{}
}

const (
	taskResultEnqueued  string = "Enqueued"
	taskResultExecuting string = "Executing"
	taskResultSuccess   string = "Success"
	timeout             int64  = 3600 // seconds
)

func (service *Service) getReportsBackgroundTaskResult(body interface{}, model interface{}) *errortools.Error {
	// get TaskID
	taskID, e := service.getReportsBackgroundTask(body)
	if e != nil {
		return e
	}
	if taskID == nil {
		return errortools.ErrorMessage("getReportsBackgroundTask returned nil")
	}

	// get content
	xmlns := "https://api.nmbrs.nl/soap/v3/ReportService"
	bodyModel := service.GetSOAPEnvelope(xmlns, Reports_BackgroundTask_Result{
		XMLNS:  xmlns,
		TaskID: *taskID,
	})

	requestConfig := go_http.RequestConfig{
		Method:    http.MethodPost,
		URL:       service.url("ReportService.asmx"),
		BodyModel: bodyModel,
	}

	now := time.Now()
	attempt := 1
	waitFor := 5
	waitForMax := 60

	for time.Since(now).Seconds() <= float64(timeout) {
		_, response, e := service.httpRequest(&requestConfig)
		if e != nil {
			return e
		}

		responseBody, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return errortools.ErrorMessagef("ioutil.ReadAll error: %s ", err.Error())
		}
		defer response.Body.Close()

		r := Reports_BackgroundTask_ResultResultResponse{}

		err = xml.Unmarshal(responseBody, &r)
		if err != nil {
			return errortools.ErrorMessagef("xml.Unmarshal error: %s ", err.Error())
		}

		if e != nil {
			if r.SoapBody.Fault != nil {
				e.SetMessage(r.SoapBody.Fault.Reason)
			}

			return e
		}

		if r.SoapBody.Response.Result.Status == taskResultSuccess {
			if r.SoapBody.Response.Result.Content == nil {
				return errortools.ErrorMessage("getReportsBackgroundTaskResult returned nil")
			}

			content := strings.ReplaceAll(strings.ReplaceAll(*r.SoapBody.Response.Result.Content, "&lt;", "<"), "&gt;", ">")

			err := xml.Unmarshal([]byte(content), model)
			if err != nil {
				return errortools.ErrorMessage(err)
			}

			return nil
		}

		if r.SoapBody.Response.Result.Status != taskResultEnqueued && r.SoapBody.Response.Result.Status != taskResultExecuting {
			return errortools.ErrorMessagef("Reports_BackgroundTask_Result returned status %s", r.SoapBody.Response.Result.Status)
		}

		fmt.Printf("%s: %v...\n", time.Now().Format("15:04:05"), r.SoapBody.Response.Result.Status)
		time.Sleep(time.Second * time.Duration(waitFor))
		if waitFor < waitForMax {
			waitFor *= 2
		}
		attempt++
	}

	return errortools.ErrorMessagef("Reports_BackgroundTask_Result is still executing after %v seconds", timeout)
}
