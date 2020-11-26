package fixtory

import (
	"reflect"
	"testing"
)

type BluePrintFunc[T any] func(i int, last T) T

type Factory[T any] struct {
	t           *testing.T
	productType reflect.Type
	// struct
	last T
	// index is the next struct index (which is equal to already built struct count in this factory)
	index int
	// v is a pointer to struct
	OnBuild func(t *testing.T, v *T)
}

type Builder[T any] struct {
	factory *Factory[T]
	// index is the next struct index in this builder
	index           int
	bluePrint       func(i int, last T) T
	traits          []T
	eachParam       []T
	zeroFields      []string
	resetAfterBuild bool
}

func NewFactory[T any](t *testing.T, v T) *Factory[T] {
	return &Factory[T]{t: t, productType: reflect.PtrTo(reflect.TypeOf(v)), index: 0, last: v}
}

func (f *Factory[T]) NewBuilder(bluePrint BluePrintFunc[T], traits ...T) *Builder[T] {
	return &Builder[T]{factory: f, bluePrint: bluePrint, traits: traits}
}

func (f *Factory[T]) Reset() {
	f.last = reflect.New(f.productType.Elem()).Elem().Interface().(T)
	f.index = 0
}

func (b *Builder[T]) EachParam(params ...T) *Builder[T] {
	b.eachParam = params
	return b
}

func (b *Builder[T]) Zero(fields ...string) *Builder[T] {
	b.zeroFields = fields
	return b
}

func (b *Builder[T]) ResetAfter() *Builder[T] {
	b.resetAfterBuild = true
	return b
}

func (b *Builder[T]) Build() *T {
	b.index = 0
	product := b.build()
	if b.resetAfterBuild {
		b.factory.Reset()
	}
	return product
}

func (b *Builder[T]) Build2() (*T, *T) {
	list := b.BuildList(2)
	return list[0], list[1]
}

func (b *Builder[T]) Build3() (*T, *T, *T) {
	list := b.BuildList(3)
	return list[0], list[1], list[2]
}

func (b *Builder[T]) BuildList(n int) []*T {
	b.index = 0
	products := make([]*T, 0, n)
	for i := 0; i < n; i++ {
		products = append(products, b.build())
	}
	if b.resetAfterBuild {
		b.factory.Reset()
	}
	return products
}

func (b *Builder[T]) build() *T {
	product := reflect.New(b.factory.productType.Elem()).Interface().(*T)

	if b.bluePrint != nil {
		MapNotZeroFields(b.bluePrint(b.factory.index, b.factory.last), product)
	}
	for _, trait := range b.traits {
		MapNotZeroFields(trait, product)
	}
	if len(b.eachParam) > b.index {
		MapNotZeroFields(b.eachParam[b.index], product)
	}
	for _, zf := range b.zeroFields {
		f := reflect.ValueOf(product).Elem().FieldByName(zf)
		f.Set(reflect.Zero(f.Type()))
	}

	b.factory.last = reflect.ValueOf(product).Elem().Interface().(T)
	b.index++
	b.factory.index++

	if b.factory.OnBuild != nil {
		b.factory.OnBuild(b.factory.t, product)
	}
	return product
}
