## mqtt-connector

## Status

Early prototype.

## Usage

```sh
go build

export TOPIC=""
export GATEWAY_PASS=""

./mqtt-connector --gateway http://192.168.0.35:8080 \
  --gw-password $GATEWAY_PASS \
  --broker tcp://test.mosquitto.org:1883 \
  --topic $TOPIC
```

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
