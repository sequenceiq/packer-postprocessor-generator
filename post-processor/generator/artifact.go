package generator

import (
	"fmt"
	"os"
)

const BuilderId = "packer.post-processor.generator"

type Artifact struct {
	Path  string
	files []string
}

func (a *Artifact) BuilderId() string {
	return fmt.Sprintf("builder.genrator")
}

func (a *Artifact) Id() string {
	return a.Path
}

func (a *Artifact) Files() []string {
	return []string{a.Path}
}

func (a *Artifact) String() string {
	return fmt.Sprintf("Generatied: %s", a.Path)
}

func (*Artifact) State(name string) interface{} {
	return nil
}

func (a *Artifact) Destroy() error {
	return os.Remove(a.Path)
}
