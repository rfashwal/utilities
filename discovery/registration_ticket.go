package discovery

import (
	"bytes"
	"fmt"
	"github.com/xdevices/utilities/net"
	"strings"
)

const (
	STATUS_UP = "UP"
)

type RegistrationTicket struct {
	Instance Instance `json:"instance"`
}

type Instance struct {
	InstanceId       string         `json:"instanceId"`
	HostName         string         `json:"hostName"`
	App              string         `json:"app"`
	IpAddr           string         `json:"ipAddr"`
	VipAddress       string         `json:"vipAddress"`
	SecureVipAddress NullableString `json:"secureVipAddress"`
	Status           string         `json:"status"`
	Port             Port           `json:"port"`
	SecurePort       Port           `json:"securePort"`
	HomePageUrl      NullableString `json:"homePageUrl"`
	StatusPageUrl    string         `json:"statusPageUrl"`
	HealthCheckUrl   NullableString `json:"healthCheckUrl"`
	DataCenterInfo   DataCenterInfo `json:"dataCenterInfo"`
	LeaseInfo        LeaseInfo      `json:"leaseInfo"`
}

type DataCenterInfo struct {
	Class string `json:"@class"`
	Name  string `json:"name"`
}

type LeaseInfo struct {
	EvictionDurationInSecs int `json:"durationInSecs"`
}

type Port struct {
	Port    int  `json:"$"`
	Enabled bool `json:"@enabled"`
}

type NullableString string

func (c NullableString) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer
	if len(string(c)) == 0 {
		buf.WriteString(`null`)
	} else {
		buf.WriteString(`"` + string(c) + `"`) // add double quation mark as json format required
	}
	return buf.Bytes(), nil
}

func (c *NullableString) UnmarshalJSON(in []byte) error {
	str := string(in)
	if str == `null` {
		*c = ""
		return nil
	}
	res := NullableString(str)
	if len(res) >= 2 {
		res = res[1 : len(res)-1] // remove the wrapped qutation
	}
	*c = res
	return nil
}

func BuildRegistrationTicket(appName string, port int, evictionDurationInSeconds int, ignoreLoopback bool) *RegistrationTicket {

	ipAddress, err := net.GetIP(ignoreLoopback)

	if err != nil {
		panic(fmt.Sprintf("Could not obtain IP adress: %s", err))
	}

	instanceId := fmt.Sprintf("%s:%s:%d", ipAddress, strings.ToUpper(appName), port)
	identifier := strings.ToUpper(appName)

	instance := Instance{
		InstanceId:       instanceId,
		HostName:         ipAddress,
		App:              identifier,
		IpAddr:           ipAddress,
		VipAddress:       identifier,
		SecureVipAddress: "",
		Status:           STATUS_UP,
		Port: Port{
			Port:    port,
			Enabled: true,
		},
		SecurePort: Port{
			Port:    8443,
			Enabled: false,
		},
		HomePageUrl:    "",
		HealthCheckUrl: "",
		DataCenterInfo: DataCenterInfo{
			Class: "com.netflix.appinfo.MyDataCenterInfo",
			Name:  "MyOwn",
		},
		LeaseInfo: LeaseInfo{
			EvictionDurationInSecs: evictionDurationInSeconds,
		},
	}

	ticket := &RegistrationTicket{
		Instance: instance,
	}

	return ticket
}
