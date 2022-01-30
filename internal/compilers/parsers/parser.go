package parsers

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"reflect"
	"unicode"

	"github.com/anqur/gbio/internal/compilers/codegens"
	"github.com/anqur/gbio/internal/compilers/langs"
	"github.com/anqur/gbio/internal/utils"
)

var (
	validFieldKinds = []reflect.Kind{
		reflect.Bool,
		reflect.String,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Float32, reflect.Float64,
	}
	validFieldTypeIdents []string
)

func init() {
	for _, k := range validFieldKinds {
		validFieldTypeIdents = append(validFieldTypeIdents, k.String())
	}
}

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
	return p.parseDecls(f.Decls)
}

func (p *Parser) parseDecls(decls []ast.Decl) *Parser {
	for _, decl := range decls {
		if g, ok := decl.(*ast.GenDecl); ok {
			p.parseGenDecl(g)
			continue
		}
		if f, ok := decl.(*ast.FuncDecl); ok {
			p.parseFuncDecl(f)
			continue
		}
		panic(p.Errorf(
			decl.Pos(),
			"unsupported declaration %v",
			decl,
		))
	}
	return p
}

func (*Parser) parseFuncDecl(fn *ast.FuncDecl) {
	// TODO
}

func (p *Parser) parseGenDecl(g *ast.GenDecl) {
	switch g.Tok {
	case token.TYPE:
		for _, spec := range g.Specs {
			p.parseTypeSpec(spec.(*ast.TypeSpec))
		}
	case token.CONST:
		for _, spec := range g.Specs {
			p.parseConstDef(spec.(*ast.ValueSpec))
		}
	default:
		panic(p.Errorf(
			g.Pos(),
			"unsupported general declaration %v, expected `type` or `const`",
			g.Tok.String(),
		))
	}
}

func (p *Parser) parseTypeSpec(s *ast.TypeSpec) {
	if st, ok := s.Type.(*ast.StructType); ok {
		p.parseStructType(s.Name, st)
		return
	}
	if i, ok := s.Type.(*ast.InterfaceType); ok {
		p.parseInterfaceType(s.Name, i)
		return
	}
	if id, ok := s.Type.(*ast.Ident); ok && s.Assign != 0 {
		p.parseTypeAlias(s.Name, id)
		return
	}
	panic(p.Errorf(
		s.Pos(),
		"unsupported type spec %v, expected `struct`, `interface` or type aliases",
		s.Name,
	))
}

func (p *Parser) parseConstDef(s *ast.ValueSpec) {
	// TODO
}

func (p *Parser) checkStructFieldType(f *ast.Field, t ast.Expr) {
	if len(f.Names) == 0 {
		panic(p.Errorf(t.Pos(), "embedded struct field not supported"))
	}

	name := f.Names[0]

	if id, ok := t.(*ast.Ident); ok {
		if obj := id.Obj; obj != nil {
			if spec, ok := obj.Decl.(*ast.TypeSpec); ok {
				p.checkStructFieldType(f, spec.Type)
				return
			}
		}
		if utils.OneOf(id.Name, validFieldTypeIdents) {
			return
		}
		panic(p.Errorf(
			t.Pos(),
			"unexpected type %v for field %v, expected primitives %v",
			id,
			name,
			validFieldTypeIdents,
		))
	}
	if arr, ok := t.(*ast.ArrayType); ok {
		if !p.IsContextField(f) {
			p.checkStructFieldType(f, arr.Elt)
			return
		}
		panic(p.Errorf(
			t.Pos(),
			"unexpected type for context field %v, expected primitives %v",
			name,
			validFieldTypeIdents,
		))
	}
	panic(p.Errorf(
		t.Pos(),
		"unexpected type for field %v, expected primitives or slices",
		name,
	))
}

func (p *Parser) parseStructType(name *ast.Ident, st *ast.StructType) {
	for _, field := range st.Fields.List {
		p.checkStructFieldType(field, field.Type)
	}
	p.AddRawDecl(&langs.StructType{Ident: name, Type: st})
}

func (*Parser) isVariantType(name *ast.Ident, i *ast.InterfaceType) bool {
	if i.Methods.NumFields() != 1 {
		return false
	}
	m := i.Methods.List[0]
	if fmt.Sprintf("is%s", name) != m.Names[0].Name {
		return false
	}
	fn, ok := m.Type.(*ast.FuncType)
	if !ok {
		return false
	}
	if fn.Params.NumFields() != 0 {
		return false
	}
	if fn.Results.NumFields() != 1 {
		return false
	}
	retType, ok := fn.Results.List[0].Type.(*ast.Ident)
	if !ok {
		return false
	}
	return retType.Name == reflect.Int.String()
}

func (p *Parser) checkInterfaceMethodType(f *ast.Field, t ast.Expr) {
	if len(f.Names) == 0 {
		panic(p.Errorf(t.Pos(), "embedded interface field not supported"))
	}

	name := f.Names[0]
	if !unicode.IsUpper(rune(name.Name[0])) {
		panic(p.Errorf(t.Pos(), "unexpected private method %v", name))
	}

	fn, ok := t.(*ast.FuncType)
	if !ok {
		panic(p.Errorf(
			t.Pos(),
			"unexpected field type of %v, expected a method",
			name,
		))
	}

	p.checkInterfaceMethodParamType(name, t, fn.Params)
	p.checkInterfaceMethodReturnType(name, t, fn.Results)
}

func (p *Parser) checkInterfaceMethodParamType(
	name *ast.Ident,
	t ast.Expr,
	fl *ast.FieldList,
) {
	if n := fl.NumFields(); n != 1 {
		panic(p.Errorf(
			t.Pos(),
			"unexpected parameter length %d of %v, expected 1",
			n,
			name,
		))
	}
	p.checkInterfaceMethodTypeReference(fl.List[0].Type)
}

func (p *Parser) checkInterfaceMethodReturnType(
	name *ast.Ident,
	t ast.Expr,
	fl *ast.FieldList,
) {
	if n := fl.NumFields(); n != 1 {
		panic(p.Errorf(
			t.Pos(),
			"unexpected return value length %d of %v, expected 1",
			n,
			name,
		))
	}
	p.checkInterfaceMethodTypeReference(fl.List[0].Type)
}

func (p *Parser) checkInterfaceMethodTypeReference(t ast.Expr) {
	// TODO
}

func (p *Parser) parseInterfaceType(name *ast.Ident, i *ast.InterfaceType) {
	if i.Methods.NumFields() == 0 {
		panic(p.Errorf(i.Pos(), "unexpected empty interface"))
	}
	if p.isVariantType(name, i) {
		p.AddRawDecl(&langs.VariantType{Ident: name, Type: i})
		return
	}
	for _, field := range i.Methods.List {
		p.checkInterfaceMethodType(field, field.Type)
	}
}

func (p *Parser) parseTypeAlias(name, typ *ast.Ident) {
	// TODO
}
