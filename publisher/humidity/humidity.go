package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/BENSARI-Fathi/mqtt/form"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var (
	clientID      = "humidity-sensor"
	server        = "broker.emqx.io"
	port          = 1883
	humidityTopic = "/sensor/humidity"
)

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	fmt.Printf("connection lost %v", err)
}

var onConnectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	fmt.Println("[+] connexion has been established.")
}

func randomHumGen() float32 {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	return r1.Float32() * 70
}

func main() {
	opts := mqtt.NewClientOptions()
	opts.SetOrderMatters(false)
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", server, port))
	opts.SetClientID(clientID)
	opts.SetOnConnectHandler(onConnectHandler)
	opts.SetConnectionLostHandler(connectLostHandler)

	c := mqtt.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	c.Subscribe(humidityTopic, 0, nil)
	payload := &form.HumidityForm{}
	for {
		payload.Device = clientID
		payload.Value = randomHumGen()
		data, err := json.Marshal(payload)
		if err != nil {
			fmt.Println("[-] Error while marshaling data:", err.Error())
		}
		t := c.Publish(humidityTopic, 0, false, data)
		go func() {
			_ = t.Wait() // Can also use '<-t.Done()' in releases > 1.2.0
			if t.Error() != nil {
				log.Fatal(t.Error())
			}
		}()
		time.Sleep(time.Second * 1)
	}
}
