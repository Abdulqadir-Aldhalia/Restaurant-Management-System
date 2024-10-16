package services

import (
	"fmt"
	"net/http"
	"server-side/model"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

const (
	broker   = "tcp://broker.emqx.io:1883"
	clientID = "go_mqtt_client"
	topic    = "esp32/notify"
)

func NotifyServer(w http.ResponseWriter, r *http.Request) {
	table_id := r.FormValue("table_id")
	var table model.Tables
	err := UpdateById("tables", table_id, map[string]interface{}{"is_needs_service": true}, &table, nil)
	if err != nil {
		SendErrorResponse(w, err)
		return
	}
	SendJsonResponse(w, http.StatusAccepted, nil)
}

func NotifyController(value int) {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(broker)
	opts.SetClientID(clientID)

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		fmt.Println("Failed to connect:", token.Error())
		return
	}
	fmt.Println("Connected to MQTT broker")

	// Publish message to the topic every 5 seconds
	for {
		text := "Hello from Go Server!"
		token := client.Publish(topic, 0, false, value)
		token.Wait()

		fmt.Printf("Published message: %s\n", text)
		time.Sleep(5 * time.Second)
	}

	// Disconnect the client (optional, if you plan to stop it)
	client.Disconnect(250)
}
