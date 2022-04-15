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
		Pkg  string

		rawDecls []RawDecl
		checked  struct{}
	}

	RawDecl   interface{ isRawDecl() }
	RawStruct struct {
		Ident *ast.Ident
		Type  *ast.StructType
	}
	RawVariant struct {
		Ident *ast.Ident
		Type  *ast.InterfaceType
		Cases []*RawCase
	}
	RawCase struct {
		Recv *ast.Field
		ID   int
	}
	RawInterface struct {
		Ident   *ast.Ident
		Methods []*ast.Field
	}
	RawEnum struct {
		Ident     *ast.Ident
		Constants []*ast.ValueSpec
	}
)

func (RawStruct) isRawDecl()    {}
func (RawVariant) isRawDecl()   {}
func (RawCase) isRawDecl()      {}
func (RawInterface) isRawDecl() {}
func (RawEnum) isRawDecl()      {}

func (g *Gbio) AddRaw(decl RawDecl) {
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
