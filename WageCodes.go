package nmbrs

import (
	"encoding/xml"

	errortools "github.com/leapforce-libraries/go_errortools"
	go_types "github.com/leapforce-libraries/go_types"
)

type WageCodeReports struct {
	XMLName xml.Name          `xml:"reports"`
	Report  *[]WageCodeReport `xml:"report"`
}

type WageCodeReport struct {
	XMLName    xml.Name      `xml:"report"`
	EmployeeID int64         `xml:"employeeid"`
	Period     int64         `xml:"period"`
	Year       int64         `xml:"year"`
	Lines      WageCodeLines `xml:"lines"`
}

type WageCodeLines struct {
	XMLName xml.Name       `xml:"lines"`
	Line    []WageCodeLine `xml:"line"`
}

type WageCodeLine struct {
	XMLName     xml.Name `xml:"line"`
	Code        int64    `xml:"code"`
	Description string   `xml:"description"`
	Value       float64  `xml:"value"`
}

type Reports_BackgroundResponse struct {
	XMLName  xml.Name `xml:"http://www.w3.org/2003/05/soap-envelope Envelope"`
	SoapBody struct {
		XMLName  xml.Name `xml:"Body"`
		Response *struct {
			XMLName xml.Name       `xml:",any"`
			Result  *go_types.GUID `xml:",any"`
		}
		Fault *Fault
	}
}

type Reports_Background struct {
	XMLName   xml.Name `xml:"Reports_Accountant_Company_CompanyWageComponentsPerPeriod_Background"`
	XMLNS     string   `xml:"xmlns,attr"`
	CompanyID int64    `xml:"companyId"`
	Year      int64    `xml:"year"`
}

func (service *Service) GetWageCodesByYear(companyID int64, year int64) (*[]WageCodeReport, *errortools.Error) {
	body := struct {
		XMLName   xml.Name `xml:"Reports_GetWageCodesByYear_Background"`
		XMLNS     string   `xml:"xmlns,attr"`
		CompanyID int64    `xml:"CompanyId"`
		Year      int64    `xml:"Year"`
	}{
		XMLNS:     "https://api.nmbrs.nl/soap/v3/ReportService",
		CompanyID: companyID,
		Year:      year,
	}

	// get Content
	wageCodeReports := WageCodeReports{}

	e := service.getReportsBackgroundTaskResult(body, &wageCodeReports)
	if e != nil {
		return nil, e
	}

	return wageCodeReports.Report, nil
}
