package fixtory

const fixtoryFileTpl = `
// Code generated by {{ .GeneratorName }}; DO NOT EDIT.

package {{ .PackageName }}

import (
{{- range .ImportPackages }}
	"{{ . }}"
{{- end}}
	"github.com/k-yomo/fixtory"
	"testing"
)

{{ .Body }}
`

const factoryTpl = `
{{$lowerStructName := .StructName | ToLower }}
{{$factoryInterface := printf "%s%s" .StructName "Factory" }}
{{$builderInterface := printf "%s%s" .StructName "Builder" }}
{{$factory := printf "%s%s" $lowerStructName "Factory" }}
{{$builder := printf "%s%s" $lowerStructName "Builder" }}
{{$fieldType := printf "%s%s" .StructName "Field" }}

type {{ $factoryInterface }} interface {
	NewBuilder(bluePrint {{ .StructName }}BluePrintFunc, traits ...{{ .StructName }}) {{ $builderInterface }}
	OnBuild(onBuild func(t *testing.T, {{ $lowerStructName }} *{{ .StructName }}))
	Reset()
}

type {{ $builderInterface }} interface {
	EachParam({{ $lowerStructName }}Params ...{{ .StructName }}) {{ $builderInterface }}
	Zero({{ $lowerStructName }}Fields ...{{ $fieldType }}) {{ $builderInterface }}
	ResetAfter() {{ $builderInterface }}

	Build() *{{ .StructName }}
	Build2() (*{{ .StructName }}, *{{ .StructName }})
	Build3() (*{{ .StructName }}, *{{ .StructName }}, *{{ .StructName }})
	BuildList(n int) []*{{ .StructName }}
}

type {{ .StructName }}BluePrintFunc func(i int, last {{ .StructName }}) {{ .StructName }}

type {{ $fieldType }} string

const (
{{- range .FieldNames }}
	{{ $.StructName }}{{ . }}Field {{ $fieldType }} = "{{ . }}"
{{- end}}
)

type {{ $factory }} struct {
	t       *testing.T
	factory *fixtory.Factory
}

type {{ $builder }} struct {
	t       *testing.T
	builder *fixtory.Builder
}

func New{{ .StructName }}Factory(t *testing.T) {{ $factoryInterface }} {
	t.Helper()

	return &{{ $factory }}{t: t, factory: fixtory.NewFactory(t, {{ .StructName }}{})}
}

func (uf *{{ $factory }}) NewBuilder(bluePrint {{ .StructName }}BluePrintFunc, {{ $lowerStructName }}Traits ...{{ .StructName }}) {{ $builderInterface }} {
	uf.t.Helper()

	var bp fixtory.BluePrintFunc
	if bluePrint != nil {
		bp = func(i int, last interface{}) interface{} { return bluePrint(i, last.({{ .StructName }})) }
	}
	builder := uf.factory.NewBuilder(bp, fixtory.ConvertToInterfaceArray({{ $lowerStructName }}Traits)...)

	return &{{ $builder }}{t: uf.t, builder: builder}
}

func (uf *{{ $factory }}) OnBuild(onBuild func(t *testing.T, {{ $lowerStructName }} *{{ .StructName }})) {
	uf.t.Helper()

	uf.factory.OnBuild = func(t *testing.T, v interface{}) { onBuild(t, v.(*{{ .StructName }})) }
}

func (uf *{{ $factory }}) Reset() {
	uf.t.Helper()

	uf.factory.Reset()
}

func (ub *{{ $builder }}) Zero({{ $lowerStructName }}Fields ...{{ $fieldType }}) {{ $builderInterface }} {
	ub.t.Helper()

	fields := make([]string, 0, len({{ $lowerStructName }}Fields))
	for _, f := range {{ $lowerStructName }}Fields {
		fields = append(fields, string(f))
	}

	ub.builder = ub.builder.Zero(fields...)
	return ub
}
func (ub *{{ $builder }}) ResetAfter() {{ $builderInterface }} {
	ub.t.Helper()

	ub.builder = ub.builder.ResetAfter()
	return ub
}

func (ub *{{ $builder }}) EachParam({{ $lowerStructName }}Params ...{{ .StructName }}) {{ $builderInterface }} {
	ub.t.Helper()

	ub.builder = ub.builder.EachParam(fixtory.ConvertToInterfaceArray({{ $lowerStructName }}Params)...)
	return ub
}

func (ub *{{ $builder }}) Build() *{{ .StructName }} {
	ub.t.Helper()

	return ub.builder.Build().(*{{ .StructName }})
}

func (ub *{{ $builder }}) Build2() (*{{ .StructName }}, *{{ .StructName }}) {
	ub.t.Helper()

	return ub.Build(), ub.Build()
}

func (ub *{{ $builder }}) Build3() (*{{ .StructName }}, *{{ .StructName }}, *{{ .StructName }}) {
	ub.t.Helper()

	return ub.Build(), ub.Build(), ub.Build()
}

func (ub *{{ $builder }}) BuildList(n int) []*{{ .StructName }} {
	ub.t.Helper()

	{{ $lowerStructName }}s := make([]*{{ .StructName }}, 0, n)
	for i := 0; i < n; i++ {
		{{ $lowerStructName }}s = append({{ $lowerStructName }}s, ub.builder.Build().(*{{ .StructName }}))
	}
	return {{ $lowerStructName }}s
}
`
