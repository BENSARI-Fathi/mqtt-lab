package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var (
	clientID         = "Samsung-S20"
	server           = "broker.emqx.io"
	port             = 1883
	temperatureTopic = "/sensor/temperature"
	humidityTopic    = "/sensor/humidity"
	scheme           = "tcp"
)

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	fmt.Printf("connection lost %v", err)
}

var onConnectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	fmt.Println("[+] connexion has been established.")
}

var defaultPubMsgHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("[+] New message from topic: %s\nmessage content: %s\n", msg.Topic(), msg.Payload())
}

func main() {
	opts := mqtt.NewClientOptions()
	opts.SetOrderMatters(false)
	opts.AddBroker(fmt.Sprintf("%s://%s:%d", scheme, server, port))
	opts.SetClientID(clientID)
	opts.SetOnConnectHandler(onConnectHandler)
	opts.SetConnectionLostHandler(connectLostHandler)
	opts.SetDefaultPublishHandler(defaultPubMsgHandler)

	c := mqtt.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	c.Subscribe(temperatureTopic, 0, nil)
	c.Subscribe(humidityTopic, 0, nil)
	//Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	c.Disconnect(250)
	log.Println("Shutdown Client ...")

}
