package fixtory

import (
	"reflect"
	"testing"
)

type BluePrintFunc func(i int, last interface{}) interface{}

type Factory struct {
	t           *testing.T
	productType reflect.Type
	last        interface{}
	index       int
	OnBuild     func(t *testing.T, v interface{})
}

type Builder struct {
	*Factory
	index           int
	bluePrint       func(i int, last interface{}) interface{}
	traits          []interface{}
	eachParams      []interface{}
	zeroFields      []string
	resetAfterBuild bool
}

func NewFactory(t *testing.T, v interface{}) *Factory {
	return &Factory{t: t, productType: reflect.TypeOf(v), index: 0, last: v}
}

func (uf *Factory) NewBuilder(bluePrint BluePrintFunc, traits ...interface{}) *Builder {
	return &Builder{Factory: uf, bluePrint: bluePrint, traits: traits}
}

func (uf *Factory) Reset() {
	uf.last = nil
	uf.index = 0
}

func (b *Builder) WithZero(fields ...string) *Builder {
	b.zeroFields = fields
	return b
}
func (b *Builder) WithReset() *Builder {
	b.resetAfterBuild = true
	return b
}

func (b *Builder) WithEachParams(traits ...interface{}) *Builder {
	b.eachParams = traits
	return b
}

func (b *Builder) Build() interface{} {
	return b.build()
}

func (b *Builder) BuildList(n int) []interface{} {
	b.index = 0
	users := make([]interface{}, 0, n)
	for i := 0; i < n; i++ {
		users = append(users, b.build())
	}
	return users
}

func (b *Builder) build() interface{} {
	product := reflect.New(b.productType).Interface()

	if b.bluePrint != nil {
		product = b.bluePrint(b.index, b.last)
	}
	for _, trait := range b.traits {
		MapNotZeroFields(trait, product)
	}
	if len(b.eachParams) > b.index {
		MapNotZeroFields(b.eachParams[b.index], product)
	}
	for _, f := range b.zeroFields {
		uf := reflect.ValueOf(product).Elem().FieldByName(f)
		uf.Set(reflect.Zero(uf.Type()))
	}

	b.last = product
	b.index++
	b.Factory.index++

	if b.OnBuild != nil {
		b.OnBuild(b.t, product)
	}
	return product
}
