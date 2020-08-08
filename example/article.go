package example

import (
	"time"
)

type ArticleStatus int

const (
	ArticleStatusDraft ArticleStatus = iota + 1
	ArticleStatusWaitReview
	ArticleStatusOpen
)

//go:generate fixtory -type=Author,Article
type Author struct {
	ID   int
	Name string
}

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

type ArticleList []*Article

func (list ArticleList) SelectPublished() ArticleList {
	var publishedArticles ArticleList
	for _, a := range list {
		if a.Status == ArticleStatusOpen && a.PublishScheduledAt.Before(time.Now()) {
			publishedArticles = append(publishedArticles, a)
		}
	}
	return publishedArticles
}

func (list ArticleList) SelectAuthoredBy(authorID int) ArticleList {
	var authoredArticles ArticleList
	for _, a := range list {
		if a.AuthorID == authorID {
			authoredArticles = append(authoredArticles, a)
		}
	}
	return authoredArticles
}
