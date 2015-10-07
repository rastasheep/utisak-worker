package article

import (
	"net/url"
	"time"

	"github.com/Machiel/slugify"
)

type Article struct {
	// gorm.Model
	ID        uint `sql:"index:idx_id_source_slug";gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time

	Title     string `json:"-"`
	Slug      string `sql:"index:idx_id_source_slug"`
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
	SourceSlug   string `sql:"index:idx_id_source_slug"`
	TotalViews   int    `sql:"default:0"`
}

func (a *Article) BeforeCreate() (err error) {
	u, err := url.Parse(a.LeadImage)
	if err != nil {
		a.LeadImage = ""
	}
	a.LeadImage = u.String()
	a.Slug = slugify.Slugify(a.Title)
	return
}
