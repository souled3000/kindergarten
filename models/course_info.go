package models

import (
	"time"

	"github.com/astaxie/beego/orm"
	"strconv"
)

type CourseInfo struct {
	Id           int       `json:"id" orm:"column(id);auto;"`
	Name         string    `json:"name" orm:"column(name);size(30)"; description:"标题"`
	CourseId     int       `json:"course_id" orm:"column(course_id)"`
	TearcherId   int       `json:"tearcher_id" orm:"column(tearcher_id)"`
	TearcherName string    `json:"tearcher_name" orm:"column(tearcher_name);size(30)"`
	Domain       string    `json:"domain" orm:"column(domain);size(30)"`
	Intro        string    `json:"intro" orm:"column(intro);size(30)"`
	Url          string    `json:"url" orm:"column(url)"`
	CoursewareId string    `json:"courseware_id" orm:"column(courseware_id)`
	Aim          string    `json:"aim" orm:"column(aim)`
	Plan         string    `json:"plan" orm:"column(plan)`
	Activity     string    `json:"activity" orm:"column(activity)`
	Job          string    `json:"job" orm:"column(job)`
	Etc          string    `json:"etc" orm:"column(etc)`
	List         string    `json:"list" orm:"column(list)`
	Type         int       `json:"type" orm:"column(type)"`
	Times        string    `json:"times" orm:"column(times)`
	CreatedAt    time.Time `json:"created_at" orm:"auto_now_add"`
}

func (t *CourseInfo) TableName() string {
	return "course_info"
}

func init() {
	orm.RegisterModel(new(CourseInfo))
}

/*
添加
*/
func AddCourseInfo(m *CourseInfo) (map[string]interface{}, error) {
	o := orm.NewOrm()
	id, err := o.Insert(m)
	if err == nil {
		paginatorMap := make(map[string]interface{})
		paginatorMap["data"] = id //返回数据
		return paginatorMap, err
	}
	return nil, err
}

/*
列表
*/
func GetCourseInfoList(class_type int, kindergarten_id int, date string) (map[string]interface{}, error) {
	var v []CourseInfo
	o := orm.NewOrm()
	sql := "select a.* from course_info a left join course b on a.course_id = b.id where b.begin_date <='" + date + "' and b.end_date >= '" + date + "' and b.class_type=" + strconv.Itoa(class_type) + " and b.kindergarten_id =" + strconv.Itoa(kindergarten_id)

	_, err := o.Raw(sql).QueryRows(&v)
	if err == nil {
		paginatorMap := make(map[string]interface{})
		paginatorMap["data"] = v //分页数据
		return paginatorMap, nil

	}
	return nil, err

}

/*
Web -详情
*/
func GetCourseInfoInfo(id int) map[string]interface{} {
	var v []CourseInfo
	o := orm.NewOrm()
	err := o.QueryTable("course_info").Filter("Id", id).One(&v)
	if err == nil {
		paginatorMap := make(map[string]interface{})
		paginatorMap["data"] = v
		return paginatorMap
	}
	return nil
}

/*
删除
*/
func DeleteCourseInfo(id int) map[string]interface{} {
	o := orm.NewOrm()
	v := CourseInfo{Id: id}
	if err := o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&CourseInfo{Id: id}); err == nil {
			paginatorMap := make(map[string]interface{})
			paginatorMap["data"] = num //返回数据
			return paginatorMap
		}
	}
	return nil
}
