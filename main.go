package main

import (
	"fmt"
	"log"
	"machine"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var (
	msgChan = make(chan [2]string)
)

func main() {
	led := machine.Pin(2)
	led.Configure(machine.PinConfig{Mode: machine.PinOutput})
	opts := mqtt.NewClientOptions()
	opts.AddBroker("192.168.1.10:1883")
	opts.SetClientID("esp_alarm_client")
	opts.SetDefaultPublishHandler(handler)

	client := mqtt.NewClient(opts)

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	if token := client.Subscribe("bedroom/esp_light", 0, nil); token.Wait() && token.Error() != nil {
		log.Fatalln(token.Error())
	}

	fmt.Println("starting infinite consumer")

	for {
		incoming := <-msgChan
		if incoming[1] == "ON" {
			led.Low()
		} else {
			led.High()
		}
	}
}

func handler(c mqtt.Client, msg mqtt.Message) {
	msgChan <- [2]string{msg.Topic(), string(msg.Payload())}
}
