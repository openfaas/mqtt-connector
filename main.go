// Copyright (c) OpenFaaS Author(s) 2019. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/openfaas/connector-sdk/types"
	"github.com/openfaas/faas-provider/auth"
)

func main() {
	var (
		gatewayUsername string
		gatewayPassword string
		gatewayFlag     string
		trimChannelKey  bool
		asyncInvoke     bool
		topic           string
		broker          string
	)

	flag.StringVar(&gatewayUsername, "gw-username", "admin", "Username for the gateway")
	flag.StringVar(&gatewayPassword, "gw-password", "", "Password for gateway")
	flag.StringVar(&gatewayFlag, "gateway", "", "gateway")
	flag.BoolVar(&trimChannelKey, "trim-channel-key", false, "Trim channel key when using emitter.io MQTT broker")
	flag.BoolVar(&asyncInvoke, "async-invoke", false, "Invoke via queueing using NATS and the function's async endpoint")
	flag.StringVar(&topic, "topic", "", "The topic name to/from which to publish/subscribe")
	flag.StringVar(&broker, "broker", "tcp://test.mosquitto.org:1883", "The broker URI. ex: tcp://test.mosquitto.org:1883")

	password := flag.String("password", "", "The password (optional)")
	user := flag.String("user", "", "The User (optional)")
	id := flag.String("id", "testgoid", "The ClientID (optional)")
	cleansess := flag.Bool("clean", false, "Set Clean Session (default false)")
	qos := flag.Int("qos", 0, "The Quality of Service 0,1,2 (default 0)")

	flag.Parse()

	var creds *auth.BasicAuthCredentials
	if len(gatewayPassword) > 0 {
		creds = &auth.BasicAuthCredentials{
			User:     gatewayUsername,
			Password: gatewayPassword,
		}
	} else {
		creds = types.GetCredentials()
	}

	contentType := "application/json"
	if v, exists := os.LookupEnv("content_type"); exists && len(v) > 0 {
		contentType = v
	}

	gatewayURL := os.Getenv("gateway_url")

	if len(gatewayFlag) > 0 {
		gatewayURL = gatewayFlag
	}

	if len(gatewayURL) == 0 {
		log.Panicln(`a value must be set for env "gatewayURL" or via the -gateway flag for your OpenFaaS gateway`)
		return
	}

	config := &types.ControllerConfig{
		RebuildInterval:          time.Millisecond * 1000,
		GatewayURL:               gatewayURL,
		PrintResponse:            true,
		PrintResponseBody:        true,
		TopicAnnotationDelimiter: ",",
		AsyncFunctionInvocation:  asyncInvoke,
		ContentType:              contentType,
	}

	log.Printf("Topic: %q\tBroker: %q\n", topic, broker)
	log.Printf("Gateway: %s\tAsync: %v\n", gatewayURL, asyncInvoke)

	controller := types.NewController(creds, config)

	receiver := ResponseReceiver{}
	controller.Subscribe(&receiver)

	controller.BeginMapBuilder()

	opts := MQTT.NewClientOptions()
	opts.AddBroker(broker)
	opts.SetClientID(*id)
	opts.SetUsername(*user)
	opts.SetPassword(*password)
	opts.SetCleanSession(*cleansess)

	receiveCount := 0
	msgCh := make(chan [2]string)

	opts.SetDefaultPublishHandler(func(client MQTT.Client, msg MQTT.Message) {
		log.Printf("Message incoming")
		msgCh <- [2]string{msg.Topic(), string(msg.Payload())}
	})

	opts.SetOnConnectHandler(func(client MQTT.Client) {
		log.Printf("Connected to %s", broker)

		if token := client.Subscribe(topic, byte(*qos), nil); token.Wait() && token.Error() != nil {
			fmt.Println(token.Error())
			os.Exit(1)
		}
		log.Printf("Subscribed to topic: %s", topic)
	})

	client := MQTT.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	log.Printf("Connection requested for broker: %s", broker)

	go func() {
		for {
			incoming := <-msgCh

			topic := incoming[0]
			data := []byte(incoming[1])

			if trimChannelKey {
				log.Printf("Topic before trim: %s", topic)
				index := strings.Index(topic, "/")
				topic = topic[index+1:]
			}

			log.Printf("Invoking (%s) on topic: %q, value: %q", gatewayURL, topic, data)

			controller.Invoke(topic, &data, http.Header{})

			receiveCount++
		}

		client.Disconnect(1250)
	}()

	select {}
}

// ResponseReceiver enables connector to receive results from the
// function invocation
type ResponseReceiver struct {
}

// Response is triggered by the controller when a message is
// received from the function invocation
func (ResponseReceiver) Response(res types.InvokerResponse) {
	if res.Error != nil {
		log.Printf("tester got error: %s", res.Error.Error())
	} else {
		log.Printf("tester got result: [%d] %s => %s (%d) bytes", res.Status, res.Topic, res.Function, len(*res.Body))
	}
}
