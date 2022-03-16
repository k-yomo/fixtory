package fixtory

import (
	"reflect"
	"testing"
)

type testStruct struct {
	String      string
	StringPtr   *string
	Int         int
	Float       float64
	Array       []int
	Map         map[string]bool
	ChildStruct *childStruct
}

type childStruct struct {
	String string
	Int    int
}

func TestNewFactory(t *testing.T) {
	type args struct {
		t *testing.T
		v testStruct
	}

	tests := []struct {
		name string
		args args
		want *Factory[testStruct]
	}{
		{
			name: "initializes new factory",
			args: args{
				t: t,
				v: testStruct{},
			},
			want: &Factory[testStruct]{
				t:           t,
				productType: reflect.PtrTo(reflect.TypeOf(testStruct{})),
				last:        testStruct{},
				index:       0,
				OnBuild:     nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewFactory(tt.args.t, tt.args.v); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewFactory() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBuilder_Build(t *testing.T) {
	fac := NewFactory(t, testStruct{})
	fac.OnBuild = func(t *testing.T, v *testStruct) {
		if v.Int == 0 {
			t.Errorf("OnBuild = %d, want not zero", v.Int)
		}
	}

	bluePrint := func(i int, last testStruct) testStruct {
		return testStruct{
			String: "setByBlueprint",
			Int:    last.Int + 1,
			Float:  0.5,
			Array:  []int{1, 2, 3},
			Map:    map[string]bool{"a": true},
			ChildStruct: &childStruct{
				String: "child",
				Int:    10,
			},
		}
	}

	tests := []struct {
		name    string
		builder *Builder[testStruct]
		want    *testStruct
	}{
		{
			name:    "struct can be initialized with nil blueprint",
			builder: fac.NewBuilder(nil, testStruct{Int: 5}).ResetAfter(),
			want:    &testStruct{Int: 5},
		},
		{
			name:    "struct is overwritten by traits, zero, each param",
			builder: fac.NewBuilder(bluePrint, testStruct{String: "setByTrait1", Int: 10}, testStruct{String: "setByTrait2", Array: []int{1, 2, 3}}).Zero("Map").EachParam(testStruct{Float: 10.9}),
			want: &testStruct{
				String: "setByTrait2",
				Int:    10,
				Float:  10.9,
				Array:  []int{1, 2, 3},
				Map:    nil,
				ChildStruct: &childStruct{
					String: "child",
					Int:    10,
				},
			},
		},
		{
			name:    "empty fields do not overwrite",
			builder: fac.NewBuilder(bluePrint, testStruct{}).EachParam(testStruct{}),
			want: &testStruct{
				String: "setByBlueprint",
				Int:    11,
				Float:  0.5,
				Array:  []int{1, 2, 3},
				Map:    map[string]bool{"a": true},
				ChildStruct: &childStruct{
					String: "child",
					Int:    10,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.builder.Build(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Build() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFactory_OnBuild(t *testing.T) {
	fac := NewFactory(t, testStruct{})

	want := testStruct{
		String:      "a",
		StringPtr:   func() *string { s := "a"; return &s }(),
		Int:         5,
		Float:       0,
		Array:       nil,
		Map:         nil,
		ChildStruct: nil,
	}

	t.Run("onBuild is called after product is all set", func(t *testing.T) {
		fac.OnBuild = func(t *testing.T, got *testStruct) {
			if !reflect.DeepEqual(got, &want) {
				t.Errorf("OnBuild() = %v, want %v", got, &want)
			}
		}
		fac.NewBuilder(nil, want).Build()
	})
}

func TestBuilder_BuildList(t *testing.T) {
	fac := NewFactory(t, testStruct{})

	type args struct {
		n int
	}
	tests := []struct {
		name    string
		builder *Builder[testStruct]
		args    args
		want    []*testStruct
	}{
		{
			name: "initializes struct list with reset",
			builder: fac.NewBuilder(func(i int, last testStruct) testStruct {
				lastChild := &childStruct{}
				if last.ChildStruct != nil {
					lastChild = last.ChildStruct
				}
				return testStruct{
					String: "test",
					Int:    i + 1,
					Map:    map[string]bool{"a": true},
					ChildStruct: &childStruct{
						String: lastChild.String + "a",
					},
				}
			}, testStruct{}).
				EachParam(testStruct{Float: 0.1}, testStruct{Float: 0.2}, testStruct{Float: 0.3}).
				Zero("Map").
				ResetAfter(),
			args: args{n: 3},
			want: []*testStruct{
				{String: "test", Int: 1, Float: 0.1, ChildStruct: &childStruct{String: "a"}},
				{String: "test", Int: 2, Float: 0.2, ChildStruct: &childStruct{String: "aa"}},
				{String: "test", Int: 3, Float: 0.3, ChildStruct: &childStruct{String: "aaa"}},
			},
		},
		{
			name: "initialize struct list with initial index 0 (since index is reset on above test)",
			builder: fac.NewBuilder(func(i int, last testStruct) testStruct {
				return testStruct{Int: i + 1}
			}).
				EachParam(testStruct{Float: 0.1}, testStruct{}, testStruct{Float: 0.3}).
				Zero("Map").
				ResetAfter(),
			args: args{n: 3},
			want: []*testStruct{
				{Int: 1, Float: 0.1},
				{Int: 2},
				{Int: 3, Float: 0.3},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.builder.BuildList(tt.args.n); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BuildList() = %v, want %v", got, tt.want)
			}
		})
	}
}
