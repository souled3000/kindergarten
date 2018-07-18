package controllers

import (
	"github.com/astaxie/beego/validation"
	"kindergarten-service-go/models"
	"strconv"
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
	class_type, _ := c.GetInt("class_type")
	kindergarten_id, _ := c.GetInt("kindergarten_id")
	date := c.GetString("date")

	if list, err := models.GetCourseInfoList(class_type, kindergarten_id, date); err == nil {
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
	types, _ := c.GetInt("type")

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
	times := c.GetString("times")
	job := c.GetString("job")
	var course models.CourseInfo
	course.CourseId = course_id
	course.TearcherId = tearcher_id
	course.TearcherName = tearcher_name
	course.Name = name
	course.Aim = aim
	course.Type = types
	course.Domain = domain
	course.Intro = intro
	course.Url = url
	course.CoursewareId = courseware_id
	course.Plan = plan
	course.Activity = activity
	course.Etc = etc
	course.Times = times
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


// GetOne ...
// @Title 			课程详情
// @Description 	课程详情
// @Param	id		path 	string	true		"课程ID"
// @Success 0 		{string} 	success
// @Failure 1004		获取失败
// @router /:id [get]
func (c *CourseInfoController) GetOne() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	if l := models.GetCourseInfoInfo(id); l != nil {
		c.Data["json"] = JSONStruct{"success", 0, l, "获取成功"}
	} else {
		c.Data["json"] = JSONStruct{"error", 1004, nil, "获取失败"}
	}
	c.ServeJSON()
}


// Delete ...
// @Title 			删除课程
// @Description 		删除课程
// @Param	id		path 	string	true		"课程ID"
// @Success 0 		{string} 	success
// @Failure 1004		删除失败
// @router /:id [delete]
func (c *CourseInfoController) Delete() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	if err := models.DeleteCourseInfo(id); err == nil {
		c.Data["json"] = JSONStruct{"success", 0, err, "删除成功"}
	} else {
		c.Data["json"] = JSONStruct{"error", 1004, nil, "删除失败"}
	}
	c.ServeJSON()
}