package models

import (
	"fmt"
	"math"
	"time"

	"github.com/astaxie/beego/orm"
)

type Organizational struct {
	Id             int       `orm:"column(id);auto"`
	KindergartenId int       `orm:"column(kindergarten_id)" description:"幼儿园id"`
	ParentId       int       `orm:"column(parent_id)" description:"父级id"`
	Name           string    `orm:"column(name);size(20)" description:"组织架构名字"`
	IsFixed        int8      `orm:"column(is_fixed)" description:"是否固定的：0不是，1是"`
	Level          int8      `orm:"column(level)" description:"等级"`
	ParentIds      string    `orm:"column(parent_ids);size(50)" description:"父级所有id"`
	Type           int8      `orm:"column(type)" description:"类型：0普通，1管理层，2年级组"`
	ClassType      int8      `orm:"column(class_type)" description:"班级类型：1小班，2中班，3大班"`
	CreatedAt      time.Time `orm:"column(created_at);type(datetime)" description:"添加时间"`
	UpdatedAt      time.Time `orm:"column(updated_at);type(datetime)" description:"修改时间"`
}

func (t *Organizational) TableName() string {
	return "organizational"
}

func init() {
	orm.RegisterModel(new(Organizational))
}

//班级搜索
func GetClassAll(kindergarten_id int, class_type int, page int, prepage int) map[string]interface{} {
	o := orm.NewOrm()
	qb, _ := orm.NewQueryBuilder("mysql")
	var condition []interface{}
	where := "1=1 "
	if kindergarten_id == 0 {
		where += " AND kindergarten_id = ?"
	} else {
		where += " AND kindergarten_id = ?"
		condition = append(condition, kindergarten_id)
	}
	if class_type > 0 {
		where += " AND class_type = ?"
		condition = append(condition, class_type)
	}
	where += " AND type = ?"
	condition = append(condition, 2)
	where += " AND level = ?"
	condition = append(condition, 3)
	// 构建查询对象
	sql := qb.Select("count(*)").From("organizational").Where(where).String()
	var total int64
	err := o.Raw(sql, condition).QueryRow(&total)
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
		sql := qb.Select("*").From("organizational").Where(where).Limit(prepage).Offset(limit).String()
		num, err := o.Raw(sql, condition).Values(&v)
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

//班级成员
func ClassMember(kindergarten_id int, class_type int, class_id int, page int, prepage int) map[string]interface{} {
	o := orm.NewOrm()
	var student []orm.Params
	var teacher []orm.Params
	var condition []interface{}
	where := "1=1 "
	if kindergarten_id == 0 {
		where += " AND o.kindergarten_id = ?"
	} else {
		where += " AND o.kindergarten_id = ?"
		condition = append(condition, kindergarten_id)
	}
	if class_type > 0 {
		where += " AND o.class_type = ?"
		condition = append(condition, class_type)
	}
	if class_id > 0 {
		where += " AND o.id = ?"
		condition = append(condition, class_id)
	}
	where += " AND o.type = ?"
	condition = append(condition, 2)
	where += " AND o.level = ?"
	condition = append(condition, 3)

	qb, _ := orm.NewQueryBuilder("mysql")
	sql := qb.Select("s.student_id", "s.name", "s.avatar", "s.number", "s.phone",
		"o.name as class_name", "o.class_type", "om.id").From("student as s").LeftJoin("organizational_member as om").
		On("s.student_id = om.member_id").LeftJoin("organizational as o").
		On("om.organizational_id = o.id").Where(where).And("om.type = 1").String()
	_, err := o.Raw(sql, condition).Values(&student)

	qb, _ = orm.NewQueryBuilder("mysql")
	sql = qb.Select("t.teacher_id", "t.name", "t.avatar", "t.number", "t.phone",
		"o.name as class_name", "om.id").From("teacher as t").LeftJoin("organizational_member as om").
		On("t.teacher_id = om.member_id").LeftJoin("organizational as o").
		On("om.organizational_id = o.id").Where(where).And("om.type = 0").String()
	_, err = o.Raw(sql, condition).Values(&teacher)
	if err == nil {
		paginatorMap := make(map[string]interface{})
		paginatorMap["student"] = student
		paginatorMap["teacher"] = teacher
		return paginatorMap
	}
	return nil
}

//删除班级
func Destroy(class_id int) map[string]interface{} {
	o := orm.NewOrm()
	err := o.Begin()
	var t []orm.Params
	var s []orm.Params
	var condition []interface{}
	where := "1=1 "
	if class_id > 0 {
		where += " AND organizational_id = ?"
		condition = append(condition, class_id)
	}
	qb, _ := orm.NewQueryBuilder("mysql")
	sql := qb.Select("member_id").From("organizational_member").Where(where).And("type = 0").String()
	_, err = o.Raw(sql, condition).Values(&t)
	//修改teacher
	for key, _ := range t {
		_, err = o.QueryTable("teacher").Filter("teacher_id", t[key]["member_id"]).Update(orm.Params{
			"status": 0,
		})
		if err != nil {
			err = o.Rollback()
		}
	}
	//修改学生
	qb, _ = orm.NewQueryBuilder("mysql")
	sql = qb.Select("member_id").From("organizational_member").Where(where).And("type = 1").String()
	_, err = o.Raw(sql, condition).Values(&s)

	for key, _ := range s {
		fmt.Println(s[key]["member_id"])
		_, err = o.QueryTable("student").Filter("student_id", s[key]["member_id"]).Update(orm.Params{
			"status": 0,
		})
		if err != nil {
			err = o.Rollback()
		}
	}
	//删除班级
	_, err = o.QueryTable("organizational").Filter("id", class_id).Delete()
	if err != nil {
		err = o.Rollback()
	}
	//删除班级成员
	num, err := o.QueryTable("organizational_member").Filter("organizational_id", class_id).Delete()

	if err != nil {
		err = o.Rollback()
	} else {
		err = o.Commit()
	}

	if err == nil {
		paginatorMap := make(map[string]interface{})
		paginatorMap["data"] = num
		return paginatorMap
	}
	return nil
}
