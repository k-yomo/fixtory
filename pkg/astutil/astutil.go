package astutil

import (
	"errors"
	"go/ast"
	"go/build"
	"go/parser"
	"go/token"
	"path/filepath"
	"sort"
	"strings"
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
		return nil, err
	}

	m := make(map[string]AstPkgWalker, len(pkgMap))
	for k, v := range pkgMap {
		m[k], err = ParseAstPkg(fileSet, v)
		if err != nil {
			return nil, err
		}
	}
	return m, nil
}

// ParseAstPkg parses package ast
func ParseAstPkg(fset *token.FileSet, pkg *ast.Package) (AstPkgWalker, error) {
	var aFilePath string
	for _, file := range pkg.Files {
		aFilePath = fset.File(file.Package).Name()
	}

	pkgPath, err := LocalPathToPackagePath(filepath.Dir(aFilePath))
	if err != nil {
		return AstPkgWalker{}, err
	}

	return AstPkgWalker{
		FileSet: fset,
		Pkg:     pkg,
		Files:   ToSortedFileListFromFileMapOfAst(pkg.Files),
		Decls:   AllDeclsFromAstPkg(pkg),

		PkgPath: pkgPath,
	}, nil
}

func LocalPathToPackagePath(s string) (string, error) {
	s, err := filepath.Abs(s)
	if err != nil {
		return "", err
	}

	s = filepath.ToSlash(s)

	for _, srcDir := range build.Default.SrcDirs() {
		srcDir = filepath.ToSlash(srcDir)
		prefix := srcDir + "/"
		if strings.HasPrefix(s, prefix) {
			return strings.TrimPrefix(s, prefix), nil
		}
	}
	return "", errors.New("failed to resolve package path")
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
