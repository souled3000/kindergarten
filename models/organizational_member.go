package models

import (
	"errors"
	"fmt"
	"strings"

	"github.com/astaxie/beego/orm"
)

type OrganizationalMember struct {
	Id               int  `json:"id" orm:"column(id);auto"`
	OrganizationalId int  `json:"organizational_id" orm:"column(organizational_id)"`
	MemberId         int  `json:"member_id" orm:"column(member_id)"`
	IsPrincipal      int8 `json:"is_principal" orm:"column(is_principal)" description:"是不是负责人：0不是，1是"`
	Type             int  `json:"type" orm:"column(type)" description:"类型：0教师，1学生"`
}

func (t *OrganizationalMember) TableName() string {
	return "organizational_member"
}

func init() {
	orm.RegisterModel(new(OrganizationalMember))
}

/*
添加成员
*/
func AddMembers(ty int, member_ids string, organizational_id int, is_principal int) (paginatorMap map[string]interface{}, err error) {
	o := orm.NewOrm()
	o.Begin()
	var v []orm.Params
	paginatorMap = make(map[string]interface{})
	s := strings.Split(member_ids, ",")
	fmt.Println(s)
	qb, _ := orm.NewQueryBuilder("mysql")
	sql := qb.Select("o.*").From("organizational as o").Where("id = ?").String()
	_, err = o.Raw(sql, organizational_id).Values(&v)
	if v == nil {
		err = errors.New("没有该班级")
		return nil, err
	}
	//组织架构为园长不能添加
	if v[0]["type"] == "1" && v[0]["is_fixed"] == "1" {
		err = errors.New("不能添加")
		return nil, err
	} else {
		for _, value := range s {
			if value == "" {
				break
			}
			sql := "insert into organizational_member set organizational_id = ?,type = ?,member_id = ?,is_principal = ?"
			_, err = o.Raw(sql, organizational_id, ty, value, is_principal).Exec()
			if err == nil {
				if v[0]["type"] == "2" && v[0]["level"] == "3" {
					//0教师 1学生
					if ty == 0 {
						_, err := o.QueryTable("teacher").Filter("teacher_id", value).Update(orm.Params{
							"status": 1,
						})
						if err != nil {
							o.Rollback()
							return nil, err
						}
					} else {
						_, err := o.QueryTable("student").Filter("student_id", value).Update(orm.Params{
							"status": 1,
						})
						if err != nil {
							o.Rollback()
							return nil, err
						}
					}
				}
			}
		}
		if err == nil {
			o.Commit()
			return nil, err
		}
	}
	err = errors.New("保存失败")
	return nil, err
}

/*
组织架构成员-admin
*/
func GetMembers(organizational_id int) (paginatorMap map[string]interface{}, err error) {
	paginatorMap = make(map[string]interface{})
	o := orm.NewOrm()
	var v []orm.Params
	qb, _ := orm.NewQueryBuilder("mysql")
	sql := qb.Select("om.*", "t.name", "t.number", "t.teacher_id", "t.phone").From("organizational_member as om").LeftJoin("teacher as t").
		On("om.member_id = t.teacher_id").Where("om.organizational_id = ?").And("om.type = 0").String()
	_, err = o.Raw(sql, organizational_id).Values(&v)
	if err == nil {
		paginatorMap["data"] = v
		return paginatorMap, nil
	}
	err = errors.New("获取失败")
	return nil, err
}

/*
组织架构成员负责人-web
*/
func GetWebMembers(organizational_id int) (paginatorMap map[string]interface{}, err error) {
	paginatorMap = make(map[string]interface{})
	o := orm.NewOrm()
	var principal []orm.Params
	var noprincipal []orm.Params
	qb, _ := orm.NewQueryBuilder("mysql")
	sql := qb.Select("t.avatar", "t.name", "o.name as title", "t.user_id", "o.is_fixed", "o.level", "o.type", "t.phone", "t.teacher_id").From("organizational_member as om").LeftJoin("teacher as t").
		On("om.member_id = t.teacher_id").LeftJoin("organizational as o").
		On("om.organizational_id = o.id").Where("om.organizational_id = ?").And("om.is_principal = 1").And("om.type = 0").String()
	_, err = o.Raw(sql, organizational_id).Values(&principal)

	qb, _ = orm.NewQueryBuilder("mysql")
	sql = qb.Select("t.avatar", "t.name", "o.name as title", "t.user_id", "o.is_fixed", "o.level", "o.type", "t.phone", "t.teacher_id").From("organizational_member as om").LeftJoin("teacher as t").
		On("om.member_id = t.teacher_id").LeftJoin("organizational as o").
		On("om.organizational_id = o.id").Where("om.organizational_id = ?").And("om.is_principal = 0").And("om.type = 0").String()
	_, err = o.Raw(sql, organizational_id).Values(&noprincipal)
	if err == nil {
		paginatorMap["principal"] = principal
		paginatorMap["noprincipal"] = noprincipal
		return paginatorMap, nil
	}
	err = errors.New("获取失败")
	return nil, err
}

