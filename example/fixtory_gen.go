// Code generated by fixtory; DO NOT EDIT.

package example

import (
	"github.com/k-yomo/fixtory"
	"testing"
)

type TestAuthorFactory interface {
	NewBuilder(bluePrint TestAuthorBluePrintFunc, traits ...*Author) TestAuthorBuilder
	OnBuild(onBuild func(t *testing.T, author *Author))
	Reset()
}

type TestAuthorBuilder interface {
	Build() *Author
	Build2() (*Author, *Author)
	Build3() (*Author, *Author, *Author)
	BuildList(n int) []*Author
	WithZero(authorFields ...TestAuthorField) TestAuthorBuilder
	WithReset() TestAuthorBuilder
	WithEachParams(authorTraits ...*Author) TestAuthorBuilder
}

type TestAuthorBluePrintFunc func(i int, last *Author) *Author

type TestAuthorField string

const (
	TestAuthorID   TestAuthorField = "ID"
	TestAuthorName TestAuthorField = "Name"
)

type testAuthorFactory struct {
	t       *testing.T
	factory *fixtory.Factory
}

type testAuthorBuilder struct {
	t       *testing.T
	builder *fixtory.Builder
}

func TestNewAuthorFactory(t *testing.T) TestAuthorFactory {
	t.Helper()

	return &testAuthorFactory{t: t, factory: fixtory.NewFactory(t, &Author{})}
}

func (uf *testAuthorFactory) NewBuilder(bluePrint TestAuthorBluePrintFunc, authorTraits ...*Author) TestAuthorBuilder {
	uf.t.Helper()

	var bp fixtory.BluePrintFunc
	if bluePrint != nil {
		bp = func(i int, last interface{}) interface{} { return bluePrint(i, last.(*Author)) }
	}
	builder := uf.factory.NewBuilder(bp, fixtory.ConvertToInterfaceArray(authorTraits)...)

	return &testAuthorBuilder{t: uf.t, builder: builder}
}

func (uf *testAuthorFactory) OnBuild(onBuild func(t *testing.T, author *Author)) {
	uf.t.Helper()

	uf.factory.OnBuild = func(t *testing.T, v interface{}) { onBuild(t, v.(*Author)) }
}

func (uf *testAuthorFactory) Reset() {
	uf.t.Helper()

	uf.factory.Reset()
}

func (ub *testAuthorBuilder) WithZero(authorFields ...TestAuthorField) TestAuthorBuilder {
	ub.t.Helper()

	fields := make([]string, 0, len(authorFields))
	for _, f := range authorFields {
		fields = append(fields, string(f))
	}

	ub.builder = ub.builder.WithZero(fields...)
	return ub
}
func (ub *testAuthorBuilder) WithReset() TestAuthorBuilder {
	ub.t.Helper()

	ub.builder = ub.builder.WithReset()
	return ub
}

func (ub *testAuthorBuilder) WithEachParams(authorTraits ...*Author) TestAuthorBuilder {
	ub.t.Helper()

	ub.builder = ub.builder.WithEachParams(fixtory.ConvertToInterfaceArray(authorTraits)...)
	return ub
}

func (ub *testAuthorBuilder) Build() *Author {
	ub.t.Helper()

	return ub.builder.Build().(*Author)
}

func (ub *testAuthorBuilder) Build2() (*Author, *Author) {
	ub.t.Helper()

	return ub.Build(), ub.Build()
}

func (ub *testAuthorBuilder) Build3() (*Author, *Author, *Author) {
	ub.t.Helper()

	return ub.Build(), ub.Build(), ub.Build()
}

func (ub *testAuthorBuilder) BuildList(n int) []*Author {
	ub.t.Helper()

	authors := make([]*Author, 0, n)
	for i := 0; i < n; i++ {
		authors = append(authors, ub.builder.Build().(*Author))
	}
	return authors
}

