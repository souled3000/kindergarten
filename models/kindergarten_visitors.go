package models

import (
	"github.com/astaxie/beego/orm"
	"math"
)

type Page struct {
	Data    []orm.Params `json:"data"`
	Total   int64        `json:"total"`
	PageNum int          `json:"page_num"`
}

func GetVisitors(page int, limit int) (Page, error) {
	o := orm.NewOrm()
	var maps []orm.Params
	totalqb, _ := orm.NewQueryBuilder("mysql")

	tatolsql := totalqb.Select("count(*)").From("visitors").String()

	var total int64
	err := o.Raw(tatolsql).QueryRow(&total)
	if err == nil {
		if limit <= 0 {
			limit = 10
		}

		if page <= 0 {
			page = 1
		}

		offset := (page - 1) * limit

		qb, _ := orm.NewQueryBuilder("mysql")
		sql := qb.Select("*").
			From("visitors").
			OrderBy("visitor_id").
			Limit(int(limit)).
			Offset(int(offset)).
			String()

		if _, err := o.Raw(sql).Values(&maps); err == nil {
			pageNum := int(math.Ceil(float64(total) / float64(limit)))
			return Page{maps, total, pageNum}, nil
		}
	}
	return Page{}, nil
}
