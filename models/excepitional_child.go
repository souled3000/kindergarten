package models

import(
	"math"
	"github.com/astaxie/beego/orm"
	"time"
)




// GetSearch 查询特殊儿童列表或者搜索
// no records exist
func GetExceptionalChild(page int, limit int, keyword string) (Page, error) {
	where := "1=1 "

	// 特殊儿童姓名或者特殊儿童过敏源
	if keyword != "" {
		where += " AND child_name like \"%" + string(keyword) + "%\" OR allergen like \"%" + string(keyword) + "%\""

	}


	o := orm.NewOrm()
	var maps []orm.Params
	totalqb, _ := orm.NewQueryBuilder("mysql")

	tatolsql := totalqb.Select("count(*)").From("exceptional_child").Where(where).String()

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
		sql := qb.Select("id, child_name, somatotype, allergen, updated_at").
			From("exceptional_child").
			Where(where).
			OrderBy("id").
			Limit(int(limit)).
			Offset(int(offset)).
			String()

		if _, err := o.Raw(sql).Values(&maps); err == nil {
			var newMap []orm.Params
			for _, v := range maps {
				t := time.Now()
				currentTime := t.Unix() - (24 * 3600 * 3) //当前时间戳
				BeCharge := v["updated_at"]
				toBeCharge := BeCharge.(string)
				timeLayout := "2006-01-02 15:04:05"
				loc, _ := time.LoadLocation("Local")
				theTime, _ := time.ParseInLocation(timeLayout, toBeCharge, loc)
				sr := theTime.Unix() //将更新时间转为时间戳
				if sr > currentTime {
					v["new"] = 1 //最新
					v["class"] = "大班二班"
				} else {
					v["new"] = 0
					v["class"] = "大班一班"
				}

				delete(v, "updated_at")
				newMap = append(newMap, v)
			}

			pageNum := int(math.Ceil(float64(total) / float64(limit)))
			return Page{newMap, total, pageNum}, nil
		}
	}
	return Page{}, nil
}