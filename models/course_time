package models

import (
	"math"
	"time"

	"github.com/astaxie/beego/orm"
)

type CourseTime struct {
	Id              int       `json:"id" orm:"column(id);auto;"`
	CourseId              int       `json:"course_id" orm:"column(course_id)"`
	Date       string    `json:"date" orm:"column(date);size(30)"`
	Type 			int			`json:"type" orm:"type"`
	CreatedAt      time.Time `json:"created_at" orm:"auto_now_add"`
	KindergartenTimeId  int	`json:"kindergarten_time_id" orm:"kindergarten_time_id"`
}

func (t *CourseTime) TableName() string {
	return "course_time"
}

func init() {
	orm.RegisterModel(new(CourseTime))
}

//添加
func AddCourseTime(m []CourseTime) (err error) {
	o := orm.NewOrm()
	if _,err := o.InsertMulti(len(m),&m); err == nil {
		return err
	}
	return err
}

/*
列表
*/
func GetCourseTimeList(parent_id int,kindergarten_id int,status int,page,per_page int) (map[string]interface{},error) {
	var v []CourseTime
	o := orm.NewOrm()
	nums, err := o.QueryTable("course").Filter("status",status).Filter("kindergarten_id",kindergarten_id).Filter("parent_id",parent_id).All(&v)
	if err == nil && nums > 0 {
		//根据nums总数，和prepage每页数量 生成分页总数
		totalpages := int(math.Ceil(float64(nums) / float64(per_page))) //page总数
		if page > totalpages {
			page = totalpages
		}
		if page <= 0 {
			page = 1
		}
		limit := (page - 1) * per_page
		num, err := o.QueryTable("course").Limit(per_page, limit).Filter("status",status).Filter("kindergarten_id",kindergarten_id).Filter("parent_id",parent_id).All(&v)
		if err == nil && num > 0 {
			paginatorMap := make(map[string]interface{})
			paginatorMap["total"] = nums          //总条数
			paginatorMap["data"] = v              //分页数据
			paginatorMap["page_num"] = totalpages //总页数
			return paginatorMap,nil
		}
	}
	return nil,err

}

/*
Web -详情
*/
func GetCourseTimeInfo(id int) map[string]interface{} {
	var v []CourseTime
	o := orm.NewOrm()
	err := o.QueryTable("course").Filter("Id", id).One(&v)
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
func DeleteCourseTime(id int) map[string]interface{} {
	o := orm.NewOrm()
	v := CourseTime{Id: id}
	// ascertain id exists in the database
	if err := o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&CourseTime{Id: id}); err == nil {
			paginatorMap := make(map[string]interface{})
			paginatorMap["data"] = num //返回数据
			return paginatorMap
		}
	}
	return nil
}
