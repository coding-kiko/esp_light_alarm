package main

import (
	"log"
	"net/http"

	"github.com/coding-kiko/esp_light_alarm/api"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var (
	addr           = "0.0.0.0:8031"
	mqttBrokerAddr = "192.168.1.10:1883"
	mqttClientId   = "localhost:8031"
	topic          = "bedroom/esp_light"
)

func main() {
	mqttClient := NewMqttClient()
	defer mqttClient.Disconnect(1000)

	service := api.NewService(mqttClient, topic)
	handler := api.NewHandler(service)
	router := api.NewRouter(handler)

	log.Println("Started listening: 8031")
	err := http.ListenAndServe(addr, router)
	log.Fatal(err)
}

func NewMqttClient() mqtt.Client {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(mqttBrokerAddr)
	opts.SetClientID(mqttClientId)
	opts.OnConnectionLost = func(client mqtt.Client, err error) {
		log.Fatalln("Mqtt Connection Lost")
	}
	client := mqtt.NewClient(opts)
	token := client.Connect()
	if token.Wait() && token.Error() != nil {
		log.Fatalln(token.Error())
	}

	return client
}
