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

func (article *Article) ParseData(data map[string]interface{}) error {
	//result.Title = data["Title"].(string)
	article.ID = data["ID"].(string)
	article.Excerpt = data["Summary"].(string)
	article.Date, _ = time.Parse(time.RFC3339, data["Date"].(string))

	return ReadabilityParse(data["Link"].(string), article)
}
