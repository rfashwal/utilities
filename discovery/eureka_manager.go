package discovery

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/carlescere/scheduler"
	"github.com/labstack/gommon/log"
)

const (
	RegistrationTicketCannotBeEmptyMessage     = "registration ticket to Eureka Discovery cannot be empty"
	RegistrationTicketMarshallingFailedMessage = "could not marshal RegistrationTicket. %s"
	EurekaUrlCannotBeEmptyMessage              = "eurekaUrl cannot be empty"
	CannotCreateNewRequestMessage              = "cannot create new request. %s"
)

type Manager struct {
	RegistrationTicket *RegistrationTicket
	EurekaService      string
}

type EurekaResponse struct {
	Status     string
	StatusCode int
	Body       string
}

type buildUrl func(*Manager) string

func (s *Manager) SendRegistrationOrFail() {
	response, err := s.register()
	if err != nil {
		panic(fmt.Sprintf("cannot register at eureka, check your EUREKA_SERVICE variable. [%s]", err))
	}

	log.Info(fmt.Sprintf("eureka response: %s", response.Status))
}

func (s *Manager) ScheduleHeartBeat(service string, intervalInSec int) {
	_, err := scheduler.Every(intervalInSec).Seconds().NotImmediately().Run(func() {
		s.sendHeartBeat(service)
	})

	if err != nil {
		panic(fmt.Sprintf("could not run job. [%s]", err))
	}
}

func (s *Manager) sendHeartBeat(service string) {
	eurekaResponse, err := s.renew()
	if err != nil {
		panic(fmt.Sprintf("could not send heartbeat. [%s]", err))
	}
	log.Info(fmt.Sprintf("[%s]: heartbeat response: [%s]", service, eurekaResponse.Status))
}

func (s *Manager) register() (*EurekaResponse, error) {
	return s.doRequest(http.MethodPost, buildRegisterUrl)
}

func (s *Manager) renew() (*EurekaResponse, error) {
	return s.doRequest(http.MethodPut, buildRenewalUrl)
}

func (s *Manager) doRequest(method string, fn buildUrl) (*EurekaResponse, error) {
	if s.RegistrationTicket == nil {
		return nil, errors.New(RegistrationTicketCannotBeEmptyMessage)
	}

	if len(s.EurekaService) == 0 {
		return nil, errors.New(EurekaUrlCannotBeEmptyMessage)
	}

	requestUrl := fn(s)

	marshal, err := json.Marshal(*s.RegistrationTicket)
	if err != nil {
		panic(fmt.Sprintf(RegistrationTicketMarshallingFailedMessage, err))
	}

	request, err := http.NewRequest(method, requestUrl, bytes.NewBufferString(string(marshal)))
	if err != nil {
		panic(fmt.Sprintf(CannotCreateNewRequestMessage, err))
	}

	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)

	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	eurekaResponse := &EurekaResponse{
		Status:     response.Status,
		StatusCode: response.StatusCode,
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return eurekaResponse, err
	}

	eurekaResponse.Body = string(body)
	return eurekaResponse, nil
}

func buildRegisterUrl(s *Manager) string {
	// NOTE: %s/eureka/apps/SERVICENAME
	return fmt.Sprintf("%s/eureka/apps/%s", s.EurekaService, s.RegistrationTicket.Instance.App)
}

func buildRenewalUrl(s *Manager) string {
	// NOTE: %s/eureka/apps/SERVICENAME/192.168.1.49:SERVICENAME:9000
	return fmt.Sprintf("%s/eureka/apps/%s/%s", s.EurekaService, s.RegistrationTicket.Instance.App, s.RegistrationTicket.Instance.InstanceId)
}
