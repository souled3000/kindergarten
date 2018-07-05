package models

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/astaxie/beego/orm"
)

type UserPermission struct {
	Id           int `orm:"column(id);auto"`
	UserId       int `orm:"column(user_id)"`
	PermissionId int `orm:"column(permission_id)"`
}

func (t *UserPermission) TableName() string {
	return "user_permission"
}

func init() {
	orm.RegisterModel(new(UserPermission))
}

/*
设置权限
*/
func AddUserPermission(user_id int, role string, permission string, group string) (paginatorMap map[string]interface{}, err error) {
	o := orm.NewOrm()
	o.Begin()
	var r map[string]interface{}
	json.Unmarshal([]byte(role), &r)
	var p map[string]interface{}
	json.Unmarshal([]byte(permission), &p)
	var g map[string]interface{}
	json.Unmarshal([]byte(group), &g)
	paginatorMap = make(map[string]interface{})
	//角色权限
	if r != nil {
		for _, v := range r {
			sql := "insert into user_role set user_id = ?,role_id = ?"
			_, err = o.Raw(sql, user_id, v).Exec()
		}
	}
	//分配用户权限
	if p != nil {
		for _, v := range p {
			sql := "insert into user_permission set user_id = ?,permission_id = ?"
			_, err = o.Raw(sql, user_id, v).Exec()
		}
	}
	//圈子权限
	if g != nil {
		timeLayout := "2006-01-02 15:04:05" //转化所需模板
		loc, _ := time.LoadLocation("")
		timenow := time.Now().Format("2006-01-02 15:04:05")
		created_at, _ := time.ParseInLocation(timeLayout, timenow, loc)
		for k, v := range g {
			sql := "insert into group_view set user_id = ?,class_type = ?,class_id = ?,created_at = ?"
			_, err = o.Raw(sql, user_id, k, v, created_at).Exec()
		}
	}
	if err == nil {
		err = o.Commit()
		return nil, nil
	} else {
		err = o.Rollback()
		err = errors.New("保存失败")
		return nil, err
	}
}

/*
查看用户权限
*/
func GetUserPermissionById(user_id int, kindergarten_id int) (paginatorMap map[string]interface{}, err error) {
	o := orm.NewOrm()
	var r []orm.Params
	var p []orm.Params
	var g []orm.Params
	var rol []orm.Params
	paginatorMap = make(map[string]interface{})
	qb, _ := orm.NewQueryBuilder("mysql")
	sql := qb.Select("r.role_id").From("user_role as r").Where("user_id = ?").String()
	_, err = o.Raw(sql, user_id).Values(&r)

	qb, _ = orm.NewQueryBuilder("mysql")
	sql = qb.Select("p.identification").From("user_permission as up").LeftJoin("permission as p").
		On("up.permission_id = p.id").Where("up.user_id = ?").String()
	_, err = o.Raw(sql, user_id).Values(&p)

	qb, _ = orm.NewQueryBuilder("mysql")
	sql = qb.Select("g.class_id").From("group_view as g").Where("user_id = ?").String()
	_, err = o.Raw(sql, user_id).Values(&g)

	qb, _ = orm.NewQueryBuilder("mysql")
	sql = qb.Select("role.id", "role.name").From("role").String()
	_, err = o.Raw(sql).Values(&rol)

	if err == nil {
		paginatorMap["user_role"] = r
		paginatorMap["user_permission"] = p
		paginatorMap["group_view"] = g
		paginatorMap["group_permission"] = GetGroupPermission(kindergarten_id)
		paginatorMap["roles"] = rol
		paginatorMap["permissions"] = PermissionOption()
		return paginatorMap, nil
	}
	return nil, err
}

/*
查看用户权限标识
*/
func GetUserIdentificationById(user_id int) (paginatorMap map[string]interface{}, err error) {
	o := orm.NewOrm()
	var p []orm.Params
	paginatorMap = make(map[string]interface{})
	qb, _ := orm.NewQueryBuilder("mysql")
	sql := qb.Select("p.identification").From("user_permission as up").LeftJoin("permission as p").
		On("up.permission_id = p.id").Where("up.user_id = ?").String()
	num, err := o.Raw(sql, user_id).Values(&p)
	if err == nil && num > 0 {
		paginatorMap["data"] = p
		return paginatorMap, nil
	}
	err = errors.New("获取失败")
	return nil, err
}

/*
查看圈子权限
*/
func GetGroupIdentificationById(user_id int) (paginatorMap map[string]interface{}, err error) {
	o := orm.NewOrm()
	var g []orm.Params
	paginatorMap = make(map[string]interface{})
	qb, _ := orm.NewQueryBuilder("mysql")
	sql := qb.Select("class_id").From("group_view as gv").Where("gv.user_id = ?").String()
	num, err := o.Raw(sql, user_id).Values(&g)
	if err == nil && num > 0 {
		paginatorMap["data"] = g
		return paginatorMap, nil
	}
	err = errors.New("获取失败")
	return nil, err
}

