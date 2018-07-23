package controllers

import (
	"kindergarten-service-go/models"

	"github.com/astaxie/beego/validation"
	"strconv"
)

type CourseController struct {
	BaseController
}

func (c *CourseController) URLMapping() {
	c.Mapping("Get", c.Get)
}

// GetAll ...
// @Title 园本课程，专题，目标
// @Description 园本课程，专题，目标
// @Param	parent_id	query	int	true	"上级id"
// @Param	kindergarten_id	query	int	true	"幼儿园id"
// @Param	status	query	int	true	"状态"
// @Param	page	query	int	true	"页"
// @Param	per_page	query	int	true	"每页条数"
// @Success 0 			{string} 	success
// @Failure 1005 获取失败
// @router / [get]
func (c *CourseController) GetAll() {
	parent_id, _ := c.GetInt("parent_id")
	kindergarten_id, _ := c.GetInt("kindergarten_id")
	status, _ := c.GetInt("status")
	page, _ := c.GetInt("page")
	per_page, _ := c.GetInt("per_page")

	if list, err := models.GetCourseList(parent_id, kindergarten_id, status, page, per_page); err == nil {
		c.Data["json"] = JSONStruct{"success", 0, list, "获取成功"}
	} else {
		c.Data["json"] = JSONStruct{"error", 1005, nil, "获取失败"}
	}

	c.ServeJSON()
}

// Post ...
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
// @router / [post]
func (c *CourseController) Post() {
	class_type, _ := c.GetInt("class_type")
	kindergarten_id, _ := c.GetInt("kindergarten_id")
	parent_id, _ := c.GetInt("parent_id")
	status, _ := c.GetInt("status")
	name := c.GetString("name")
	aim := c.GetString("aim")
	begin_date := c.GetString("begin_date")
	end_date := c.GetString("end_date")
	var course models.Course
	course.KindergartenId = kindergarten_id
	course.ClassType = class_type
	course.ParentId = parent_id
	course.Status = status
	course.Name = name
	course.Aim = aim
	course.BeginDate = begin_date
	course.EndDate = end_date
	valid := validation.Validation{}
	valid.Required(kindergarten_id, "kindergarten_id").Message("幼儿园ID不能为空")
	if valid.HasErrors() {
		c.Data["json"] = JSONStruct{"error", 1001, nil, valid.Errors[0].Message}
	} else {
		if _, err := models.AddCourse(&course); err == nil {
			c.Data["json"] = JSONStruct{"success", 0, err, "保存成功"}
		} else {
			c.Data["json"] = JSONStruct{"error", 1003, nil, "保存失败"}
		}
	}
	c.ServeJSON()
}

// Add_use ...
// @Title 专题时间安排
// @Description 专题时间安排
// @param 		id				query  	int    	true		"专题id"
// @param 		begin_date			query  	string 	true		"开始时间"
// @param 		end_date			query  	string 	true		"结束时间"
// @param 		name			query  	string 	true		"名称"
// @param 		url			query  	string 	true		"封面"
// @router /adduse [post]
func (c *CourseController) Add_use() {
	id, _ := c.GetInt("id")
	kindergarten_id, _ := c.GetInt("kindergarten_id")

	valid := validation.Validation{}
	valid.Required(id, "id").Message("专题id不能为空")
	if valid.HasErrors() {
		c.Data["json"] = JSONStruct{"error", 1001, nil, valid.Errors[0].Message}
	} else {
		if err := models.UseCourse(id, kindergarten_id); err == nil {
			c.Data["json"] = JSONStruct{"success", 0, err, "保存成功"}
		} else {
			c.Data["json"] = JSONStruct{"error", 1003, nil, "保存失败"}
		}
	}
	c.ServeJSON()
}


