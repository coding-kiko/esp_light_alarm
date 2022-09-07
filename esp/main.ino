#include <ESP8266WiFi.h>
#include <WiFiClient.h>
#include <ESP8266HTTPClient.h>
#include <ArduinoMqttClient.h>

WiFiClient wifiClient;
MqttClient mqttClient(wifiClient);

int pin =0;
const char* ssid = "Apto307-2.4GHz";
const char*  pwd = "Calixtos";
const char* topic = "bedroom/esp_light";
const char* broker_addr = "192.168.1.10";
const int broker_port = 1883;
char* buffer = 0;
unsigned long startAlarmMillis;
unsigned long period;
bool startedAlarm = false;

void setup() {
    pinMode(pin, OUTPUT);
    digitalWrite(pin, HIGH); // turn down on reset (its inverted)

    // connect to wifi
    WiFi.begin(ssid, pwd);
    while (WiFi.status() != WL_CONNECTED) { // wait until connection is established
        delay(1000);
        if (WiFi.status() == 1) {
          break;
        }
    }

    // connect to mqtt broker
    if (!mqttClient.connect(broker_addr, broker_port)) {
      mqttClient.connectError();
      while (1);
    }

    // set the message receive callback
    mqttClient.onMessage(onMqttMessage);
    mqttClient.subscribe(topic);

    // [DEBUG] show that connections are established by blinking
    digitalWrite(pin, LOW);
    delay(500);
    digitalWrite(pin, HIGH);
    // [DEBUG]
}

void loop() {
  mqttClient.poll();
  if (startedAlarm) {
    if (millis() - startAlarmMillis >= period) {
      digitalWrite(pin, LOW);
      startedAlarm = false;
    }
  }
}

void onMqttMessage(int messageSize) {
  int i;

  // if buffer is not clean, empty it
  if (buffer != 0) {
    delete [] buffer;
  }
  // declare new buffer with messageSize size
  buffer = new char [messageSize];
  // fill new buffer with message comming from mosquitto
  for (i=0; i<messageSize; i++) {
    buffer[i] = (char)mqttClient.read(); 
  }
  buffer[i] = '\0';

  if (strcmp(buffer,(char*)"ON") == 0) {
    digitalWrite(pin, LOW);
    return;
  }
  if (strcmp(buffer,(char*)"OFF") == 0) {
    digitalWrite(pin, HIGH);
    return;
  }
  if (strcmp(buffer,(char*)"CANCEL") == 0) {
    startedAlarm = false;
    return;
  }

  int seconds = atoi(buffer);
  period = 1000 * seconds;
  startAlarmMillis = millis();
  startedAlarm = true;
}
