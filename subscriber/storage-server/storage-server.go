package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"gorm.io/gorm"

	database "github.com/BENSARI-Fathi/mqtt/db"
	"github.com/BENSARI-Fathi/mqtt/form"
)

var (
	clientID         = "Storage-server"
	server           = "broker.emqx.io"
	port             = 1883
	temperatureTopic = "/sensor/temperature"
	humidityTopic    = "/sensor/humidity"
	scheme           = "tcp"
	humidityQueue    *Queue
	temperatureQueue *Queue
	sum              float32
	db               *gorm.DB
	err              error
	humiditydata     *form.HumidityForm
	temperaturedata  *form.TemperatureForm
)

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	fmt.Printf("[-] connection lost %v", err)
}

var onConnectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	fmt.Println("[+] connexion has been established.")
}

var defaultPubMsgHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("[+] New message from topic: %s\nmessage content: %s\n", msg.Topic(), msg.Payload())
}

var humidityMsgHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	data := &form.HumidityForm{}
	json.Unmarshal(msg.Payload(), data)
	humidityQueue.Put(data)
}

var temperatureMsgHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	data := &form.TemperatureForm{}
	json.Unmarshal(msg.Payload(), data)
	temperatureQueue.Put(data)
}

func init() {
	// initialize the queue
	humidityQueue = NewQueue()
	temperatureQueue = NewQueue()
	// initialize the db
	db, err = database.NewSqliteCLient()
	if err != nil {
		log.Fatal("[-] failed to connect database")
	}
	db.AutoMigrate(&database.Humidity{})
	db.AutoMigrate(&database.Temperature{})
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
	// c.Subscribe(temperatureTopic, 0, nil)
	c.Subscribe(humidityTopic, 0, humidityMsgHandler)
	c.Subscribe(temperatureTopic, 0, temperatureMsgHandler)
	// humidity process and save into the db
	go func() {
		for {
			if humidityQueue.Len() == 10 {
				sum = 0
				for !humidityQueue.Empty() {
					humiditydata = humidityQueue.Get().Value.(*form.HumidityForm)
					sum += humiditydata.Value
				}
				log.Printf("[+] The humidity value is: %.2f%%\n", sum/10)
				humidity := &database.Humidity{
					Device: humiditydata.Device,
					Value:  sum / 10,
				}
				db.Create(humidity)
			}
		}
	}()
	// temperature process and save into the db
	go func() {
		for {
			if temperatureQueue.Len() == 10 {
				sum = 0
				for !temperatureQueue.Empty() {
					temperaturedata = temperatureQueue.Get().Value.(*form.TemperatureForm)
					sum += temperaturedata.Value
				}
				log.Printf("[+] The temperature value is: %.2f°C\n", sum/10)
				temperature := &database.Temperature{
					Device: temperaturedata.Device,
					Value:  sum / 10,
				}
				db.Create(temperature)
			}
		}
	}()
	//Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	c.Disconnect(250)
	log.Println("Shutdown Client ...")
}
