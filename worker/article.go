package main

import "time"

type Article struct {
	// gorm.Model
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time

	Title     string
	Domain    string
	Url       string
	Author    string
	Excerpt   string
	WordCount int `json:"word_count"`
	//Content    string
	Date      time.Time
	LeadImage string `json:"lead_image_url"`
}

func (article *Article) FetchDetails() error {
	return ReadabilityParse(article.Url, article)
}

func LatestArticle() *Article {
	var article Article

	db.Order("date desc").Limit(1).Find(&article)
	return &article
}
