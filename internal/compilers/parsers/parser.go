package parsers

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"

	"github.com/anqur/gbio/internal/compilers/codegens"
)

type Parser struct {
	codegens.Codegen

	File string
}

func (p *Parser) Parse() *Parser {
	p.FSet = token.NewFileSet()
	f, err := parser.ParseFile(p.FSet, p.File, nil, parser.AllErrors)
	if err != nil {
		panic(err)
	}
	if len(f.Imports) != 0 {
		panic(p.Errorf(f.Package, "sorry, imports not supported yet"))
	}
	return p.parseDecls(f)
}

func (p *Parser) parseDecls(f *ast.File) *Parser {
	fmt.Println(f)
	return p
}
