package langs

import (
	"fmt"
	"go/ast"
	"go/token"
)

type (
	Gbio struct {
		FSet *token.FileSet

		Decls []Decl
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
		Case  []*CaseType
	}
	CaseType struct {
		Ident *ast.Ident
		Type  *ast.StructType
		ID    int
	}

	InterfaceType struct {
		Ident   *ast.Ident
		Methods []*ast.Field
	}

	EnumType struct {
		Ident     *ast.Ident
		Constants []*ast.ValueSpec
	}
)

func (StructType) isDecl()    {}
func (VariantType) isDecl()   {}
func (InterfaceType) isDecl() {}
func (EnumType) isDecl()      {}

func (l *Gbio) Errorf(pos token.Pos, format string, args ...interface{}) error {
	msg := fmt.Sprintf(format, args...)
	return fmt.Errorf("%s: %s", l.FSet.Position(pos), msg)
}
