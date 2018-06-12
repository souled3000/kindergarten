package models

import (
	"math"
	"time"

	"github.com/astaxie/beego/orm"
)

type Role struct {
	Id        int       `json:"id" orm:"column(id);auto"`
	Name      string    `json:"name" orm:"column(name);size(15)" description:"名称"`
	CreatedAt time.Time `json:"created_at" orm:"auto_now_add"`
	UpdatedAt time.Time `json:"updated_at" orm:"auto_now"`
}

func (t *Role) TableName() string {
	return "role"
}

func init() {
	orm.RegisterModel(new(Role))
}

/*
添加角色
*/
func AddRole(name string) map[string]interface{} {
	o := orm.NewOrm()
	var r Role
	r.Name = name
	id, err := o.Insert(&r)
	if err == nil {
		paginatorMap := make(map[string]interface{})
		paginatorMap["data"] = id
		return paginatorMap
	}
	return nil
}

/*
角色详情
*/
func GetRoleById(id int) map[string]interface{} {
	o := orm.NewOrm()
	var v []orm.Params
	qb, _ := orm.NewQueryBuilder("mysql")
	sql := qb.Select("r.id", "r.name").From("role as r").Where("id = ?").String()
	num, err := o.Raw(sql, id).Values(&v)
	if err == nil && num > 0 {
		paginatorMap := make(map[string]interface{})
		paginatorMap["data"] = v
		return paginatorMap
	}
	return nil
}

/*
角色列表
*/
func GetAllRole(page int, prepage int) map[string]interface{} {
	o := orm.NewOrm()
	qb, _ := orm.NewQueryBuilder("mysql")

	// 构建查询对象
	sql := qb.Select("count(*)").From("role as r").String()
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
		sql := qb.Select("r.id", "r.name").From("role as r").Limit(prepage).Offset(limit).String()
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
编辑角色
*/
func UpdateRoleById(id int, name string) map[string]interface{} {
	o := orm.NewOrm()
	v := Role{Id: id}
	if err := o.Read(&v); err == nil {
		v.Name = name
		if num, err := o.Update(&v); err == nil {
			paginatorMap := make(map[string]interface{})
			paginatorMap["data"] = num
			return paginatorMap
		}
	}
	return nil
}
