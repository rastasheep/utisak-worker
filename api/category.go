package api

type Category struct {
	Category     string `json:"title"`
	CategorySlug string `json:"slug"`
	Count        int    `json:"count"`
}

func (c Category) TableName() string {
	return "articles"
}

func GetCategories() []Category {
	// SELECT category, category_slug, COUNT(category)
	// FROM "articles"
	// WHERE DATE(created_at) = DATE(NOW())
	// GROUP BY category, category_slug;

	var categories []Category

	db.Select("category, category_slug, COUNT(category)").
		Where("DATE(created_at) = DATE(NOW())").
		Group("category, category_slug").
		Find(&categories)

	return categories
}