/*
我的幼儿园教师-web
*/
func MyKinderTeacher(organizational_id int) (paginatorMap map[string]interface{}, err error) {
	paginatorMap = make(map[string]interface{})
	o := orm.NewOrm()
	var mk []orm.Params

	qb, _ := orm.NewQueryBuilder("mysql")
	sql := qb.Select("t.avatar", "t.name", "o.name as title", "t.user_id", "o.is_fixed", "o.level", "o.type").From("organizational_member as om").LeftJoin("teacher as t").
		On("om.member_id = t.teacher_id").LeftJoin("organizational as o").
		On("om.organizational_id = o.id").Where("o.parent_id = ?").And("om.type = 0").String()
	num, err := o.Raw(sql, organizational_id).Values(&mk)
	if err == nil && num > 0 {
		paginatorMap["data"] = mk
		return paginatorMap, nil
	}
	err = errors.New("暂无教师信息")
	return nil, err
}

/*
我的幼儿园列表-web
*/
func MyKinder(kindergarten_id int) (paginatorMap map[string]interface{}, err error) {
	paginatorMap = make(map[string]interface{})
	o := orm.NewOrm()
	var class []orm.Params
	var manage []orm.Params

	qb, _ := orm.NewQueryBuilder("mysql")
	sql := qb.Select("o.*").From("organizational as o").Where("o.kindergarten_id = ?").And("o.type = 1").And("o.is_fixed = 1").And("o.level = 1").String()
	_, err = o.Raw(sql, kindergarten_id).Values(&manage)

	qb, _ = orm.NewQueryBuilder("mysql")
	sql = qb.Select("o.*").From("organizational as o").Where("o.kindergarten_id = ?").And("o.type = 2").And("o.is_fixed = 0").And("o.level = 2").String()
	num, err := o.Raw(sql, kindergarten_id).Values(&class)
	for _, v := range class {
		manage = append(manage, v)
	}
	if err == nil && num > 0 {
		paginatorMap["class"] = manage
		return paginatorMap, nil
	}
	err = errors.New("暂无幼儿园信息")
	return nil, err
}

/*
组织框架移除教师
*/
func DestroyMember(teacher_id int, class_id int) error {
	o := orm.NewOrm()
	var or []orm.Params
	var om []orm.Params
	qb, _ := orm.NewQueryBuilder("mysql")
	sql := qb.Select("om.*").From("organizational_member as om").Where("om.organizational_id = ?").And("om.member_id = ?").String()
	_, err := o.Raw(sql, class_id, teacher_id).Values(&om)
	if om == nil {
		err = errors.New("成员不存在")
		return err
	}
	qb, _ = orm.NewQueryBuilder("mysql")
	sql = qb.Select("o.*").From("organizational as o").Where("o.id = ?").String()
	_, err = o.Raw(sql, class_id).Values(&or)
	if or == nil {
		err = errors.New("班级不存在")
		return err
	}
	//组织架构为园长不能删除
	if or[0]["type"].(string) == "1" && or[0]["is_fixed"] == "1" {
		err = errors.New("不能删除")
		return err
	} else {
		//组织架构类型是年级组并且是第三级需要设置教师或者学生状态为未分配班级
		if or[0]["type"].(string) == "2" && or[0]["level"] == "3" {
			if om[0]["type"] == 0 {
				o.QueryTable("teacher").Filter("teacher_id", teacher_id).Update(orm.Params{
					"status": 0,
				})
			} else {
				o.QueryTable("student").Filter("student_id", teacher_id).Update(orm.Params{
					"status": 0,
				})
			}
			_, err = o.QueryTable("organizational_member").Filter("organizational_id", class_id).Filter("member_id", teacher_id).Delete()
		} else {
			_, err = o.QueryTable("organizational_member").Filter("organizational_id", class_id).Filter("member_id", teacher_id).Delete()
		}
	}
	if err == nil {
		return nil
	}
	err = errors.New("移除失败")
	return err
}
