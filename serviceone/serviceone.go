package main

func main() {

	ch := make(chan bool)
	<-ch

}

func init() {
	c := ServiceOneConfigManager{}
	c.Init()

	e := ServiceOneEurekaManager{}
	e.EurekaService = c.EurekaService()
	e.RegistrationTicket = c.RegistrationTicket()

	e.SendRegistrationOrFail()
	e.ScheduleHeartBeat(c.ServiceName(), 5)
}
