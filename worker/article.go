package main

import "time"

type Article struct {
	ID        string
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
