package worker

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
	Date      time.Time `sql:"index:idx_category_source"`
	LeadImage string    `json:"lead_image_url"`
	Catogory  string    `sql:"index:idx_category_source"`
	Source    string    `sql:"index:idx_category_source"`
}

func (article *Article) FetchDetails() error {
	return ReadabilityParse(article.Url, article)
}
