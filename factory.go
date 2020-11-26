package fixtory

import (
	"github.com/k-yomo/fixtory/pkg/reflectutil"
	"reflect"
	"testing"
)

type BluePrintFunc func(i int, last interface{}) interface{}

type Factory struct {
	t           *testing.T
	productType reflect.Type
	// struct
	last interface{}
	// index is the next struct index (which is equal to already built struct count in this factory)
	index int
	// v is a pointer to struct
	OnBuild func(t *testing.T, v interface{})
}

type Builder struct {
	*Factory
	// index is the next struct index in this builder
	index           int
	bluePrint       func(i int, last interface{}) interface{}
	traits          []interface{}
	eachParam       []interface{}
	zeroFields      []string
	resetAfterBuild bool
}

func NewFactory(t *testing.T, v interface{}) *Factory {
	return &Factory{t: t, productType: reflect.PtrTo(reflect.TypeOf(v)), index: 0, last: v}
}

func (uf *Factory) NewBuilder(bluePrint BluePrintFunc, traits ...interface{}) *Builder {
	return &Builder{Factory: uf, bluePrint: bluePrint, traits: traits}
}

func (uf *Factory) Reset() {
	uf.last = reflect.New(uf.productType.Elem()).Elem().Interface()
	uf.index = 0
}

func (b *Builder) EachParam(params ...interface{}) *Builder {
	b.eachParam = params
	return b
}

func (b *Builder) Zero(fields ...string) *Builder {
	b.zeroFields = fields
	return b
}

func (b *Builder) ResetAfter() *Builder {
	b.resetAfterBuild = true
	return b
}

func (b *Builder) Build() interface{} {
	b.index = 0
	product := b.build()
	if b.resetAfterBuild {
		b.Factory.Reset()
	}
	return product
}

func (b *Builder) BuildList(n int) []interface{} {
	b.index = 0
	products := make([]interface{}, 0, n)
	for i := 0; i < n; i++ {
		products = append(products, b.build())
	}
	if b.resetAfterBuild {
		b.Factory.Reset()
	}
	return products
}

func (b *Builder) build() interface{} {
	product := reflect.New(b.productType.Elem()).Interface()

	if b.bluePrint != nil {
		reflectutil.MapNotZeroFields(b.bluePrint(b.Factory.index, b.last), product)
	}
	for _, trait := range b.traits {
		reflectutil.MapNotZeroFields(trait, product)
	}
	if len(b.eachParam) > b.index {
		reflectutil.MapNotZeroFields(b.eachParam[b.index], product)
	}
	for _, f := range b.zeroFields {
		uf := reflect.ValueOf(product).Elem().FieldByName(f)
		uf.Set(reflect.Zero(uf.Type()))
	}

	b.last = reflect.ValueOf(product).Elem().Interface()
	b.index++
	b.Factory.index++

	if b.OnBuild != nil {
		b.OnBuild(b.t, product)
	}
	return product
}
