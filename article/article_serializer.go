package article

import (
	"fmt"
	"time"
)

const ShareUrlTmpl = "http://local.go:8080/posts/%d"

type SerializedArticle struct {
	ID        uint      `json:"-"`
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
	ShareUrl     string `json:"share_url"`
	TotalViews   int    `json:"total_views" `
}

func (sa SerializedArticle) TableName() string {
	return "articles"
}

func (sa *SerializedArticle) AfterFind() (err error) {
	sa.ShareUrl = fmt.Sprintf(ShareUrlTmpl, sa.ID)
	return
}
