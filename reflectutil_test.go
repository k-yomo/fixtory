package fixtory

import (
	"reflect"
	"testing"
)

func TestMapNotZeroFields(t *testing.T) {
	type childStruct struct {
		String string
		Int    int
	}

	type testStruct struct {
		String      string
		StringPtr   *string
		Int         int
		Float       float64
		Array       []int
		Map         map[string]bool
		ChildStruct *childStruct
	}

	type args struct {
		from interface{}
		to   interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
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
		{
			name: "panic when from is not struct",
			args: args{
				from: "a",
				to:   &testStruct{},
			},
			wantErr: true,
		},
		{
			name: "panic when to is not struct pointer",
			args: args{
				from: testStruct{},
				to:   "a",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if err := recover(); (err != nil) != tt.wantErr {
					t.Errorf("MapNotZeroFields() = %v, wantErr %v", err, tt.wantErr)
				}
			}()
			MapNotZeroFields(tt.args.from, tt.args.to)
			if !reflect.DeepEqual(tt.args.to, tt.want) {
				t.Errorf("MapNotZeroFields() = %v, want %v", tt.args.to, tt.want)
			}
		})
	}
}

func TestConvertToInterfaceArray(t *testing.T) {
	type testStruct struct {
		Val string
	}

	type args struct {
		from interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    []interface{}
		wantErr bool
	}{
		{
			name: "convert type of string array to interface array",
			args: args{
				from: []string{"a", "b", "c"},
			},
			want: []interface{}{"a", "b", "c"},
		},
		{
			name: "convert type of int array to interface array",
			args: args{
				from: []int{1, 2, 3},
			},
			want: []interface{}{1, 2, 3},
		},
		{
			name: "convert type of struct array to interface array",
			args: args{
				from: []testStruct{{Val: "a"}, {Val: "b"}, {Val: "c"}},
			},
			want: []interface{}{testStruct{Val: "a"}, testStruct{Val: "b"}, testStruct{Val: "c"}},
		},
		{
			name:    "panic when from is not array",
			args:    args{from: "a"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if err := recover(); (err != nil) != tt.wantErr {
					t.Errorf("ConvertToInterfaceArray() = %v, wantErr %v", err, tt.wantErr)
				}
			}()
			if got := ConvertToInterfaceArray(tt.args.from); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ConvertToInterfaceArray() = %v, want %v", got, tt.want)
			}
		})
	}
}
