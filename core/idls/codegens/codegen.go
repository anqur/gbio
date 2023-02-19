package codegens

import "github.com/anqur/gbio/core/idls/specs"

type Codegen struct {
	OutDir        string
	Target        string
	Marshaller    string
	Discriminator string

	specs.Gbio
}

func (c *Codegen) Generate() {
}
