package astutil

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"golang.org/x/mod/modfile"
	"golang.org/x/xerrors"
	"io/ioutil"
	"os"
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
		return nil, xerrors.Errorf("parser.ParseDir: %w", err)
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

	dir := filepath.Dir(aFilePath)
	if dir == "." {
		var err error
		dir, err = os.Getwd()
		if err != nil {
			return AstPkgWalker{}, xerrors.Errorf("get current directory: %w", err)
		}
	}
	pkgPath, err := packageNameOfDir(dir)
	if err != nil {
		return AstPkgWalker{}, err
	}
	return AstPkgWalker{
		Pkg:     pkg,
		Decls:   allDeclsFromAstPkg(pkg),
		PkgPath: pkgPath,
	}, nil
}

func allDeclsFromAstPkg(pkg *ast.Package) []ast.Decl {
	decls := make([]ast.Decl, 0)
	for _, file := range getSortedFileListFromFileAstMap(pkg.Files) {
		decls = append(decls, file.Decls...)
	}
	return decls
}

func getSortedFileListFromFileAstMap(s map[string]*ast.File) []*ast.File {
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

// packageNameOfDir get package import path via dir
func packageNameOfDir(srcDir string) (string, error) {
	files, err := ioutil.ReadDir(srcDir)
	if err != nil {
		return "", err
	}

	var goFilePath string
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".go") {
			goFilePath = file.Name()
			break
		}
	}
	if goFilePath == "" {
		return "", fmt.Errorf("go source file not found %s", srcDir)
	}

	packageImport, err := parsePackageImport(srcDir)
	if err != nil {
		return "", err
	}
	return packageImport, nil
}

// parseImportPackage get package import path via source file
// an alternative implementation is to use:
// cfg := &packages.Config{Mode: packages.NeedName, Tests: true, Dir: srcDir}
// pkgs, err := packages.Load(cfg, "file="+source)
// However, it will call "go list" and slow down the performance
func parsePackageImport(srcDir string) (string, error) {
	moduleMode := os.Getenv("GO111MODULE")
	// trying to find the module
	if moduleMode != "off" {
		currentDir := srcDir
		for {
			dat, err := ioutil.ReadFile(filepath.Join(currentDir, "go.mod"))
			if os.IsNotExist(err) {
				if currentDir == filepath.Dir(currentDir) {
					// at the root
					break
				}
				currentDir = filepath.Dir(currentDir)
				continue
			} else if err != nil {
				return "", err
			}
			modulePath := modfile.ModulePath(dat)
			return filepath.ToSlash(filepath.Join(modulePath, strings.TrimPrefix(srcDir, currentDir))), nil
		}
	}
	// fall back to GOPATH mode
	goPaths := os.Getenv("GOPATH")
	if goPaths == "" {
		return "", xerrors.New("GOPATH is not set")
	}
	goPathList := strings.Split(goPaths, string(os.PathListSeparator))
	for _, goPath := range goPathList {
		sourceRoot := filepath.Join(goPath, "src") + string(os.PathSeparator)
		if strings.HasPrefix(srcDir, sourceRoot) {
			return filepath.ToSlash(strings.TrimPrefix(srcDir, sourceRoot)), nil
		}
	}
	return "", xerrors.New("Source directory is outside GOPATH")
}
