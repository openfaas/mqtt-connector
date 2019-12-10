# OpenFaaS MQTT Connector chart

## Installation

```sh
git clone https://github.com/openfaas-incubator/mqtt-connector
cd mqtt-connector/chart

# Template and apply:
helm template -n openfaas --namespace openfaas mqtt-connector/ | kubectl apply -f -
```

You can watch the Connector logs to see it invoke your functions:

```
kubectl logs deployment.apps/openfaas-mqtt-connector -n openfaas -f
```

## Configuration

Configure via `values.yaml`.

<!-- | Parameter                | Description                                                                            | Default                        |
| ------------------------ | -------------------------------------------------------------------------------------- | ------------------------------ |
| `upstream_timeout`       | Maximum timeout for upstream function call, must be a Go formatted duration string.    | `30s`                          | -->
