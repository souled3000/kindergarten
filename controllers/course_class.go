package controllers

import (
	"kindergarten-service-go/models"
	"github.com/astaxie/beego/validation"
)

type CourseClassController struct {
	BaseController
}

func (c *CourseClassController) URLMapping() {
	c.Mapping("Get", c.Get)
}

// PostTimeClass ...
// @Title 班级课程安排
// @Description 班级课程安排
// @param 		class_id		query  	int    	true		"班级id"
// @param 		begin_date				query  	string    	true		"开始时间"
// @param 		end_date				query  	string    	true		"结束时间"
// @param 		content				query  	string    	true		"时间安排json"
// @router / [post]
func (c *CourseClassController) Post() {
	content := c.GetString("content")
	class_id,_ := c.GetInt("class_id")
	begin_date := c.GetString("begin_date")
	end_date := c.GetString("end_date")
	var course models.CourseClass
	course.ClassId = class_id
	course.Content = content
	course.BeginDate = begin_date
	course.EndDate = end_date
	valid := validation.Validation{}
	valid.Required(content, "content").Message("参数不能为空")
	if valid.HasErrors() {
		c.Data["json"] = JSONStruct{"error", 1001, nil, valid.Errors[0].Message}
	} else {
		if _,err := models.AddCourseClass(&course); err == nil {
			c.Data["json"] = JSONStruct{"success", 0, err, "保存成功"}
		} else {
			c.Data["json"] = JSONStruct{"error", 1003, nil, "保存失败"}
		}
	}
	c.ServeJSON()
}

// GetTimeOne ...
// @Title 获取班级某一天课程
// @Description 获取班级某一天课程
// @Param	kindergarten_id	query	int	true	"幼儿园id"
// @Param	class_id	query	int	true	"班级id"
// @Success 0 			{json} 	JSONStruct
// @Failure 1005 获取失败
// @router /class_day [get]
func (c *CourseClassController) GetTimeOne() {

	kindergarten_id,_ := c.GetInt("kindergarten_id")
	class_id,_ := c.GetInt("class_id")
	types,_ := c.GetInt("type")
	time := c.GetString("time")

	if list := models.GetCourseClassInfo(class_id,time,types,kindergarten_id); list == nil {
		c.Data["json"] = JSONStruct{"error", 1005, nil, "获取失败"}
	} else {
		c.Data["json"] = JSONStruct{"success", 0, list, "获取成功"}
	}

	c.ServeJSON()
}

// GetTimelist ...
// @Title 获取班级课程表
// @Description 获取班级课程表
// @Param	kindergarten_id	query	int	true	"幼儿园id"
// @Param	class_id	query	int	true	"班级id"
// @Success 0 			{json} 	JSONStruct
// @Failure 1005 获取失败
// @router /class_course [get]
func (c *CourseClassController) GetTimelist() {

	kindergarten_id,_ := c.GetInt("kindergarten_id")
	class_id,_ := c.GetInt("class_id")
	time := c.GetString("time")

	if list := models.GetClassTime(class_id,kindergarten_id,time); list == nil {
		c.Data["json"] = JSONStruct{"error", 1005, nil, "获取失败"}
	} else {
		c.Data["json"] = JSONStruct{"success", 0, list, "获取成功"}
	}

	c.ServeJSON()
}

// GetPlan ...
// @Title 计划列表
// @Description 计划列表
// @Param	kindergarten_id	query	int	true	"幼儿园id"
// @Param	class_id	query	int	true	"班级id"
// @Success 0 			{json} 	JSONStruct
// @Failure 1005 获取失败
// @router /plan [get]
func (c *CourseClassController) GetPlan() {
	class_id,_ := c.GetInt("class_id")
	time := c.GetString("time")

	if list := models.PlanCourseClass(class_id,time); list == nil {
		c.Data["json"] = JSONStruct{"error", 1005, nil, "获取失败"}
	} else {
		c.Data["json"] = JSONStruct{"success", 0, list, "获取成功"}
	}

	c.ServeJSON()
}

// GetPlan ...
// @Title 计划详情
// @Description 计划详情
// @Param	id	query	int	true	"计划列表id"
// @Success 0 			{json} 	JSONStruct
// @Failure 1005 获取失败
// @router /plan_info [get]
func (c *CourseClassController) GetPlanInfo() {
	id,_ := c.GetInt("id")

	if list := models.PlanInfoCourseClass(id); list == nil {
		c.Data["json"] = JSONStruct{"error", 1005, nil, "获取失败"}
	} else {
		c.Data["json"] = JSONStruct{"success", 0, list, "获取成功"}
	}

	c.ServeJSON()
}