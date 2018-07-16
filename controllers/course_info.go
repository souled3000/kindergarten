package controllers

import (
	"kindergarten-service-go/models"
	"github.com/astaxie/beego/validation"
)

type CourseInfoController struct {
	BaseController
}

func (c *CourseInfoController) URLMapping() {
	c.Mapping("Get", c.Get)
}

// GetInfoList ...
// @Title 课程列表
// @Description 课程列表
// @Param	kindergarten_id	query	int	true	"幼儿园id"
// @Param	date	query	string	true	"日期"
// @Param	class_type	query	int	true	"班级类型"
// @Success 0 			{json} 	JSONStruct
// @Failure 1005 获取失败
// @router / [get]
func (c *CourseInfoController) GetAll() {
	class_type,_ := c.GetInt("class_type")
	kindergarten_id,_ := c.GetInt("kindergarten_id")
	date := c.GetString("date")

	if list, err := models.GetCourseInfoList(class_type,kindergarten_id,date); err == nil {
		c.Data["json"] = JSONStruct{"success", 0, list, "获取成功"}
	} else {
		c.Data["json"] = JSONStruct{"error", 1005, nil, "获取失败"}
	}

	c.ServeJSON()
}

// @Title 添加园本课程，专题，课程
// @Description 添加园本课程，专题，课程
// @param 		parent_id				query  	int    	true		"上级id"
// @param 		kindergarten_id		query  	int    	true		"幼儿园ID"
// @param 		leval				query  	int    	true		"园本课程1，专题2，课程3"
// @param 		status			query  	int    	true		""
// @param 		name				query  	string    	true		"名称"
// @param		aim			query	string	true		"目标"
// @param 		begin_date			query  	string 	true		"开始时间"
// @param 		end_date			query  	string 	true		"结束时间"
// @router /addinfo [post]
func (c *CourseInfoController) Add_info() {
	kindergarten_id, _ := c.GetInt("kindergarten_id")
	course_id, _ := c.GetInt("course_id")
	tearcher_id, _ := c.GetInt("tearcher_id")
	name := c.GetString("name")
	aim := c.GetString("aim")
	tearcher_name := c.GetString("tearcher_name")
	domain := c.GetString("domain")
	intro := c.GetString("intro")
	url := c.GetString("url")
	courseware_id := c.GetString("courseware_id")
	plan := c.GetString("plan")
	activity := c.GetString("activity")
	etc := c.GetString("etc")
	list := c.GetString("list")
	job := c.GetString("job")
	var course models.CourseInfo
	course.CourseId = course_id
	course.TearcherId = tearcher_id
	course.TearcherName = tearcher_name
	course.Name = name
	course.Aim = aim
	course.Domain = domain
	course.Intro = intro
	course.Url = url
	course.CoursewareId = courseware_id
	course.Plan = plan
	course.Activity = activity
	course.Etc = etc
	course.List = list
	course.Job = job
	valid := validation.Validation{}
	valid.Required(kindergarten_id, "kindergarten_id").Message("幼儿园ID不能为空")
	if valid.HasErrors() {
		c.Data["json"] = JSONStruct{"error", 1001, nil, valid.Errors[0].Message}
	} else {
		if _, err := models.AddCourseInfo(&course); err == nil {
			c.Data["json"] = JSONStruct{"success", 0, err, "保存成功"}
		} else {
			c.Data["json"] = JSONStruct{"error", 1003, err.Error(), "保存失败"}
		}
	}
	c.ServeJSON()
}
