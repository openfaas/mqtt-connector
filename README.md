## mqtt-connector

[![Build Status](https://travis-ci.com/openfaas-incubator/mqtt-connector.svg?branch=master)](https://travis-ci.com/openfaas-incubator/mqtt-connector)

## Status

This is an MQTT connector for OpenFaaS. Once configured and deployed it will deliver messages from selected topics to OpenFaaS functions.

There are various other connectors available for OpenFaaS which form ["triggers"](https://docs.openfaas.com/reference/triggers/) for event-driven architectures.

Prior work:

This is inspired by prior work by [Alex Ellis](https://www.alexellis.io): [Collect, plot and analyse sensor readings from your IoT devices with OpenFaaS](https://github.com/alexellis/iot-sensors-mqtt-openfaas)

Component parts:

* [connector-sdk](https://github.com/openfaas-incubator/connector-sdk/blob/) from OpenFaaS
* The Eclipse provides [a test broker](https://mosquitto.org)
* Eclipse's [paho.mqtt.golang package](https://github.com/eclipse/paho.mqtt.golang) provides the connection to MQTT.

## Deploy in-cluster with Kubernetes

See [helm chart](chart/mqtt-connector) for deployment instructions. Then continue at "Test the connector".

```sh
TAG=0.2.0 make build push
```

## Deploy out of cluster

```sh
go build

export TOPIC=""
export GATEWAY_PASS=""

./mqtt-connector --gateway http://192.168.0.35:8080 \
  --gw-password $GATEWAY_PASS \
  --broker tcp://test.mosquitto.org:1883 \
  --topic $TOPIC
```

## Test the connector

Annotate a function with the annotation `topic: $TOPIC` <- where `$TOPIC` is the MQTT topic you care about.

```sh
2019/12/03 16:43:26 Topic: topic        Broker: tcp://test.mosquitto.org:1883
2019/12/03 16:43:29 Invoking (http://192.168.0.35:8080) on topic: "topic", value: "{\"sensor\": \"s1\", \"humidity\": \"52.09\", \"temp\": \"23.200\", \"ip\": \"192.168.0.40\", \"vdd33\": \"65535\", \"rssi\": -45}"
2019/12/03 16:43:29 Invoke function: print-out
Send: "{\"sensor\": \"s1\", \"humidity\": \"52.09\", \"temp\": \"23.200\", \"ip\": \"192.168.0.40\", \"vdd33\": \"65535\", \"rssi\": -45}"
2019/12/03 16:43:29 connector-sdk got result: [200] topic => print-out (24) bytes
[200] topic => print-out
{"temperature":"23.200"}
2019/12/03 16:43:29 tester got result: [200] topic => print-out (24) bytes
```

This data was generated on the topic `topic` by my NodeMCU device which publishes sensor data.

A `node12` function named `print-out` returned the temperature as reported.

## Custom builds

```sh
export NAMESPACE="alexellis2" # Or set your own registry/username
TAG=0.1.1 make build push
```

## License

MIT
