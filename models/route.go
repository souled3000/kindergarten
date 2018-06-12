package models

import (
	"math"
	"time"

	"github.com/astaxie/beego/orm"
)

type Route struct {
	Id        int       `json:"id" orm:"column(id);auto"`
	Name      string    `json:"name" orm:"column(name);size(10)" description:"名称"`
	Route     string    `json:"route" orm:"column(route);size(50)" description:"路由"`
	CreatedAt time.Time `json:"created_at" orm:"auto_now_add"`
	UpdatedAt time.Time `json:"updated_at" orm:"auto_now"`
}

func (t *Route) TableName() string {
	return "route"
}

func init() {
	orm.RegisterModel(new(Route))
}

/*
添加路由
*/
func AddRoute(name string, route string) map[string]interface{} {
	o := orm.NewOrm()
	var r Route
	r.Name = name
	r.Route = route
	id, err := o.Insert(&r)
	if err == nil {
		paginatorMap := make(map[string]interface{})
		paginatorMap["data"] = id
		return paginatorMap
	}
	return nil
}

/*
路由详情
*/
func GetRouteById(id int) map[string]interface{} {
	o := orm.NewOrm()
	var v []orm.Params
	qb, _ := orm.NewQueryBuilder("mysql")
	// 构建查询对象
	sql := qb.Select("*").From("route").Where("id = ?").String()
	_, err := o.Raw(sql, id).Values(&v)
	if err == nil {
		paginatorMap := make(map[string]interface{})
		paginatorMap["data"] = v
		return paginatorMap
	}
	return nil
}

/*
路由列表
*/
func GetAllRoute(page int, prepage int) map[string]interface{} {
	o := orm.NewOrm()
	qb, _ := orm.NewQueryBuilder("mysql")
	// 构建查询对象
	sql := qb.Select("count(*)").From("route").String()
	var total int64
	err := o.Raw(sql).QueryRow(&total)
	if err == nil {
		var v []orm.Params
		//根据nums总数，和prepage每页数量 生成分页总数
		totalpages := int(math.Ceil(float64(total) / float64(prepage))) //page总数
		if page > totalpages {
			page = totalpages
		}
		if page <= 0 {
			page = 1
		}
		limit := (page - 1) * prepage
		qb, _ := orm.NewQueryBuilder("mysql")
		sql := qb.Select("*").From("route").Limit(prepage).Offset(limit).String()
		num, err := o.Raw(sql).Values(&v)
		if err == nil && num > 0 {
			paginatorMap := make(map[string]interface{})
			paginatorMap["total"] = total         //总条数
			paginatorMap["data"] = v              //分页数据
			paginatorMap["page_num"] = totalpages //总页数
			return paginatorMap
		}
	}
	return nil
}

/*
编辑路由
*/
func UpdateRouteById(id int, name string, route string) map[string]interface{} {
	o := orm.NewOrm()
	num, err := o.QueryTable("route").Filter("id", id).Update(orm.Params{
		"name": name, "route": route,
	})
	if err == nil && num > 0 {
		paginatorMap := make(map[string]interface{})
		paginatorMap["data"] = num
		return paginatorMap
	}
	return nil
}

/*
删除路由
*/
func DeleteRoute(id int) map[string]interface{} {
	o := orm.NewOrm()
	v := Route{Id: id}
	if err := o.Read(&v); err == nil {
		if num, err := o.Delete(&Route{Id: id}); err == nil {
			paginatorMap := make(map[string]interface{})
			paginatorMap["data"] = num
			return paginatorMap
		}
	}
	return nil
}
