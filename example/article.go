package example

import (
	"time"
)

// ArticleStatus is the status of article
type ArticleStatus int

const (
	ArticleStatusDraft ArticleStatus = iota + 1
	ArticleStatusOpen
)

//go:generate fixtory -type=Author,Article -output=article.fixtory.go
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

// Article represents list of article
type ArticleList []*Article

// SelectPublished returns only published articles
func (list ArticleList) SelectPublished() ArticleList {
	var publishedArticles ArticleList
	for _, a := range list {
		if a.Status == ArticleStatusOpen && a.PublishScheduledAt.Before(time.Now()) {
			publishedArticles = append(publishedArticles, a)
		}
	}
	return publishedArticles
}

// SelectAuthoredBy returns only articles authored by given author's id
func (list ArticleList) SelectAuthoredBy(authorID int) ArticleList {
	var authoredArticles ArticleList
	for _, a := range list {
		if a.AuthorID == authorID {
			authoredArticles = append(authoredArticles, a)
		}
	}
	return authoredArticles
}
