// Copyright (c) OpenFaaS Author(s) 2019. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/openfaas-incubator/connector-sdk/types"
	"github.com/openfaas/faas-provider/auth"
)

func main() {

	var gatewayUsername, gatewayPassword, gatewayFlag string
	var trimChannelKey bool

	flag.StringVar(&gatewayUsername, "gw-username", "", "Username for the gateway")
	flag.StringVar(&gatewayPassword, "gw-password", "", "Password for gateway")
	flag.StringVar(&gatewayFlag, "gateway", "", "gateway")
	flag.BoolVar(&trimChannelKey, "trim-channel-key", false, "Trim channel key when using emitter.io MQTT broker")

	topic := flag.String("topic", "", "The topic name to/from which to publish/subscribe")
	broker := flag.String("broker", "tcp://iot.eclipse.org:1883", "The broker URI. ex: tcp://10.10.1.1:1883")
	password := flag.String("password", "", "The password (optional)")
	user := flag.String("user", "", "The User (optional)")
	id := flag.String("id", "testgoid", "The ClientID (optional)")
	cleansess := flag.Bool("clean", false, "Set Clean Session (default false)")
	qos := flag.Int("qos", 0, "The Quality of Service 0,1,2 (default 0)")
	num := flag.Int("num", 1000, "The number of messages to publish or subscribe (default 1)")

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

	gatewayURL := os.Getenv("gateway_url")

	if len(gatewayFlag) > 0 {
		gatewayURL = gatewayFlag
	}

	if len(gatewayURL) == 0 {
		log.Panicln(`a value must be set for env "gatewayURL" or via the -gateway flag for your OpenFaaS gateway`)
		return
	}

	config := &types.ControllerConfig{
		RebuildInterval:   time.Millisecond * 1000,
		GatewayURL:        gatewayURL,
		PrintResponse:     true,
		PrintResponseBody: true,
	}

	log.Printf("Topic: %s\tBroker: %s\n", *topic, *broker)

	controller := types.NewController(creds, config)

	receiver := ResponseReceiver{}
	controller.Subscribe(&receiver)

	controller.BeginMapBuilder()

	opts := MQTT.NewClientOptions()
	opts.AddBroker(*broker)
	opts.SetClientID(*id)
	opts.SetUsername(*user)
	opts.SetPassword(*password)
	opts.SetCleanSession(*cleansess)

	receiveCount := 0
	choke := make(chan [2]string)

	opts.SetDefaultPublishHandler(func(client MQTT.Client, msg MQTT.Message) {
		choke <- [2]string{msg.Topic(), string(msg.Payload())}
	})

	client := MQTT.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	if token := client.Subscribe(*topic, byte(*qos), nil); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}

	for receiveCount < *num {
		incoming := <-choke

		topic := incoming[0]
		data := []byte(incoming[1])

		if trimChannelKey {
			log.Printf("Topic before trim: %s\n", topic)
			index := strings.Index(topic, "/")
			topic = topic[index+1:]
		}

		log.Printf("Invoking (%s) on topic: %q, value: %q\n", gatewayURL, topic, data)

		controller.Invoke(topic, &data)

		receiveCount++
	}

	client.Disconnect(250)
	fmt.Println("Sample Subscriber Disconnected")
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
