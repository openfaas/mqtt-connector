TAG?=latest-dev
NAMESPACE?=alexellis2
.PHONY: build push

build:
	docker build -t $(NAMESPACE)/mqtt-connector:$(TAG) .

push:
	docker push $(NAMESPACE)/mqtt-connector:$(TAG)
