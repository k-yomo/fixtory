package astutil

import (
	"go/ast"
	"go/parser"
	"go/token"
	"golang.org/x/xerrors"
	"path/filepath"
	"sort"
)

// AstPkgWalker represents ast package walker
type AstPkgWalker struct {
	FileSet *token.FileSet
	Pkg     *ast.Package
	Files   []*ast.File
	Decls   []ast.Decl

	PkgPath string
}

// AllGenDecls returns all generic declaration nodes
func (w AstPkgWalker) AllGenDecls() []*ast.GenDecl {
	decls := w.Decls
	l := make([]*ast.GenDecl, 0, len(decls))
	for _, decl := range decls {
		decl, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}
		l = append(l, decl)
	}
	return l
}

// AllStructSpecs returns all struct type specs
func (w AstPkgWalker) AllStructSpecs() []*ast.TypeSpec {
	decls := w.AllGenDecls()
	l := make([]*ast.TypeSpec, 0)
	for _, decl := range decls {
		if decl.Tok != token.TYPE {
			continue
		}
		for _, spec := range decl.Specs {
			typeSpec := spec.(*ast.TypeSpec)
			_, ok := typeSpec.Type.(*ast.StructType)
			if !ok {
				continue
			}
			l = append(l, typeSpec)
		}
	}
	return l
}

// DirToAstWalker initializes map of AstPkgWalker from given directory
func DirToAstWalker(targetDir string) (map[string]AstPkgWalker, error) {
	fileSet := token.NewFileSet()
	pkgMap, err := parser.ParseDir(
		fileSet,
		filepath.FromSlash(targetDir),
		nil,
		parser.ParseComments,
	)
	if err != nil {
		return nil, xerrors.Errorf("parser.ParseDir: %w", err)
	}

	m := make(map[string]AstPkgWalker, len(pkgMap))
	for k, v := range pkgMap {
		m[k] = ParseAstPkg(v)
	}
	return m, nil
}

// ParseAstPkg parses package ast
func ParseAstPkg(pkg *ast.Package) AstPkgWalker {
	return AstPkgWalker{
		Pkg:     pkg,
		Decls:   AllDeclsFromAstPkg(pkg),
	}
}

func AllDeclsFromAstPkg(pkg *ast.Package) []ast.Decl {
	decls := make([]ast.Decl, 0)
	for _, file := range ToSortedFileListFromFileMapOfAst(pkg.Files) {
		decls = append(decls, file.Decls...)
	}
	return decls
}

func ToSortedFileListFromFileMapOfAst(s map[string]*ast.File) []*ast.File {
	names := make([]string, 0, len(s))
	for name := range s {
		names = append(names, name)
	}
	sort.Strings(names)

	l := make([]*ast.File, len(names))
	for i, name := range names {
		l[i] = s[name]
	}
	return l
}
