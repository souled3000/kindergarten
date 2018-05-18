package controllers

import (
	"kindergarten-service-go/models"

	"github.com/astaxie/beego"
)

//学生
type StudentController struct {
	beego.Controller
}

// URLMapping ...
func (c *StudentController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
}

// GetStudent ...
// @Title 学生列表
// @Description 学生列表
// @Param	kindergarten_id       query	int	     true		"幼儿园ID"
// @Param	status                query	int	     false		"状态"
// @Param	search                query	int	     false		"搜索条件"
// @Param	page                  query	int	     false		"页数"
// @Param	per_page              query	int	     false		"每页显示条数"
// @Success 200 {object} models.Student
// @Failure 403
// @router / [get]
func (c *StudentController) GetStudent() {
	var prepage int = 20
	var page int
	var kindergarten_id int
	var status int
	var search string
	search = c.GetString("search")
	if v, err := c.GetInt("per_page"); err == nil {
		prepage = v
	}
	if v, err := c.GetInt("page"); err == nil {
		page = v
	}
	if v, err := c.GetInt("kindergarten_id"); err == nil {
		kindergarten_id = v
	}
	if v, err := c.GetInt("status", -1); err == nil {
		status = v
	}
	v := models.GetStudent(kindergarten_id, status, search, page, prepage)
	if v == nil {
		c.Data["json"] = JSONStruct{"error", 1005, nil, "获取失败"}
	} else {
		c.Data["json"] = JSONStruct{"success", 0, v, "获取成功"}
	}
	c.ServeJSON()
}

// GetStudentClass ...
// @Title 学生班级搜索
// @Description 学生班级搜索
// @Param	kindergarten_id       query	int	     true		"幼儿园ID"
// @Param	class_type            query	int	     true		"班级类型"
// @Param	page                  query	int	     false		"页数"
// @Param	per_page              query	int	     false		"每页显示条数"
// @Success 200 {object} models.Student
// @Failure 403
// @router /class [get]
func (c *TeacherController) GetStudentClass() {
	var prepage int = 20
	var page int
	var kindergarten_id int
	var class_type int
	if v, err := c.GetInt("per_page"); err == nil {
		prepage = v
	}
	if v, err := c.GetInt("page"); err == nil {
		page = v
	}
	if v, err := c.GetInt("kindergarten_id"); err == nil {
		kindergarten_id = v
	}
	if v, err := c.GetInt("class_type"); err == nil {
		class_type = v
	}
	v := models.GetStudentClass(kindergarten_id, class_type, page, prepage)
	if v == nil {
		c.Data["json"] = JSONStruct{"error", 1005, nil, "获取失败"}
	} else {
		c.Data["json"] = JSONStruct{"success", 0, v, "获取成功"}
	}
	c.ServeJSON()
}