// Add_time ...
// @Title 专题时间安排
// @Description 专题时间安排
// @param 		id				query  	int    	true		"专题id"
// @param 		begin_date			query  	string 	true		"开始时间"
// @param 		end_date			query  	string 	true		"结束时间"
// @param 		name			query  	string 	true		"名称"
// @param 		url			query  	string 	true		"封面"
// @router /addtime [post]
func (c *CourseController) Add_time() {
	id, _ := c.GetInt("id")
	begin_date := c.GetString("begin_date")
	end_date := c.GetString("end_date")
	name := c.GetString("name")
	url := c.GetString("url")
	valid := validation.Validation{}
	valid.Required(id, "id").Message("专题id不能为空")
	if valid.HasErrors() {
		c.Data["json"] = JSONStruct{"error", 1001, nil, valid.Errors[0].Message}
	} else {
		if err := models.UpdataCourse(id, begin_date, end_date, name, url); err == nil {
			c.Data["json"] = JSONStruct{"success", 0, err, "保存成功"}
		} else {
			c.Data["json"] = JSONStruct{"error", 1003, nil, "保存失败"}
		}
	}
	c.ServeJSON()
}

// DelCourse ...
// @Title 			删除
// @Description 		删除
// @Param	id		path 	string	true		"特殊儿童ID"
// @Success 0 		{string} 	success
// @Failure 1004		删除失败
// @router /:id [delete]
func (c *CourseController) DelCourse() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	if err := models.DeleteCourse(id); err == nil {
		c.Data["json"] = JSONStruct{"success", 0, err, "删除成功"}
	} else {
		c.Data["json"] = JSONStruct{"error", 1004, nil, "删除失败"}
	}
	c.ServeJSON()
}

// GetCourse ...
// @Title 专题详情
// @Description 专题详情
// @Param	id	query	int	true	"专题id"
// @Success 0 			{string} 	success
// @Failure 1005 获取失败
// @router /courseinfo [get]
func (c *CourseController) GetCourse() {
	id, _ := c.GetInt("id")
	sass, _ := c.GetInt("sass")
	if list, err := models.InfoCourse(id,sass); err == nil {
		c.Data["json"] = JSONStruct{"success", 0, list, "获取成功"}
	} else {
		c.Data["json"] = JSONStruct{"error", 1005, nil, "获取失败"}
	}

	c.ServeJSON()
}

// PostTime ...
// @Title 幼儿园添加时间安排(时间段添加)
// @Description 幼儿园添加时间安排(时间段)
// @param 		kindergarten_id		query  	int    	true		"幼儿园id"
// @param 		class_id		query  	int    	true		"班级id"
// @param 		name		query  	string    	true		"时间安排名"
// @param 		end_time				query  	string    	true		"结束时间"
// @param 		begin_time				query  	string    	true		"开始时间"
// @router /add_time [post]
func (c *CourseController) PostTime() {
	kindergarten_id, _ := c.GetInt("kindergarten_id")
	name := c.GetString("name")
	begin_time := c.GetString("begin_time")
	end_time := c.GetString("end_time")
	types, _ := c.GetInt("type")
	class_type, _ := c.GetInt("class_type")
	class_id, _ := c.GetInt("class_id")
	var course models.KindergartenTime
	course.KindergartenId = kindergarten_id
	course.Name = name
	course.ClassType = class_type
	course.ClassId = class_id
	course.Type = types
	course.BeginTime = begin_time
	course.EndTime = end_time
	valid := validation.Validation{}
	valid.Required(kindergarten_id, "kindergarten_id").Message("幼儿园ID不能为空")
	if valid.HasErrors() {
		c.Data["json"] = JSONStruct{"error", 1001, nil, valid.Errors[0].Message}
	} else {
		if l, err := models.AddKindergartenTime(course); err == nil {
			c.Data["json"] = JSONStruct{"success", 0, l, "保存成功"}
		} else {
			c.Data["json"] = JSONStruct{"error", 1003, err, "保存失败"}
		}
	}
	c.ServeJSON()
}

// GetTimeInfo ...
// @Title 幼儿园时间安排列表
// @Description 幼儿园时间安排列表
// @Param	class_id	query	int	true	"班级id"
// @Success 0 			{string} 	success
// @Param	class_type	query	int	true	"班级类型  1，2，3"
// @Success 0 			{json} 	JSONStruct
// @Failure 1005 获取失败
// @router /time_list [get]
func (c *CourseController) GetTimeInfo() {
	class_type, _ := c.GetInt("class_type")
	class_id, _ := c.GetInt("class_id")
	if list := models.GetKindergartenTimeInfo(class_type, class_id); list == nil {
		c.Data["json"] = JSONStruct{"error", 1005, nil, "获取失败"}
	} else {
		c.Data["json"] = JSONStruct{"success", 0, list, "获取成功"}
	}

	c.ServeJSON()
}
