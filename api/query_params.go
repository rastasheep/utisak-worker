package api

import (
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

const (
	perPage    = 20
	allCat     = "all"
	dateFormat = "2006-01-02"
)

type QueryParams struct {
	CategoryStr  string
	Categories   []string
	PageStr      string
	Page         int64
	Sort         string
	StartDateStr string
	StartDate    time.Time
	EndDateStr   string
	EndDate      time.Time
}

func (p *QueryParams) Parse() {
	p.Categories = deleteEmpty(strings.Split(p.CategoryStr, ","))

	var err error
	p.Page, err = strconv.ParseInt(p.PageStr, 10, 64)
	if err != nil {
		p.Page = 0
	}

	p.StartDate, _ = time.Parse(dateFormat, p.StartDateStr)
	p.EndDate, _ = time.Parse(dateFormat, p.EndDateStr)
}

func (p *QueryParams) PrepareQuery(searchRelation *gorm.DB) *gorm.DB {
	if len(p.Categories) != 0 && !contains(p.Categories, allCat) {
		searchRelation = searchRelation.Where("category_slug IN (?)", p.Categories)
	}

	if p.Page > 1 {
		searchRelation = searchRelation.Offset((p.Page - 1) * perPage)
	}

	if p.Sort == "newest" {
		searchRelation = searchRelation.Order("date desc")
	} else {
		searchRelation = searchRelation.Order("(total_views / POW(((EXTRACT(EPOCH FROM (now()-date)) / 3600)::integer + 2), 1.5)) desc")
	}

	if !p.StartDate.IsZero() {
		searchRelation = searchRelation.Where("date > ?", p.StartDate)
	}
	if !p.EndDate.IsZero() {
		searchRelation = searchRelation.Where("date < ?", p.EndDate)
	}

	return searchRelation.Limit(perPage)
}

func deleteEmpty(s []string) []string {
	var r []string
	if len(s) == 0 {
		return s
	}
	for _, str := range s {
		if str != "" {
			r = append(r, str)
		}
	}
	return r
}

func contains(slice interface{}, val interface{}) bool {
	sv := reflect.ValueOf(slice)

	for i := 0; i < sv.Len(); i++ {
		if sv.Index(i).Interface() == val {
			return true
		}
	}
	return false
}
