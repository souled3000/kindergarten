package models

import (
	"errors"
	"math"
	"time"

	"github.com/astaxie/beego/orm"
)

type TeachersShow struct {
	Id             int       `orm:"column(id);auto"`
	TeacherId      int       `orm:"column(teacher_id)" description:"教师ID"`
	Introduction   string    `orm:"column(introduction);size(100)" description:"介绍"`
	KindergartenId int       `orm:"column(kindergarten_id)" description:"幼儿园ID"`
	CreatedAt      time.Time `orm:"auto_now_add;type(datetime)"`
	UpdatedAt      time.Time `orm:"auto_now;type(datetime)"`
}

func (t *TeachersShow) TableName() string {
	return "teachers_show"
}

func init() {
	orm.RegisterModel(new(TeachersShow))
}

/*
添加教师展示
*/
func AddTeachersShow(teacher_id int, introduction string, kindergarten_id int) (err error) {
	o := orm.NewOrm()
	show := TeachersShow{TeacherId: teacher_id, Introduction: introduction, KindergartenId: kindergarten_id}
	_, err = o.Insert(&show)
	return err
}

/*
教师展示列表
*/
func TeachersShowAll(page int, prepage int, kindergarten_id int) (ml map[string]interface{}, err error) {
	o := orm.NewOrm()
	qb, _ := orm.NewQueryBuilder("mysql")
	// 构建查询对象
	sql := qb.Select("count(*)").From("teachers_show as tw").LeftJoin("teacher as t").
		On("tw.teacher_id = t.teacher_id").Where("tw.kindergarten_id = ?").String()
	var total int64
	err = o.Raw(sql, kindergarten_id).QueryRow(&total)
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
		sql := qb.Select("*").From("teachers_show as tw").LeftJoin("teacher as t").
			On("tw.teacher_id = t.teacher_id").Where("tw.kindergarten_id = ?").Limit(prepage).Offset(limit).String()
		num, err := o.Raw(sql, kindergarten_id).Values(&v)
		if err == nil && num > 0 {
			paginatorMap := make(map[string]interface{})
			paginatorMap["total"] = total         //总条数
			paginatorMap["data"] = v              //分页数据
			paginatorMap["page_num"] = totalpages //总页数
			return paginatorMap, nil
		}
	}
	err = errors.New("没有教师介绍")
	return nil, err
}

/*
教师展示详情
*/
func TeachersShowOne(id int) (ml map[string]interface{}, err error) {
	o := orm.NewOrm()
	var v []orm.Params
	qb, _ := orm.NewQueryBuilder("mysql")
	sql := qb.Select("*").From("teachers_show as tw").LeftJoin("teacher as t").
		On("tw.teacher_id = t.teacher_id").Where("tw.id = ?").String()
	_, err = o.Raw(sql, id).Values(&v)
	ml = make(map[string]interface{})
	ml["data"] = v
	return ml, err
}

/*
删除教师展示
*/
func DeleteTeachersShow(id int) (err error) {
	o := orm.NewOrm()
	_, err = o.QueryTable("teachers_show").Filter("id", id).Delete()
	return err
}

/*
编辑教师展示
*/
func UpdateTeachersShow(id int, teacher_id int, introduction string, kindergarten_id int) (err error) {
	o := orm.NewOrm()
	_, err = o.QueryTable("teachers_show").Filter("id", id).Update(orm.Params{
		"teacher_id": teacher_id, "introduction": introduction, "kindergarten_id": kindergarten_id,
	})
	return err
}
