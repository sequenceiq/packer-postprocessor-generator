package main

import (
	"github.com/mitchellh/packer/packer/plugin"
	"github.com/sequenceiq/packer-postprocessor-template/post-processor/generator"
)

func main() {
	server, err := plugin.Server()
	if err != nil {
		panic(err)
	}
	server.RegisterPostProcessor(new(generator.PostProcessor))
	server.Serve()
}
