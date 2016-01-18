package api

import (
	"reflect"
	"strconv"
	"strings"

	"github.com/jinzhu/gorm"
)

const (
	perPage = 20
	allCat  = "all"
)

type QueryParams struct {
	CategoryStr string
	Categories  []string
	PageStr     string
	Page        int64
	Sort        string
}

func (p *QueryParams) Parse() {
	p.Categories = deleteEmpty(strings.Split(p.CategoryStr, ","))

	var err error
	p.Page, err = strconv.ParseInt(p.PageStr, 10, 64)
	if err != nil {
		p.Page = 0
	}
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
