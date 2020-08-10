package fixtory

import (
	"bytes"
	"github.com/k-yomo/fixtory/pkg/astutil"
	"go/ast"
	"go/format"
	"golang.org/x/xerrors"
	"io"
	"strings"
	"text/template"
)

type tmplParam struct {
	StructName string
	FieldNames []string
}

func Generate(targetDir string, types []string, pkgName *string, newWriter func() (writer io.Writer, close func() error, err error)) error {
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
		body := new(bytes.Buffer)
		for _, spec := range walker.AllStructSpecs() {
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
			params := tmplParam{StructName: spec.Name.Name, FieldNames: fieldNames}
			if err := tpl.Execute(body, params); err != nil {
				return xerrors.Errorf("execute factory template: %w", err)
			}
		}
		if body.Len() == 0 {
			continue
		}

		if pkgName == nil {
			pkgName = &walker.Pkg.Name
		}

		out := new(bytes.Buffer)
		err = template.Must(template.New("fixtoryFile").Parse(fixtoryFileTpl)).Execute(out, map[string]string{
			"GeneratorName": "fixtory",
			"PackageName":   *pkgName,
			"Body":          body.String(),
		})
		if err != nil {
			return xerrors.Errorf("execute fixtoryFile template: %w", err)
		}

		str, err := format.Source(out.Bytes())
		if err != nil {
			return xerrors.Errorf("format.Source: %w", err)
		}

		writer, closeWriter, err := newWriter()
		defer closeWriter()

		if _, err := writer.Write(str); err != nil {
			return xerrors.Errorf("write fixtory file: %w", err)
		}
	}

	return nil
}
