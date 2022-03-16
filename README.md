![Fixtory Logo](https://user-images.githubusercontent.com/24503508/89726870-a4803980-da5a-11ea-9b84-d06eb73c7fdf.png)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://github.com/k-yomo/fixtory/blob/master/LICENSE)
[![Main Workflow](https://github.com/k-yomo/fixtory/workflows/Main%20Workflow/badge.svg)](https://github.com/k-yomo/fixtory/actions?query=workflow%3A%22Main+Workflow%22)
[![codecov](https://codecov.io/gh/k-yomo/fixtory/branch/master/graph/badge.svg)](https://codecov.io/gh/k-yomo/fixtory)
[![Go Report Card](https://goreportcard.com/badge/github.com/k-yomo/fixtory)](https://goreportcard.com/report/github.com/k-yomo/fixtory)

Fixtory is a test fixture factory which initializes type-safe, DRY, flexible fixtures with the power of Generics.

By using Fixtory...
- No more redundant repeated field value setting
- No more type assertion to convert from interface
- No more error handling when building fixture
- No more exhaustion to just prepare test data

## Installation
```sh
$ go get github.com/k-yomo/fixtory/cmd/fixtory/v2
```

## Getting started
Complete code is in [example](example).

Use factory to initialize fixtures
```go
// Author represents article's author
type Author struct {
    ID   int
    Name string
}

// Article represents article
type Article struct {
    ID                 int
    Title              string
    Body               string
    AuthorID           int
    PublishScheduledAt time.Time
    PublishedAt        time.Time
    Status             ArticleStatus
    LikeCount          int
}

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


func TestArticleList_SelectAuthoredBy(t *testing.T) {
    authorFactory := fixtory.NewFactory(t, Author{})
    articleFactory := fixtory.NewFactory(t, Article{})

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
```

## How it works
There are 4 layers in the process of initializing fixture in fixtory. 
The layers are stacked and each one is a delta of the changes from the previous layer like Dockerfile.

### 1. Blueprint
Blueprint is the base of fixture(like FROM in Dockerfile) and called first.
You need to implement blueprint function to meet generated blueprint type (like below)
It should return instance with generic field values.
```
type TestArticleBluePrintFunc func(i int, last Article) Article
```

### 2. Traits
To overwrite some fields, you can use traits.
Traits are applied in the order of arguments to all fixtures.
※ Only non-zero value will be set.
```go
//  Repeatedly used trait would be better to define as global variable.
var articleTraitPublished = Article{
	Status:             ArticleStatusOpen,
	PublishScheduledAt: time.Now().Add(-1 * time.Hour),
	PublishedAt:        time.Now().Add(-1 * time.Hour),
	LikeCount:          15,
}

// recently published articles
articles:= articleFactory.NewBuilder(
               nil, 
               articleTraitPublished,
               Article{AuthorID: 5, PublishedAt: time.Now().Add(-1 * time.Minute)},
           ).BuildList(2)
```

### 3. Each Param
When you want to overwrite a specific fixture in the list, use EachParam.
Each Param overwrites the same index struct as parameter.
※ Only non-zero value will be set.
```go
articleFactory := NewArticleFactory(t)
articles := articleFactory.NewBuilder(nil, Article{Title: "test article"})
                .EachParam(Article{AuthorID: 1}, Article{AuthorID: 2}, Article{AuthorID: 2})
                .BuildList(3)
```

### 4. Zero
Since there is no way to distinguish default zero value or intentionally set zero in params,
you can overwrite fields with zero value like below, and it will be applied at the last minute.
```go
articleFactory := NewArticleFactory(t)
// AuthorID will be overwritten with zero value.
articles := articleFactory.NewBuilder(articleBluePrint, Article{AuthorID: author1.ID}).
                Zero(ArticleAuthorIDField).
                BuildList(4)
```
