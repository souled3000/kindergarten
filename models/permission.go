package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/astaxie/beego/orm"
)

type Permission struct {
	Id             int       `json:"id" orm:"column(id);auto"`
	Name           string    `json:"name" orm:"column(name);size(15)" description:"名称"`
	Identification string    `json:"identification" orm:"column(identification);size(50)" description:"标识"`
	ParentId       int       `json:"parent_id" orm:"column(parent_id)" description:"上一级ID"`
	Level          int8      `json:"level" orm:"column(level)" description:"等级"`
	CreatedAt      time.Time `json:"created_at" orm:"column(created_at);type(datetime)" description:"创建时间"`
	UpdatedAt      time.Time `json:"updated_at" orm:"column(updated_at);type(datetime)" description:"修改时间"`
}

type PermissionTree struct {
	Id             int              `json:"id" orm:"column(id);auto"`
	Name           string           `json:"name" orm:"column(name);size(15)" description:"名称"`
	Identification string           `json:"identification" orm:"column(identification);size(50)" description:"标识"`
	ParentId       int              `json:"parent_id" orm:"column(parent_id)" description:"上一级ID"`
	Level          int8             `json:"level" orm:"column(level)" description:"等级"`
	CreatedAt      time.Time        `json:"created_at" orm:"column(created_at);type(datetime)" description:"创建时间"`
	UpdatedAt      time.Time        `json:"updated_at" orm:"column(updated_at);type(datetime)" description:"修改时间"`
	Child          []PermissionTree `json:"child"`
}

func (t *Permission) TableName() string {
	return "permission"
}

func init() {
	orm.RegisterModel(new(Permission))
}

//保存权限
func AddPermission(name string, identification string, parent_id int, route string) (paginatorMap map[string]interface{}, err error) {
	o := orm.NewOrm()
	var v []orm.Params
	var ident []orm.Params
	var r map[string]interface{}
	json.Unmarshal([]byte(route), &r)
	paginatorMap = make(map[string]interface{})

	qb, _ := orm.NewQueryBuilder("mysql")
	sql := qb.Select("p.*").From("permission as p").Where("p.identification = ?").String()
	num, err := o.Raw(sql, identification).Values(&ident)
	if num > 0 {
		err = errors.New("标识已存在")
		return nil, err
	}
	if parent_id > -1 {
		qb, _ := orm.NewQueryBuilder("mysql")
		sql := qb.Select("p.*").From("permission as p").Where("p.id = ?").String()
		num, err := o.Raw(sql, parent_id).Values(&v)
		if num < 0 {
			err = errors.New("上一级不存在")
			return nil, err
		}
		le := v[0]["level"].(string)
		leve, _ := strconv.Atoi(le)
		lev := leve + 1
		if leve >= 2 {
			err = errors.New("最多两级")
			return nil, err
		}
		qb, _ = orm.NewQueryBuilder("mysql")
		sql = "insert into permission set parent_id = ?,name = ?,level = ?,identification = ?"
		res, err := o.Raw(sql, v[0]["id"], name, lev, identification).Exec()
		id, _ := res.LastInsertId()
		for _, v := range r {
			qb, _ = orm.NewQueryBuilder("mysql")
			sql = "insert into permission_route set permission_id = ?,route_id = ?"
			_, _ = o.Raw(sql, id, v).Exec()
		}
		if err == nil {
			return nil, nil
		}
	} else {
		qb, _ = orm.NewQueryBuilder("mysql")
		sql = "insert into permission set name = ?,identification = ?"
		res, err := o.Raw(sql, name, identification).Exec()
		id, _ := res.LastInsertId()
		for _, v := range r {
			qb, _ = orm.NewQueryBuilder("mysql")
			sql = "insert into permission_route set permission_id = ?,route_id = ?"
			_, _ = o.Raw(sql, id, v).Exec()
		}
		if err == nil {
			return nil, nil
		}
	}
	err = errors.New("保存失败")
	return nil, err
}

//权限详情
func GetPermissionById(id int) map[string]interface{} {
	o := orm.NewOrm()
	var v []orm.Params
	var permission_route []orm.Params
	var parent []orm.Params
	qb, _ := orm.NewQueryBuilder("mysql")
	sql := qb.Select("p.*").From("permission as p").Where("p.id = ?").String()
	_, err := o.Raw(sql, id).Values(&v)

	qb, _ = orm.NewQueryBuilder("mysql")
	sql = qb.Select("pr.route_id").From("permission_route as pr").Where("pr.permission_id = ?").String()
	_, err = o.Raw(sql, id).Values(&permission_route)

	qb, _ = orm.NewQueryBuilder("mysql")
	sql = qb.Select("p.id", "p.name").From("permission as p").Where("p.id = ?").String()
	_, err = o.Raw(sql, v[0]["parent_id"]).Values(&parent)
	if err == nil {
		paginatorMap := make(map[string]interface{})
		paginatorMap["data"] = v
		paginatorMap["route"] = permission_route
		paginatorMap["parent"] = parent
		return paginatorMap
	}
	return nil
}

//权限列表
func GetAllPermission(page int, prepage int) map[string]interface{} {
	o := orm.NewOrm()
	qb, _ := orm.NewQueryBuilder("mysql")
	// 构建查询对象
	sql := qb.Select("count(*)").From("permission").String()
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
		sql := qb.Select("*").From("permission").Limit(prepage).Offset(limit).String()
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

//权限选项
func PermissionOption() map[string]interface{} {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Permission))
	var posts []Permission
	var Permission []PermissionTree
	if _, err := qs.All(&posts); err == nil {
		fmt.Println(posts)
		ml := make(map[string]interface{})
		for _, val := range posts {
			if val.ParentId == 0 {
				next := getNexts(posts, val.Id)
				var tree PermissionTree
				tree.Id = val.Id
				tree.Identification = val.Identification
				tree.Level = val.Level
				tree.Name = val.Name
				tree.ParentId = val.ParentId
				tree.CreatedAt = val.CreatedAt
				tree.UpdatedAt = val.UpdatedAt
				tree.Child = next
				Permission = append(Permission, tree)
			}
		}
		if err == nil {
			ml["data"] = Permission
			return ml
		}
	}
	return nil
}

//获取下级
func getNexts(posts []Permission, id int) (Permission []PermissionTree) {
	for _, val := range posts {
		if val.ParentId == id {
			next := getNexts(posts, val.Id)
			var tree PermissionTree
			tree.Id = val.Id
			tree.Identification = val.Identification
			tree.Level = val.Level
			tree.Name = val.Name
			tree.ParentId = val.ParentId
			tree.CreatedAt = val.CreatedAt
			tree.UpdatedAt = val.UpdatedAt
			tree.Child = next
			Permission = append(Permission, tree)
		}
	}
	return Permission
}
