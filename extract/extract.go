package extract

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strconv"
	"strings"
)

type NotFoundError struct {
	Declaration string
}

func (e NotFoundError) Error() string {
	return fmt.Sprintf("declaration %q not found", e.Declaration)
}

func (m NotFoundError) Is(target error) bool {
	_, ok := target.(*NotFoundError)
	_, ok2 := target.(NotFoundError)
	return ok || ok2
}

func ExtractDeclarationFromFile(path, declName string) (any, error) {
	f, err := parser.ParseFile(token.NewFileSet(), path, nil, parser.Mode(0))
	if err != nil {
		return nil, fmt.Errorf("failed to parse file %q: %w", path, err)
	}

	for _, decl := range f.Decls {
		if gen, ok := decl.(*ast.GenDecl); ok && (gen.Tok == token.CONST || gen.Tok == token.VAR) {
			for _, spec := range gen.Specs {
				spec := spec.(*ast.ValueSpec)
				for i, name := range spec.Names {
					if name.Name == declName {
						v, ok := spec.Values[i].(*ast.BasicLit)
						if !ok {
							return nil, fmt.Errorf("declaration %q is not a basic literal found %T", declName, spec.Values[i])
						}
						return basicLitValue(v)
					}
				}
			}
		}
	}

	return nil, &NotFoundError{Declaration: declName}
}

func basicLitValue(b *ast.BasicLit) (any, error) {
	switch b.Kind {
	case token.INT:
		return strconv.ParseInt(b.Value, 0, 64)
	case token.FLOAT:
		return strconv.ParseFloat(b.Value, 64)
	case token.IMAG:
		return strconv.ParseComplex(b.Value, 128)
	case token.CHAR:
		r, _, _, err := strconv.UnquoteChar(strings.Trim(b.Value, "'"), 0)
		return r, err
	case token.STRING:
		return strconv.Unquote(b.Value)
	}
	// Should never happen
	return nil, fmt.Errorf("unknown token kind %q", b.Kind.String())
}
