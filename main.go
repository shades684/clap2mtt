package main

//cat /proc/sys/kernel/random/uuid

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"time"

	"clap2mqtt/clapping"
	"clap2mqtt/detection"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/google/uuid"
	"github.com/stianeikeland/go-rpio"
)

type Config struct {
	Broker string `json:"broker"`
	Uuid   string `json:"uuid"`
	Pin    int    `json:"pin"`
}

func exists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func configure() Config {
	const fileName = "config.json"
	var config Config

	if !exists(fileName) {
		config.Broker = "tcp://127.0.0.1:1883"
		config.Uuid = uuid.NewString()
		config.Pin = 18

		b, _ := json.MarshalIndent(config, "", " ")
		_ = ioutil.WriteFile(fileName, b, 0644)

	} else {
		jsonFile, _ := os.Open(fileName)
		defer jsonFile.Close()

		byteValue, _ := ioutil.ReadAll(jsonFile)
		json.Unmarshal(byteValue, &config)
	}

	return config
}

func setup() (rpio.Pin, func(int)) {
	config := configure()
	objectId := config.Uuid + "_" + strconv.Itoa(config.Pin)

	rpio.Open()
	defer rpio.Close()

	pin := rpio.Pin(config.Pin)
	pin.Input()

	options := mqtt.NewClientOptions()
	options.AddBroker(config.Broker)
	client := mqtt.NewClient(options)
	client.Connect()

	defer client.Disconnect(0)

	message := `{
        "value": 3,
        "name": "Pi Clap",
        "icon": "mdi:gesture-double-tap",
        "unique_id: "sensor.%s",
        "device": {
            "identifiers": ["sensor.%s"]
            "name": "Raspberry Pi Clapper", 
            "model": "DIY",
            "manufacturer": "DIY"
        }
    }`

	token := client.Publish("homeassistant/sensor/"+objectId+"/config", 0, true, fmt.Sprintf(message, objectId, objectId))
	token.Wait()

	return pin, func(claps int) {
		token := client.Publish("clap2mqtt/"+objectId, 0, true, claps)
		token.Wait()

		time.Sleep(500 * time.Millisecond)

		token = client.Publish("clap2mqtt/"+objectId, 0, true, 0)
		token.Wait()
	}
}

func main() {
	pin, publish := setup()
	clapping := clapping.NewClapping()

	var soundDetection *detection.Detection = nil

	for {
		signal := pin.Read() == rpio.Low

		if soundDetection == nil && signal {
			soundDetection = detection.NewDetection()
		} else {
			soundDetection.Update(signal)

			if soundDetection.HasStopped() {
				clapping.AddDetection(*soundDetection)
				soundDetection = nil
			}
		}

		if clapping.HasStopped() {
			publish(clapping.Count())
			clapping.Reset()
		}

		time.Sleep(time.Millisecond)
	}
}
