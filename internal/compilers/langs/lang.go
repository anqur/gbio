package langs

import (
	"fmt"
	"go/ast"
	"go/token"
	"reflect"
	"strings"
)

type (
	Gbio struct {
		FSet *token.FileSet

		rawDecls []Decl
	}

	Decl interface {
		isDecl()
	}
	StructType struct {
		Ident *ast.Ident
		Type  *ast.StructType
	}
	VariantType struct {
		Ident *ast.Ident
		Type  *ast.InterfaceType
		Cases []*Case
	}
	Case struct {
		Recv *ast.Field
		ID   int
	}
	InterfaceType struct {
		Ident   *ast.Ident
		Methods []*ast.Field
	}
	EnumType struct {
		Ident     *ast.Ident
		Constants []*ast.ValueSpec
	}
	Constant struct {
		// TODO: Type
		Value *ast.ValueSpec
	}
)

func (StructType) isDecl()    {}
func (VariantType) isDecl()   {}
func (Case) isDecl()          {}
func (InterfaceType) isDecl() {}
func (EnumType) isDecl()      {}
func (Constant) isDecl()      {}

func (g *Gbio) AddRawDecl(decl Decl) {
	g.rawDecls = append(g.rawDecls, decl)
}

func (g *Gbio) Errorf(pos token.Pos, format string, args ...interface{}) error {
	msg := fmt.Sprintf(format, args...)
	return fmt.Errorf("%s: %s", g.FSet.Position(pos), msg)
}

func (g *Gbio) IsContextField(f *ast.Field) bool {
	return f.Tag != nil &&
		reflect.
			StructTag(strings.Trim(f.Tag.Value, "`")).
			Get("json") == "-"
}
