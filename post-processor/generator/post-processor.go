package generator

import (
	"fmt"
	"log"

	"github.com/mitchellh/packer/common"
	"github.com/mitchellh/packer/helper/config"
	"github.com/mitchellh/packer/packer"
	"github.com/mitchellh/packer/template/interpolate"
)

type Config struct {
	common.PackerConfig `mapstructure:",squash"`

	Template string `mapstructure:"template"`
	Output   string `mapstructure:"output"`

	ctx interpolate.Context
}

type PostProcessor struct {
	config Config
}

func (p *PostProcessor) Configure(raws ...interface{}) error {
	err := config.Decode(&p.config, &config.DecodeOpts{
		Interpolate:        true,
		InterpolateContext: &p.config.ctx,
		InterpolateFilter: &interpolate.RenderFilter{
			Exclude: []string{},
		},
	}, raws...)
	if err != nil {
		return err
	}

	// Accumulate any errors
	errs := new(packer.MultiError)

	if p.config.Template == "" {
		errs = packer.MultiErrorAppend(errs, fmt.Errorf("Template must be set"))
	}
	if p.config.Output == "" {
		errs = packer.MultiErrorAppend(errs, fmt.Errorf("Output must be set"))
	}

	log.Printf("Configure(): template:%s output:%s", p.config.Template, p.config.Output)
	if len(errs.Errors) > 0 {
		return errs
	}

	return nil
}

func (p *PostProcessor) PostProcess(ui packer.Ui, artifact packer.Artifact) (packer.Artifact, bool, error) {

	ui.Message(fmt.Sprintf("Generating: '%s' from: '%s'", p.config.Output, p.config.Template))

	ui.Message(fmt.Sprintf("Artifact: id:%s string:%s files:%#v", artifact.Id(), artifact.String(), artifact.Files()))
	return artifact, true, nil
	a := &Artifact{
		Path: p.config.Output,
	}
	return a, true, nil
}
