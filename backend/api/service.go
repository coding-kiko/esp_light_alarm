package api

import (
	"strconv"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type Service interface {
	turnOn() error
	turnOff() error
	cancelAlarm() error
	setAlarm(hour, min int) error
}

type service struct {
	mqttTopic  string
	mqttClient mqtt.Client
}

func (s *service) setAlarm(hour, min int) error {
	seconds := calculateSecondsUntil(hour, min)

	ok := s.mqttClient.Publish(s.mqttTopic, 0, false, strconv.Itoa(seconds))
	if ok.Wait() && ok.Error() != nil {
		return ok.Error()
	}
	return nil
}

func (s *service) cancelAlarm() error {
	ok := s.mqttClient.Publish(s.mqttTopic, 0, false, "CANCEL")
	if ok.Wait() && ok.Error() != nil {
		return ok.Error()
	}
	return nil
}

func (s *service) turnOn() error {
	ok := s.mqttClient.Publish(s.mqttTopic, 0, false, "ON")
	if ok.Wait() && ok.Error() != nil {
		return ok.Error()
	}
	return nil
}

func (s *service) turnOff() error {
	ok := s.mqttClient.Publish(s.mqttTopic, 0, false, "OFF")
	if ok.Wait() && ok.Error() != nil {
		return ok.Error()
	}
	return nil
}

func NewService(mc mqtt.Client, topic string) Service {
	return &service{
		mqttClient: mc,
		mqttTopic:  topic,
	}
}

// calculate time until the next hour:minute, takes into account if it falls in the day after
// returns result in seconds            REFACTOR THIS SHI UGLY
func calculateSecondsUntil(hour, minute int) int {
	now := time.Now()
	if now.Hour() > hour {
		return tomorrow(hour, minute)
	}
	if now.Hour() < hour {
		return today(hour, minute)
	}
	if now.Minute() >= minute {
		return tomorrow(hour, minute)
	}
	return today(hour, minute)
}

func today(hour, min int) int {
	return 60 * (((hour - time.Now().Hour()) * 60) + min - time.Now().Minute())
}

func tomorrow(hour, min int) int {
	return (24 * 60 * 60) + today(hour, min)
}
