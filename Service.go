package nmbrs

import (
	"fmt"
	"net/http"
	"text/template"

	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
)

const (
	apiName string = "Nmbrs"
	apiURL  string = "https://api.nmbrs.nl/soap/v3"
)

type Service struct {
	apiToken    string
	username    string
	domain      string
	httpService *go_http.Service
	templates   *template.Template
}

type ServiceConfig struct {
	APIToken string
	Username string
	Domain   string
}

func NewService(serviceConfig *ServiceConfig) (*Service, *errortools.Error) {
	if serviceConfig == nil {
		return nil, errortools.ErrorMessage("ServiceConfig must not be a nil pointer")
	}

	if serviceConfig.APIToken == "" {
		return nil, errortools.ErrorMessage("APIToken not provided")
	}

	if serviceConfig.Username == "" {
		return nil, errortools.ErrorMessage("Username not provided")
	}

	if serviceConfig.Domain == "" {
		return nil, errortools.ErrorMessage("Domain not provided")
	}

	httpService, e := go_http.NewService(&go_http.ServiceConfig{})
	if e != nil {
		return nil, e
	}

	return &Service{
		apiToken:    serviceConfig.APIToken,
		username:    serviceConfig.Username,
		domain:      serviceConfig.Domain,
		httpService: httpService,
		templates:   template.Must(template.ParseGlob("templates/*.xml")),
	}, nil
}

func (service *Service) get(requestConfig *go_http.RequestConfig) (*http.Request, *http.Response, *errortools.Error) {
	return service.httpRequest(http.MethodGet, requestConfig)
}

func (service *Service) post(requestConfig *go_http.RequestConfig) (*http.Request, *http.Response, *errortools.Error) {
	return service.httpRequest(http.MethodPost, requestConfig)
}

func (service *Service) url(path string) string {
	return fmt.Sprintf("%s/%s", apiURL, path)
}

func (service *Service) httpRequest(httpMethod string, requestConfig *go_http.RequestConfig) (*http.Request, *http.Response, *errortools.Error) {
	return nil, nil, nil
}