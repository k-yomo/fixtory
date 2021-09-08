package fixtory

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"io"
	"strings"
	"text/template"

	"github.com/k-yomo/fixtory/pkg/astutil"
	"golang.org/x/xerrors"
)

func Generate(targetDir string, types []string, pkgName string, newWriter func() (writer io.Writer, close func(), err error)) error {
	targetTypeMap := map[string]bool{}
	for _, t := range types {
		targetTypeMap[t] = true
	}

	walkerMap, err := astutil.DirToAstWalker(targetDir)
	if err != nil {
		return err
	}
	for _, walker := range walkerMap {
		if len(targetTypeMap) == 0 {
			break
		}
		walkerPkgName := walker.Pkg.Name
		allStructSpecs := walker.AllStructSpecs()
		for _, spec := range allStructSpecs {
			if strings.ToLower(spec.Name.Name) == walkerPkgName {
				walkerPkgName = "_" + walkerPkgName
				break
			}
		}
		body := new(bytes.Buffer)
		for _, spec := range allStructSpecs {
			if len(targetTypeMap) == 0 {
				break
			}
			if !targetTypeMap[spec.Name.Name] {
				continue
			}
			delete(targetTypeMap, spec.Name.Name)

			structType := spec.Type.(*ast.StructType)
			fieldNames := make([]string, 0, len(structType.Fields.List))
			for _, field := range structType.Fields.List {
				fieldNames = append(fieldNames, field.Names[0].Name)
			}
			tpl := template.Must(template.New("factory").Funcs(template.FuncMap{"ToLower": strings.ToLower}).Parse(factoryTpl))
			st := spec.Name.Name
			if pkgName != "" {
				st = fmt.Sprintf("%s.%s", walkerPkgName, st)
			}
			params := struct {
				StructName string
				Struct     string
				FieldNames []string
			}{
				StructName: spec.Name.Name,
				Struct:     st,
				FieldNames: fieldNames,
			}
			if err := tpl.Execute(body, params); err != nil {
				return xerrors.Errorf("execute factory template: %w", err)
			}
		}
		if body.Len() == 0 {
			continue
		}

		var importPackages []string
		if pkgName == "" {
			pkgName = walkerPkgName
		} else {
			importPackage := fmt.Sprintf(`"%s"`, walker.PkgPath)
			if walkerPkgName != walker.Pkg.Name {
				importPackage = fmt.Sprintf(`%s "%s"`, walkerPkgName, walker.PkgPath)
			}

			importPackages = append(importPackages, importPackage)
		}

		out := new(bytes.Buffer)
		params := struct {
			GeneratorName  string
			PackageName    string
			ImportPackages []string
			Body           string
		}{
			GeneratorName:  "fixtory",
			PackageName:    pkgName,
			ImportPackages: importPackages,
			Body:           body.String(),
		}
		err = template.Must(template.New("fixtoryFile").Parse(fixtoryFileTpl)).Execute(out, params)
		if err != nil {
			return xerrors.Errorf("execute fixtoryFile template: %w", err)
		}

		str, err := format.Source(out.Bytes())
		if err != nil {
			return xerrors.Errorf("format.Source: %w", err)
		}

		writer, closeWriter, err := newWriter()
		if err != nil {
			return xerrors.Errorf("initialize writer: %w", err)
		}
		defer closeWriter()

		if _, err := writer.Write(str); err != nil {
			return xerrors.Errorf("write fixtory file: %w", err)
		}
	}

	return nil
}
