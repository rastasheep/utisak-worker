package api

import "time"

type SerializedArticle struct {
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Date      time.Time `json:"published_at"`

	Title        string `json:"title"`
	Domain       string `json:"domain"`
	Url          string `json:"url"`
	Author       string `json:"-"`
	Excerpt      string `json:"excerpt"`
	WordCount    int    `json:"word_count"`
	LeadImage    string `json:"hero_image_url"`
	Category     string `json:"category"`
	CategorySlug string `json:"category_slug"`
	Source       string `json:"author"`
}

func (sa SerializedArticle) TableName() string {
	return "articles"
}
