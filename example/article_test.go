package example

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

var authorBluePrint = func(i int, last Author) Author {
	num := i + 1
	return Author{
		ID:   num,
		Name: fmt.Sprintf("Author %d", num),
	}
}

var articleBluePrint = func(i int, last Article) Article {
	num := i + 1
	return Article{
		ID:                 num,
		Title:              fmt.Sprintf("Article %d", i+1),
		AuthorID:           num,
		PublishScheduledAt: time.Now().Add(-1 * time.Hour),
		PublishedAt:        time.Now().Add(-1 * time.Hour),
		LikeCount:          15,
	}
}

var articleTraitDraft = Article{
	Status: ArticleStatusDraft,
}

var articleTraitPublishScheduled = Article{
	Status:             ArticleStatusOpen,
	PublishScheduledAt: time.Now().Add(1 * time.Hour),
}

var articleTraitPublished = Article{
	Status:             ArticleStatusOpen,
	PublishScheduledAt: time.Now().Add(-1 * time.Hour),
	PublishedAt:        time.Now().Add(-1 * time.Hour),
	LikeCount:          15,
}

func TestArticleList_SelectPublished(t *testing.T) {
	articleFactory := TestNewArticleFactory(t)
	// if you want to persist articles, set OnBuild func here
	articleFactory.OnBuild(func(t *testing.T, article *Article) { fmt.Println("Insert to db here") })

	// creates 3 different articles
	waitReview, publishedScheduled, published := articleFactory.NewBuilder(articleBluePrint).
		WithEachParams(articleTraitDraft, articleTraitPublishScheduled, articleTraitPublished).
		WithZero(TestArticleLikeCount).WithReset().
		Build3()

	tests := []struct {
		name string
		list ArticleList
		want ArticleList
	}{
		{
			name: "returns only published articles",
			list: ArticleList{waitReview, publishedScheduled, published},
			want: ArticleList{published},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.list.SelectPublished(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SelectPublished() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestArticleList_SelectAuthoredBy(t *testing.T) {
	authorFactory := TestNewAuthorFactory(t)
	articleFactory := TestNewArticleFactory(t)

	author1, author2 := authorFactory.NewBuilder(authorBluePrint).Build2()
	articlesAuthoredBy1 := articleFactory.NewBuilder(articleBluePrint, Article{AuthorID: author1.ID}).BuildList(4)
	articleAuthoredBy2 := articleFactory.NewBuilder(articleBluePrint, Article{AuthorID: author2.ID}).Build()

	type args struct {
		authorID int
	}
	tests := []struct {
		name string
		list ArticleList
		args args
		want ArticleList
	}{
		{
			name: "returns articles authored by author 1",
			list: append(articlesAuthoredBy1, articleAuthoredBy2),
			args: args{authorID: author1.ID},
			want: articlesAuthoredBy1,
		},
		{
			name: "returns articles authored by author 2",
			list: append(articlesAuthoredBy1, articleAuthoredBy2),
			args: args{authorID: author2.ID},
			want: ArticleList{articleAuthoredBy2},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.list.SelectAuthoredBy(tt.args.authorID); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SelectAuthoredBy() = %v, want %v", got, tt.want)
			}
		})
	}
}
