package models

import (
	"math"
	"time"

	"github.com/astaxie/beego/orm"
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
func AddCourseInfo(m *CourseInfo, class_course_time []CourseTime) (map[string]interface{}, error) {
	o := orm.NewOrm()
	id, err := o.Insert(m)
	if err == nil {
		if len(class_course_time) > 0 {
			for key, _ := range class_course_time {
				class_course_time[key].CourseId = int(id)
			}
			AddCourseTime(class_course_time)
		}
		paginatorMap := make(map[string]interface{})
		paginatorMap["data"] = id //返回数据
		return paginatorMap, err
	}
	return nil, err
}

/*
列表
*/
func GetCourseInfoList(parent_id int, kindergarten_id int, status int, page, per_page int) (map[string]interface{}, error) {
	var v []CourseInfo
	o := orm.NewOrm()
	nums, err := o.QueryTable("course").Filter("status", status).Filter("kindergarten_id", kindergarten_id).Filter("parent_id", parent_id).All(&v)
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
		num, err := o.QueryTable("course").Limit(per_page, limit).Filter("status", status).Filter("kindergarten_id", kindergarten_id).Filter("parent_id", parent_id).All(&v)
		if err == nil && num > 0 {
			paginatorMap := make(map[string]interface{})
			paginatorMap["total"] = nums          //总条数
			paginatorMap["data"] = v              //分页数据
			paginatorMap["page_num"] = totalpages //总页数
			return paginatorMap, nil
		}
	}
	return nil, err

}

/*
Web -详情
*/
func GetCourseInfoInfo(id int) map[string]interface{} {
	var v []CourseInfo
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
func DeleteCourseInfo(id int) map[string]interface{} {
	o := orm.NewOrm()
	v := CourseInfo{Id: id}
	// ascertain id exists in the database
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