type TestArticleFactory interface {
	NewBuilder(bluePrint TestArticleBluePrintFunc, traits ...*Article) TestArticleBuilder
	OnBuild(onBuild func(t *testing.T, article *Article))
	Reset()
}

type TestArticleBuilder interface {
	Build() *Article
	Build2() (*Article, *Article)
	Build3() (*Article, *Article, *Article)
	BuildList(n int) []*Article
	WithZero(articleFields ...TestArticleField) TestArticleBuilder
	WithReset() TestArticleBuilder
	WithEachParams(articleTraits ...*Article) TestArticleBuilder
}

type TestArticleBluePrintFunc func(i int, last *Article) *Article

type TestArticleField string

const (
	TestArticleID                 TestArticleField = "ID"
	TestArticleTitle              TestArticleField = "Title"
	TestArticleBody               TestArticleField = "Body"
	TestArticleAuthorID           TestArticleField = "AuthorID"
	TestArticlePublishScheduledAt TestArticleField = "PublishScheduledAt"
	TestArticlePublishedAt        TestArticleField = "PublishedAt"
	TestArticleStatus             TestArticleField = "Status"
	TestArticleLikeCount          TestArticleField = "LikeCount"
)

type testArticleFactory struct {
	t       *testing.T
	factory *fixtory.Factory
}

type testArticleBuilder struct {
	t       *testing.T
	builder *fixtory.Builder
}

func TestNewArticleFactory(t *testing.T) TestArticleFactory {
	t.Helper()

	return &testArticleFactory{t: t, factory: fixtory.NewFactory(t, &Article{})}
}

func (uf *testArticleFactory) NewBuilder(bluePrint TestArticleBluePrintFunc, articleTraits ...*Article) TestArticleBuilder {
	uf.t.Helper()

	var bp fixtory.BluePrintFunc
	if bluePrint != nil {
		bp = func(i int, last interface{}) interface{} { return bluePrint(i, last.(*Article)) }
	}
	builder := uf.factory.NewBuilder(bp, fixtory.ConvertToInterfaceArray(articleTraits)...)

	return &testArticleBuilder{t: uf.t, builder: builder}
}

func (uf *testArticleFactory) OnBuild(onBuild func(t *testing.T, article *Article)) {
	uf.t.Helper()

	uf.factory.OnBuild = func(t *testing.T, v interface{}) { onBuild(t, v.(*Article)) }
}

func (uf *testArticleFactory) Reset() {
	uf.t.Helper()

	uf.factory.Reset()
}

func (ub *testArticleBuilder) WithZero(articleFields ...TestArticleField) TestArticleBuilder {
	ub.t.Helper()

	fields := make([]string, 0, len(articleFields))
	for _, f := range articleFields {
		fields = append(fields, string(f))
	}

	ub.builder = ub.builder.WithZero(fields...)
	return ub
}
func (ub *testArticleBuilder) WithReset() TestArticleBuilder {
	ub.t.Helper()

	ub.builder = ub.builder.WithReset()
	return ub
}

func (ub *testArticleBuilder) WithEachParams(articleTraits ...*Article) TestArticleBuilder {
	ub.t.Helper()

	ub.builder = ub.builder.WithEachParams(fixtory.ConvertToInterfaceArray(articleTraits)...)
	return ub
}

func (ub *testArticleBuilder) Build() *Article {
	ub.t.Helper()

	return ub.builder.Build().(*Article)
}

func (ub *testArticleBuilder) Build2() (*Article, *Article) {
	ub.t.Helper()

	return ub.Build(), ub.Build()
}

func (ub *testArticleBuilder) Build3() (*Article, *Article, *Article) {
	ub.t.Helper()

	return ub.Build(), ub.Build(), ub.Build()
}

func (ub *testArticleBuilder) BuildList(n int) []*Article {
	ub.t.Helper()

	articles := make([]*Article, 0, n)
	for i := 0; i < n; i++ {
		articles = append(articles, ub.builder.Build().(*Article))
	}
	return articles
}
