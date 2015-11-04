package generator

import (
	"fmt"
	"log"
	"os"
	"strings"
	"text/template"

	"github.com/mitchellh/packer/common"
	"github.com/mitchellh/packer/helper/config"
	"github.com/mitchellh/packer/packer"
	"github.com/mitchellh/packer/provisioner/shell-local"
	"github.com/mitchellh/packer/template/interpolate"
)

type Config struct {
	common.PackerConfig `mapstructure:",squash"`

	Template string `mapstructure:"template"`
	Output   string `mapstructure:"output"`

	// ExecuteCommand is the command used to execute the command.
	ExecuteCommand []string `mapstructure:"execute_command"`

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

func getAmiMap(artifact packer.Artifact) map[string]string {

	meta := artifact.State("atlas.artifact.metadata")
	if meta == nil {
		log.Println("Artifact has no AWS info: artifact.State(atlas.artifact.metadata)")
		return nil
	}

	regionMap, ok := meta.(map[interface{}]interface{})
	if !ok {
		return nil
	}

	amiMap := make(map[string]string)
	for reg, ami := range regionMap {
		r, ok := reg.(string)
		if !ok {
			log.Printf("Couldnt convert Region to string: %#v \n", reg)
		}
		a, ok := ami.(string)
		if !ok {
			log.Printf("Couldnt convert Ami to string: %#v \n", ami)
		}
		r = strings.TrimPrefix(r, "region.")
		amiMap[r] = a
	}
	return amiMap
}

func (p *PostProcessor) PostProcess(ui packer.Ui, artifact packer.Artifact) (packer.Artifact, bool, error) {

	ui.Message(fmt.Sprintf("Generating: '%s' from: '%s'", p.config.Output, p.config.Template))

	//ui.Message(fmt.Sprintf("Artifact: id:%s string:%s files:%#v", artifact.Id(), artifact.String(), artifact.Files()))

	tmpl, err := template.ParseFiles(p.config.Template)
	if err != nil {
		return nil, true, fmt.Errorf("Failed to parse template: %s ", err)
	}

	info, err := os.Stat(p.config.Template)
	if err != nil {
		return nil, true, fmt.Errorf("Couldn't get Stat of file: %s", err)
	}

	out, err := os.Create(p.config.Output)
	if err != nil {
		return nil, true, fmt.Errorf("Failed to create file: %s", err)
	}

	out.Chmod(info.Mode())

	amiMap := getAmiMap(artifact)
	ui.Say(fmt.Sprintf("AWS amimap: %#v", amiMap))
	data := struct {
		Test     string
		Artifact packer.Artifact
		Config   Config
		Meta     interface{}
	}{
		Test:     fmt.Sprintf("ok"),
		Artifact: artifact,
		Config:   p.config,
		Meta:     amiMap,
	}

	ui.Message("Generating ...")
	err = tmpl.Execute(out, data)
	if err != nil {
		return nil, true, fmt.Errorf("Template execution failed: %s", err)
	}

	ui.Message(fmt.Sprintf("p.config.ExecuteCommand: %#v", p.config.ExecuteCommand))
	if len(p.config.ExecuteCommand) > 0 {
		// Make another communicator for local
		comm := &shell.Communicator{
			Ctx:            p.config.ctx,
			ExecuteCommand: p.config.ExecuteCommand,
		}

		// Build the remote command
		cmd := &packer.RemoteCmd{Command: p.config.Output}

		ui.Say(fmt.Sprintf(
			"Executing local command: %s",
			p.config.Output))
		if err := cmd.StartWithUi(comm, ui); err != nil {
			return nil, true, fmt.Errorf(
				"Error executing command: %s\n\n"+
					"Please see output above for more information.",
				p.config.Output)
		}
		if cmd.ExitStatus != 0 {
			return nil, true, fmt.Errorf(
				"Erroneous exit code %d while executing command: %s\n\n"+
					"Please see output above for more information.",
				cmd.ExitStatus,
				p.config.Output)
		}

	}

	a := &Artifact{
		Path: p.config.Output,
	}
	return a, true, nil
}
