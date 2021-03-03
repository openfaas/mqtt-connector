## mqtt-connector

This is an MQTT connector for OpenFaaS.

[![build](https://github.com/openfaas/mqtt-connector/actions/workflows/build.yaml/badge.svg)](https://github.com/openfaas/mqtt-connector/actions/workflows/build.yaml)

## About

Once configured and deployed it will deliver messages from selected topics to OpenFaaS functions.

There are various other connectors available for OpenFaaS which form ["triggers"](https://docs.openfaas.com/reference/triggers/) for event-driven architectures.

Prior work:

This is inspired by prior work by [Alex Ellis](https://www.alexellis.io): [Collect, plot and analyse sensor readings from your IoT devices with OpenFaaS](https://github.com/alexellis/iot-sensors-mqtt-openfaas)

Component parts:

* [connector-sdk](https://github.com/openfaas/connector-sdk/blob/) from OpenFaaS
* The Eclipse provides [a test broker](https://mosquitto.org)
* Eclipse's [paho.mqtt.golang package](https://github.com/eclipse/paho.mqtt.golang) provides the connection to MQTT.

## Deploy in-cluster with Kubernetes

See [helm chart](https://github.com/openfaas/faas-netes/tree/master/chart/mqtt-connector) for deployment instructions. Then continue at "Test the connector".

```sh
export TAG=0.3.1

make build push
```

## Deploy out of cluster

```sh
go build

export GATEWAY_PASSWORD=""
export BROKER="tcp://test.mosquitto.org:1883"
export TOPIC="openfaas-sensor-data"

./mqtt-connector --gateway http://127.0.0.1:8080 \
  --broker $BROKER \
  --gw-username admin \
  --gw-password $GATEWAY_PASSWORD \
  --topic $TOPIC
```

Deploy a function:

```bash
faas-cli deploy --name echo --image ghcr.io/openfaas/alpine:latest \
  --fprocess=cat \
  --annotation topic="openfaas-sensor-data"
````

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

## License

MIT
