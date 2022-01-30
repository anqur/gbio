package codegens

import "github.com/anqur/gbio/internal/compilers/langs"

type Codegen struct {
	OutDir        string
	Target        string
	Marshaller    string
	Discriminator string

	langs.Gbio
}

func (c *Codegen) Generate() {
}
