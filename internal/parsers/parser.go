package parsers

import (
	"github.com/anqur/gbio/internal/codegens"
)

type Parser struct {
	codegens.Codegen

	File string
}

func (p *Parser) Parse() *Parser {
	return p
}

func (p *Parser) Generate() {
}
