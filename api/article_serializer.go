package api

import "time"

type SerializedArticle struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	Date      time.Time `json:"published_at"`

	Title        string
	Domain       string
	Url          string
	Author       string
	Excerpt      string `json:"excerpt"`
	WordCount    int    `json:"word_count"`
	LeadImage    string `json:"hero_image_url"`
	Category     string
	CategorySlug string `json:"category_slug"`
	Source       string
}

func (sa SerializedArticle) TableName() string {
	return "articles"
}
