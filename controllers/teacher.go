package controllers

import (
	"encoding/json"
	"kindergarten-service-go/models"
	"strconv"

	"github.com/astaxie/beego"
)

//教师
type TeacherController struct {
	beego.Controller
}

// URLMapping ...
func (c *TeacherController) URLMapping() {
	c.Mapping("GetTeacherDown", c.GetTeacherDown)
	c.Mapping("GetTeacher", c.GetTeacher)
	c.Mapping("Delete", c.Delete)
}

// GetTeacherDown ...
// @Title 教师下拉菜单
// @Description 教师下拉菜单
// @Param	id		path 	string	true		"幼儿园ID"
// @Success 200 {object} models.Teacher
// @Failure 403 :id is empty
// @router /teacher_down/:id [get]
func (c *TeacherController) GetTeacherDown() {
	var prepage int = 20
	var page int
	if v, err := c.GetInt("per_page"); err == nil {
		prepage = v
	}
	if v, err := c.GetInt("page"); err == nil {
		page = v
	}
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	v := models.GetTeacherById(id, page, prepage)
	if v == nil {
		c.Data["json"] = JSONStruct{"error", 1005, v, "获取失败"}
	} else {
		c.Data["json"] = JSONStruct{"success", 0, v, "获取成功"}
	}
	c.ServeJSON()
}

// GetTeacher ...
// @Title 全部教师列表
// @Description 全部教师列表
// @Param	kindergarten_id       query	int	     true		"幼儿园ID"
// @Param	class_type            query	int	     true		"班级类型"
// @Param	status                query	int	     false		"状态"
// @Param	search                query	int	     false		"搜索条件"
// @Param	page                  query	int	     false		"页数"
// @Param	per_page              query	int	     false		"每页显示条数"
// @Success 200 {object} models.Teacher
// @Failure 403
// @router / [get]
func (c *TeacherController) GetTeacher() {
	var prepage int = 20
	var page int
	var kindergarten_id int
	var class_type int
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
	if v, err := c.GetInt("class_type"); err == nil {
		class_type = v
	}
	if v, err := c.GetInt("status", -1); err == nil {
		status = v
	}
	v := models.GetTeacher(kindergarten_id, class_type, status, search, page, prepage)
	if v == nil {
		c.Data["json"] = JSONStruct{"error", 1005, v, "获取失败"}
	} else {
		c.Data["json"] = JSONStruct{"success", 0, v, "获取成功"}
	}
	c.ServeJSON()
}

// Delete ...
// @Title Delete
// @Description 删除教师
// @Param	teacher_id		path 	string	true		"教师ID"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (c *TeacherController) Delete() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	v := models.DeleteTeacher(id)
	if v == nil {
		c.Data["json"] = JSONStruct{"error", 1004, nil, "删除失败"}
	} else {
		c.Data["json"] = JSONStruct{"success", 0, nil, "删除成功"}
	}
	c.ServeJSON()
}

// GetTeacherInfo ...
// @Title Get Teacher Info
// @Description 教师详情
// @Param	teacher_id       query	int	 true		"教师编号"
// @Param	page             query	int	 false		"页数"
// @Param	per_page         query	int	 false		"每页显示条数"
// @Success 200 {object} models.Teacher
// @Failure 403 :教师编号为空
// @router /:id [get]
func (c *TeacherController) GetTeacherInfo() {
	var prepage int = 20
	var page int
	if v, err := c.GetInt("per_page"); err == nil {
		prepage = v
	}
	if v, err := c.GetInt("page"); err == nil {
		page = v
	}
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	v := models.GetTeacherInfo(id, page, prepage)
	if v == nil {
		c.Data["json"] = JSONStruct{"error", 1005, v, "获取失败"}
	} else {
		c.Data["json"] = JSONStruct{"success", 0, v, "获取成功"}
	}
	c.ServeJSON()
}

// Put ...
// @Title 编辑教师
// @Description 编辑教师
// @Param	id		    path 	int	               true		    "教师编号"
// @Param	body		body 	models.Animation	true		"param(json)"
// @Success 200 {object} models.Animation
// @Failure 403 :id is not int
// @router /:id [put]
func (c *TeacherController) Put() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	v := models.Teacher{Id: id}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		v := models.UpdateTeacher(&v)
		if v == nil {
			c.Data["json"] = JSONStruct{"error", 1003, err.Error(), "编辑失败"}
		} else {
			c.Data["json"] = JSONStruct{"success", 0, nil, "编辑成功"}
		}
		c.ServeJSON()
	} else {
		c.Data["json"] = JSONStruct{"error", 1001, err.Error(), "字段必须为json格式"}
		c.ServeJSON()
	}
}
