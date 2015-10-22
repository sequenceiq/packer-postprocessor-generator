PLUGIN_NAME=post-processor-generator
BINARY_NAME=packer-$(PLUGIN_NAME)
GH_PROJECT=sequenceiq/packer-postprocessor-generator
VERSION=0.8.7

deps:
	go get ./...

build:
	rm -rf build
	GOOS=linux go build -o build/Linux/$(BINARY_NAME)  plugin/$(PLUGIN_NAME)/main.go
	GOOS=darwin go build -o build/Darwin/$(BINARY_NAME)  plugin/$(PLUGIN_NAME)/main.go

gh-release-prepare: build
	rm -rf release; mkdir -p release

	#cp build/Darwin/$(BINARY_NAME) release/$(BINARY_NAME)-Darwin
	#cp build/Linux/$(BINARY_NAME) release/$(BINARY_NAME)-Linux

	tar czvf release/$(BINARY_NAME)-Darwin.tgz -C build/Darwin/ $(BINARY_NAME)
	tar czvf release/$(BINARY_NAME)-Linux.tgz -C build/Linux/ $(BINARY_NAME)

gh-release: build
	gh-release create $(GH_PROJECT) $(VERSION)

dev-install:
	go build -v -o ~/.packer.d/plugins/$(BINARY_NAME) ./plugin/$(PLUGIN_NAME)/

integration-test: build
	packer build packer-test.json

.PHONY: build
