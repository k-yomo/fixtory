// Code generated by fixtory; DO NOT EDIT.

package example

import (
	"github.com/k-yomo/fixtory"
	"testing"
)

type AuthorFactory interface {
	NewBuilder(bluePrint AuthorBluePrintFunc, traits ...Author) AuthorBuilder
	OnBuild(onBuild func(t *testing.T, author *Author))
	Reset()
}

type AuthorBuilder interface {
	EachParam(authorParams ...Author) AuthorBuilder
	Zero(authorFields ...AuthorField) AuthorBuilder
	ResetAfter() AuthorBuilder

	Build() *Author
	Build2() (*Author, *Author)
	Build3() (*Author, *Author, *Author)
	BuildList(n int) []*Author
}

type AuthorBluePrintFunc func(i int, last Author) Author

type AuthorField string

const (
	AuthorIDField   AuthorField = "ID"
	AuthorNameField AuthorField = "Name"
)

type authorFactory struct {
	t       *testing.T
	factory *fixtory.Factory
}

type authorBuilder struct {
	t       *testing.T
	builder *fixtory.Builder
}

func NewAuthorFactory(t *testing.T) AuthorFactory {
	t.Helper()

	return &authorFactory{t: t, factory: fixtory.NewFactory(t, Author{})}
}

func (uf *authorFactory) NewBuilder(bluePrint AuthorBluePrintFunc, authorTraits ...Author) AuthorBuilder {
	uf.t.Helper()

	var bp fixtory.BluePrintFunc
	if bluePrint != nil {
		bp = func(i int, last interface{}) interface{} { return bluePrint(i, last.(Author)) }
	}
	builder := uf.factory.NewBuilder(bp, fixtory.ConvertToInterfaceArray(authorTraits)...)

	return &authorBuilder{t: uf.t, builder: builder}
}

func (uf *authorFactory) OnBuild(onBuild func(t *testing.T, author *Author)) {
	uf.t.Helper()

	uf.factory.OnBuild = func(t *testing.T, v interface{}) { onBuild(t, v.(*Author)) }
}

func (uf *authorFactory) Reset() {
	uf.t.Helper()

	uf.factory.Reset()
}

func (ub *authorBuilder) Zero(authorFields ...AuthorField) AuthorBuilder {
	ub.t.Helper()

	fields := make([]string, 0, len(authorFields))
	for _, f := range authorFields {
		fields = append(fields, string(f))
	}

	ub.builder = ub.builder.Zero(fields...)
	return ub
}
func (ub *authorBuilder) ResetAfter() AuthorBuilder {
	ub.t.Helper()

	ub.builder = ub.builder.ResetAfter()
	return ub
}

func (ub *authorBuilder) EachParam(authorParams ...Author) AuthorBuilder {
	ub.t.Helper()

	ub.builder = ub.builder.EachParam(fixtory.ConvertToInterfaceArray(authorParams)...)
	return ub
}

func (ub *authorBuilder) Build() *Author {
	ub.t.Helper()

	return ub.builder.Build().(*Author)
}

func (ub *authorBuilder) Build2() (*Author, *Author) {
	ub.t.Helper()

	list := ub.BuildList(2)
	return list[0], list[1]
}

func (ub *authorBuilder) Build3() (*Author, *Author, *Author) {
	ub.t.Helper()

	list := ub.BuildList(3)
	return list[0], list[1], list[2]
}

func (ub *authorBuilder) BuildList(n int) []*Author {
	ub.t.Helper()

	authors := make([]*Author, 0, n)
	for _, author := range ub.builder.BuildList(n) {
		authors = append(authors, author.(*Author))
	}
	return authors
}

type ArticleFactory interface {
	NewBuilder(bluePrint ArticleBluePrintFunc, traits ...Article) ArticleBuilder
	OnBuild(onBuild func(t *testing.T, article *Article))
	Reset()
}

type ArticleBuilder interface {
	EachParam(articleParams ...Article) ArticleBuilder
	Zero(articleFields ...ArticleField) ArticleBuilder
	ResetAfter() ArticleBuilder

	Build() *Article
	Build2() (*Article, *Article)
	Build3() (*Article, *Article, *Article)
	BuildList(n int) []*Article
}

type ArticleBluePrintFunc func(i int, last Article) Article

type ArticleField string

const (
	ArticleIDField                 ArticleField = "ID"
	ArticleTitleField              ArticleField = "Title"
	ArticleBodyField               ArticleField = "Body"
	ArticleAuthorIDField           ArticleField = "AuthorID"
	ArticlePublishScheduledAtField ArticleField = "PublishScheduledAt"
	ArticlePublishedAtField        ArticleField = "PublishedAt"
	ArticleStatusField             ArticleField = "Status"
	ArticleLikeCountField          ArticleField = "LikeCount"
)

type articleFactory struct {
	t       *testing.T
	factory *fixtory.Factory
}

type articleBuilder struct {
	t       *testing.T
	builder *fixtory.Builder
}

func NewArticleFactory(t *testing.T) ArticleFactory {
	t.Helper()

	return &articleFactory{t: t, factory: fixtory.NewFactory(t, Article{})}
}

func (uf *articleFactory) NewBuilder(bluePrint ArticleBluePrintFunc, articleTraits ...Article) ArticleBuilder {
	uf.t.Helper()

	var bp fixtory.BluePrintFunc
	if bluePrint != nil {
		bp = func(i int, last interface{}) interface{} { return bluePrint(i, last.(Article)) }
	}
	builder := uf.factory.NewBuilder(bp, fixtory.ConvertToInterfaceArray(articleTraits)...)

	return &articleBuilder{t: uf.t, builder: builder}
}

func (uf *articleFactory) OnBuild(onBuild func(t *testing.T, article *Article)) {
	uf.t.Helper()

	uf.factory.OnBuild = func(t *testing.T, v interface{}) { onBuild(t, v.(*Article)) }
}

func (uf *articleFactory) Reset() {
	uf.t.Helper()

	uf.factory.Reset()
}

func (ub *articleBuilder) Zero(articleFields ...ArticleField) ArticleBuilder {
	ub.t.Helper()

	fields := make([]string, 0, len(articleFields))
	for _, f := range articleFields {
		fields = append(fields, string(f))
	}

	ub.builder = ub.builder.Zero(fields...)
	return ub
}
func (ub *articleBuilder) ResetAfter() ArticleBuilder {
	ub.t.Helper()

	ub.builder = ub.builder.ResetAfter()
	return ub
}

func (ub *articleBuilder) EachParam(articleParams ...Article) ArticleBuilder {
	ub.t.Helper()

	ub.builder = ub.builder.EachParam(fixtory.ConvertToInterfaceArray(articleParams)...)
	return ub
}

func (ub *articleBuilder) Build() *Article {
	ub.t.Helper()

	return ub.builder.Build().(*Article)
}

func (ub *articleBuilder) Build2() (*Article, *Article) {
	ub.t.Helper()

	list := ub.BuildList(2)
	return list[0], list[1]
}

func (ub *articleBuilder) Build3() (*Article, *Article, *Article) {
	ub.t.Helper()

	list := ub.BuildList(3)
	return list[0], list[1], list[2]
}

func (ub *articleBuilder) BuildList(n int) []*Article {
	ub.t.Helper()

	articles := make([]*Article, 0, n)
	for _, article := range ub.builder.BuildList(n) {
		articles = append(articles, article.(*Article))
	}
	return articles
}
