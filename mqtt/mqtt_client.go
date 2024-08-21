package mqtt

import (
	paho_mqtt_golang "github.com/eclipse/paho.mqtt.golang"
	"log"
)

type MqttClient struct {
	client paho_mqtt_golang.Client
}

func NewMqttClient(broker string, username string, password string, clientId string) *MqttClient {
	opts := paho_mqtt_golang.NewClientOptions()
	opts.AddBroker(broker)

	opts.SetClientID(clientId)
	opts.SetAutoReconnect(true)
	opts.SetCleanSession(false)
	opts.SetUsername(username)
	opts.SetPassword(password)
	opts.SetOnConnectHandler(func(client paho_mqtt_golang.Client) {
		log.Println("ReConnect to broker")
	})
	//opts.SetDefaultPublishHandler(messageHandler)

	mqttClient := &MqttClient{}
	mqttClient.client = paho_mqtt_golang.NewClient(opts)
	return mqttClient
}

func (this *MqttClient) Connect() paho_mqtt_golang.Token {
	return this.client.Connect()
}

func (this *MqttClient) Disconnect() {
	this.client.Disconnect(250)
}

func (this *MqttClient) Subscribe(topic string, qos byte, callback paho_mqtt_golang.MessageHandler) paho_mqtt_golang.Token {
	return this.client.Subscribe(topic, qos, callback)
}

func (this *MqttClient) Publish(topic string, qos byte, retained bool, payload interface{}) paho_mqtt_golang.Token {
	return this.client.Publish(topic, qos, retained, payload)
}