/*
修改权限
*/
func UpdateUserPermissionById(user_id int, role string, permission string, group string) (paginatorMap map[string]interface{}, err error) {
	o := orm.NewOrm()
	err = o.Begin()
	var r map[string]interface{}
	json.Unmarshal([]byte(role), &r)
	var p map[string]interface{}
	json.Unmarshal([]byte(permission), &p)
	var g map[string]interface{}
	json.Unmarshal([]byte(group), &g)
	paginatorMap = make(map[string]interface{})
	//角色权限
	if r != nil {
		_, err = o.QueryTable("user_role").Filter("user_id", user_id).Delete()

		for _, v := range r {
			sql := "insert into user_role set user_id = ?,role_id = ?"
			_, err = o.Raw(sql, user_id, v).Exec()
		}
	}
	//分配用户权限
	if p != nil {
		_, err = o.QueryTable("user_permission").Filter("user_id", user_id).Delete()
		for _, v := range p {
			sql := "insert into user_permission set user_id = ?,permission_id = ?"
			_, err = o.Raw(sql, user_id, v).Exec()
		}
	}
	//圈子权限
	if g != nil {
		_, err = o.QueryTable("group_view").Filter("user_id", user_id).Delete()
		timeLayout := "2006-01-02 15:04:05" //转化所需模板
		loc, _ := time.LoadLocation("")
		timenow := time.Now().Format("2006-01-02 15:04:05")
		created_at, _ := time.ParseInLocation(timeLayout, timenow, loc)
		for k, v := range g {
			sql := "insert into group_view set user_id = ?,class_type = ?, class_id = ?,created_at = ?"
			_, err = o.Raw(sql, user_id, k, v, created_at).Exec()
		}
	}
	if err == nil {
		o.Commit()
		return nil, nil
	} else {
		o.Rollback()
		err = errors.New("更新失败")
		return nil, err
	}
}

/*
查看用户权限
*/
func GetPermissionRoute(user_id int) ([]orm.Params, error) {
	o := orm.NewOrm()
	var r []orm.Params
	qb, _ := orm.NewQueryBuilder("mysql")
	sql := qb.Select("r.route").From("user_permission as up").LeftJoin("permission as p").
		On("up.permission_id = p.id").LeftJoin("permission_route as pr").
		On("p.id = pr.permission_id").LeftJoin("route as r").
		On("r.id = pr.route_id").Where("up.user_id = ?").String()
	_, err := o.Raw(sql, user_id).Values(&r)
	if err == nil {
		return r, nil
	}
	return nil, err
}

/*
筛选圈子
*/
func GetGroupAll(user_id int, class_type int, role int, kindergarten_id int) (paginatorMap map[string]interface{}, err error) {
	o := orm.NewOrm()
	var g []orm.Params
	paginatorMap = make(map[string]interface{})
	if role == 1 {
		qb, _ := orm.NewQueryBuilder("mysql")
		sql := qb.Select("name as class_name", "id as class_id", "class_type").From("organizational").Where("kindergarten_id = ?").And("class_type = ?").And("level = 3").String()
		_, err := o.Raw(sql, kindergarten_id, class_type).Values(&g)
		if err != nil {
			err = errors.New("未创建班级")
			return nil, err
		}
		for key, val := range g {
			if val["class_type"].(string) == "3" {
				g[key]["class"] = "大班" + val["class_name"].(string) + ""
			} else if val["class_type"].(string) == "2" {
				g[key]["class"] = "中班" + val["class_name"].(string) + ""
			} else {
				g[key]["class"] = "小班" + val["class_name"].(string) + ""
			}
		}
	} else if role == 5 {
		qb, _ := orm.NewQueryBuilder("mysql")
		sql := qb.Select("o.name as class_name", "gv.class_id", "gv.class_type").From("group_view as gv").LeftJoin("organizational as o").
			On("o.id = gv.class_id").Where("gv.user_id = ?").And("gv.class_type = ?").String()
		_, err := o.Raw(sql, user_id, class_type).Values(&g)
		if err != nil {
			err = errors.New("未创建班级")
			return nil, err
		}
		for key, val := range g {
			if val["class_type"].(string) == "3" {
				g[key]["class"] = "大班" + val["class_name"].(string) + ""
			} else if val["class_type"].(string) == "2" {
				g[key]["class"] = "中班" + val["class_name"].(string) + ""
			} else {
				g[key]["class"] = "小班" + val["class_name"].(string) + ""
			}
		}
	}
	if err == nil {
		paginatorMap["data"] = g
		return paginatorMap, nil
	}
	err = errors.New("未创建班级")
	return nil, err
}
