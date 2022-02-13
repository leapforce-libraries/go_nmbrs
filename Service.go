package nmbrs

import (
	"fmt"
	"net/http"

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

	accept := go_http.AcceptXML
	httpServiceConfig := go_http.ServiceConfig{
		Accept: &accept,
	}
	httpService, e := go_http.NewService(&httpServiceConfig)
	if e != nil {
		return nil, e
	}

	return &Service{
		apiToken:    serviceConfig.APIToken,
		username:    serviceConfig.Username,
		domain:      serviceConfig.Domain,
		httpService: httpService,
	}, nil
}

func (service *Service) url(path string) string {
	return fmt.Sprintf("%s/%s", apiURL, path)
}

func (service *Service) httpRequest(requestConfig *go_http.RequestConfig) (*http.Request, *http.Response, *errortools.Error) {
	// set content type header
	header := http.Header{}
	header.Set("Content-Type", "application/soap+xml; charset=utf-8")

	requestConfig.NonDefaultHeaders = &header

	return service.httpService.HTTPRequest(requestConfig)
}

func (service *Service) APIName() string {
	return apiName
}

func (service *Service) APIKey() string {
	return service.apiToken
}

func (service *Service) APICallCount() int64 {
	return service.httpService.RequestCount()
}

func (service *Service) APIReset() {
	service.httpService.ResetRequestCount()
}
