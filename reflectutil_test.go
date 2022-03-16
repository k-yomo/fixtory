package fixtory

import (
	"reflect"
	"testing"
)

func TestMapNotZeroFields(t *testing.T) {
	type args struct {
		from testStruct
		to   *testStruct
	}
	tests := []struct {
		name    string
		args    args
		want    *testStruct
		wantErr bool
	}{
		{
			name: "fields of from struct overwrite fields of to struct",
			args: args{
				from: testStruct{
					String:    "after",
					StringPtr: func() *string { s := "after"; return &s }(),
					Int:       10,
					Float:     10.5,
					Array:     []int{9, 8, 7},
					Map:       map[string]bool{"after": true},
					ChildStruct: &childStruct{
						String: "after",
						Int:    10,
					},
				},
				to: &testStruct{
					String:    "a",
					StringPtr: func() *string { s := "a"; return &s }(),
					Int:       1,
					Float:     0.5,
					Array:     []int{1, 2, 3},
					Map:       map[string]bool{"a": true},
					ChildStruct: &childStruct{
						String: "b",
						Int:    2,
					},
				},
			},
			want: &testStruct{
				String:    "after",
				StringPtr: func() *string { s := "after"; return &s }(),
				Int:       10,
				Float:     10.5,
				Array:     []int{9, 8, 7},
				Map:       map[string]bool{"after": true},
				ChildStruct: &childStruct{
					String: "after",
					Int:    10,
				},
			},
		},
		{
			name: "zero value fields of from struct do not overwrite",
			args: args{
				from: testStruct{},
				to: &testStruct{
					String:    "a",
					StringPtr: func() *string { s := "a"; return &s }(),
					Int:       1,
					Float:     0.5,
					Array:     []int{1, 2, 3},
					Map:       map[string]bool{"a": true},
					ChildStruct: &childStruct{
						String: "b",
						Int:    2,
					},
				},
			},
			want: &testStruct{
				String:    "a",
				StringPtr: func() *string { s := "a"; return &s }(),
				Int:       1,
				Float:     0.5,
				Array:     []int{1, 2, 3},
				Map:       map[string]bool{"a": true},
				ChildStruct: &childStruct{
					String: "b",
					Int:    2,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if err := recover(); (err != nil) != tt.wantErr {
					t.Errorf("MapNotZeroFields() = %v, wantErr %v", err, tt.wantErr)
				}
			}()
			MapNotZeroFields[testStruct](tt.args.from, tt.args.to)
			if !reflect.DeepEqual(tt.args.to, tt.want) {
				t.Errorf("MapNotZeroFields() = %v, want %v", tt.args.to, tt.want)
			}
		})
	}
}
