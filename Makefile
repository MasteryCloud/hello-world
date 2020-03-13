.PHONY: clean build

export USERNAME = masterycloud
export IMAGE = hello-world
export VERSION = 1.0.0

clean:
	rm -rf ./build

build: clean
	mkdir ./build
	docker build -t go -f Dockerfile.build .
	docker run --rm --volume $(PWD):/opt/build go build -o build/main main.go
	docker build -t $(USERNAME)/$(IMAGE):$(VERSION) .
	docker image tag $(USERNAME)/$(IMAGE):$(VERSION) $(USERNAME)/$(IMAGE):latest

publish: build
	docker login --username=$(USERNAME)
	docker push $(USERNAME)/$(IMAGE):$(VERSION)
	docker push $(USERNAME)/$(IMAGE):latest
