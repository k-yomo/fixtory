package fixtory

import (
	"bytes"
	"io"
	"testing"
)

func TestGenerate(t *testing.T) {
	type args struct {
		targetDir string
		types     []string
		pkgName   string
		newWriter func() (io.Writer, func(), error)
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "generate factory for Author, Article in example directory",
			args: args{
				targetDir: "example",
				types:     []string{"Article", "Author"},
				newWriter: func() (io.Writer, func(), error) {
					var b bytes.Buffer
					return &b, func() { }, nil
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Generate(tt.args.targetDir, tt.args.types, tt.args.pkgName, tt.args.newWriter); (err != nil) != tt.wantErr {
				t.Errorf("Generate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
