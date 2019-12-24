package config

import (
	"fmt"
	"github.com/xdevices/utilities/discovery"
	"github.com/xdevices/utilities/net"
	"os"
	"strconv"

	"github.com/labstack/gommon/log"

)

type Manager struct {
	httpPort                 int
	address                  string
	serviceName              string
	eurekaService            string
	hostName                 string
	registrationTicket       *discovery.RegistrationTicket
	serviceDurationInSeconds int
	ignoreLoopback           bool
}

func (c *Manager) Init() {

	if port, err := os.LookupEnv("HTTP_PORT"); !err {
		panic("Set HTTP_PORT and try again")
	} else {
		if port, e := strconv.Atoi(port); e != nil {
			panic("please set HTTP_PORT variable and try again")
		} else {
			c.httpPort = port
		}
	}

	if serviceName, err := os.LookupEnv("SERVICE_NAME"); !err {
		panic("please set SERVICE_NAME variable and try again")
	} else {
		c.serviceName = serviceName
	}

	if ignoreLoopback, err := os.LookupEnv("IGNORE_LOOPBACK"); !err {
		// set default behaviour in case variable not set
		c.ignoreLoopback = true
	} else {
		if result, err := strconv.ParseBool(ignoreLoopback); err != nil {
			log.Warn(fmt.Sprintf("could not convert %s to bool", ignoreLoopback))
			c.ignoreLoopback = true
		} else {
			c.ignoreLoopback = result
		}
	}

	address, err := net.GetIP(c.ignoreLoopback)
	if err != nil {
		panic(fmt.Sprintf("Could not resolve ip address: %s", err))
	}
	c.address = address

	hostname, err := os.Hostname()
	if err != nil {
		panic(fmt.Sprintf("Could not obtain hostname: %s", err))
	}
	c.hostName = hostname

	if serviceDuration, err := os.LookupEnv("SERVICE_DURATION_IN_SECONDS"); !err {
		c.serviceDurationInSeconds = 5
	} else {
		if port, e := strconv.Atoi(serviceDuration); e != nil {
			c.serviceDurationInSeconds = port
		}
	}

	if eurekaHost, err := os.LookupEnv("EUREKA_SERVICE"); !err {
		panic(fmt.Sprintf("set EUREKA_SERVICE variable and try again"))
	} else {
		c.eurekaService = eurekaHost
	}

	registrationTicket := discovery.BuildRegistrationTicket(c.serviceName, c.httpPort, c.serviceDurationInSeconds, c.ignoreLoopback)
	c.registrationTicket = registrationTicket
}

func (c *Manager) HttpPort() int {
	return c.httpPort
}

func (c *Manager) Address() string {
	return fmt.Sprintf("%s:%d", c.address, c.httpPort)
}

func (c *Manager) ServiceName() string {
	return c.serviceName
}

func (c *Manager) EurekaService() string {
	return c.eurekaService
}

func (c *Manager) Hostname() string {
	return c.hostName
}

func (c *Manager) RegistrationTicket() *discovery.RegistrationTicket {
	return c.registrationTicket
}

func (c *Manager) IgnoreLoopback() bool {
	return c.ignoreLoopback
}
