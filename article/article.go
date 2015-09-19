package article

import (
	"net/url"
	"time"
)

type Article struct {
	// gorm.Model
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time

	Title     string `json:"-"`
	Domain    string
	Url       string
	Author    string
	Excerpt   string
	WordCount int `json:"word_count"`
	//Content    string
	Date         time.Time `sql:"index:idx_category_source"`
	LeadImage    string    `json:"lead_image_url"`
	Category     string
	CategorySlug string `sql:"index:idx_category_source"`
	Source       string `sql:"index:idx_category_source"`
	TotalViews   int    `sql:"default:0"`
}

func (a *Article) BeforeCreate() (err error) {
	u, _ := url.Parse(a.LeadImage)
	a.LeadImage = u.String()
	return
}
