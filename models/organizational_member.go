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

//添加成员
func AddMembers(ty int, member_ids string, organizational_id int, is_principal int) (paginatorMap map[string]interface{}, err error) {
	paginatorMap = make(map[string]interface{})
	o := orm.NewOrm()
	var v []orm.Params
	err = o.Begin()
	s := strings.Split(member_ids, ",")
	qb, _ := orm.NewQueryBuilder("mysql")
	sql := qb.Select("o.*").From("organizational as o").Where("id = ?").String()
	_, err = o.Raw(sql, organizational_id).Values(&v)
	//组织架构为园长不能添加
	if v[0]["type"] == "1" && v[0]["is_fixed"] == "1" {
		err = errors.New("不能添加")
		return nil, err
	}
	for _, value := range s {
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
						err = o.Rollback()
					}
				} else {
					_, err := o.QueryTable("student").Filter("student_id", value).Update(orm.Params{
						"status": 1,
					})
					if err != nil {
						err = o.Rollback()
					}
				}
			}
		}
	}
	if err == nil {
		err = o.Commit()
	} else {
		err = o.Rollback()
	}
	if err == nil {
		paginatorMap["data"] = nil //返回数据
		return paginatorMap, nil
	}
	err = errors.New("保存失败")
	return nil, err
}

//组织架构成员
func GetMembers(organizational_id int) (paginatorMap map[string]interface{}, err error) {
	paginatorMap = make(map[string]interface{})
	o := orm.NewOrm()
	var v []orm.Params
	qb, _ := orm.NewQueryBuilder("mysql")
	sql := qb.Select("om.*", "t.name", "t.number", "t.teacher_id", "t.phone").From("organizational_member as om").LeftJoin("teacher as t").
		On("om.member_id = t.teacher_id").Where("om.organizational_id = ?").And("om.type = 0").String()
	_, err = o.Raw(sql, organizational_id).Values(&v)
	fmt.Println(v)
	if err == nil {
		paginatorMap["data"] = v //返回数据
		return paginatorMap, nil
	}
	err = errors.New("获取失败")
	return nil, err
}
