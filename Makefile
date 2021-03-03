.PHONY: build push manifest test verify-codegen charts
TAG?=latest
NAMESPACE?=openfaas

# docker manifest command will work with Docker CLI 18.03 or newer
# but for now it's still experimental feature so we need to enable that
export DOCKER_CLI_EXPERIMENTAL=enabled

build:
	docker build -t $(NAMESPACE)/mqtt-connector:$(TAG)-amd64 . -f Dockerfile
	docker build --build-arg OPTS="GOARCH=arm64" -t $(NAMESPACE)/mqtt-connector:$(TAG)-arm64 . -f Dockerfile
	docker build --build-arg OPTS="GOARCH=arm GOARM=6" -t $(NAMESPACE)/mqtt-connector:$(TAG)-armhf . -f Dockerfile

push:
	docker push $(NAMESPACE)/mqtt-connector:$(TAG)-amd64
	docker push $(NAMESPACE)/mqtt-connector:$(TAG)-arm64
	docker push $(NAMESPACE)/mqtt-connector:$(TAG)-armhf

manifest:
	docker manifest create --amend $(NAMESPACE)/mqtt-connector:$(TAG) \
		$(NAMESPACE)/mqtt-connector:$(TAG)-amd64 \
		$(NAMESPACE)/mqtt-connector:$(TAG)-arm64 \
		$(NAMESPACE)/mqtt-connector:$(TAG)-armhf
	docker manifest annotate $(NAMESPACE)/mqtt-connector:$(TAG) $(NAMESPACE)/mqtt-connector:$(TAG)-arm64 --os linux --arch arm64
	docker manifest annotate $(NAMESPACE)/mqtt-connector:$(TAG) $(NAMESPACE)/mqtt-connector:$(TAG)-armhf --os linux --arch arm --variant v6
	docker manifest push -p $(NAMESPACE)/mqtt-connector:$(TAG)

test:
	go test ./...

charts:
	cd chart && helm package mqtt-connector
	mv chart/*.tgz docs/
	helm repo index docs --url https://openfaas.github.io/mqtt-connector/ --merge ./docs/index.yaml
