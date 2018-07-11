package controllers

import (
	"kindergarten-service-go/models"
	"strconv"

	"github.com/astaxie/beego/validation"
)

//学生
type StudentController struct {
	BaseController
}

// Student ...
// @Title 学生列表-前台
// @Description 学生列表-前台
// @Param	kindergarten_id       query	int	     true		"幼儿园ID"
// @Param	page                  query	int	     false		"页数"
// @Param	per_page              query	int	     false		"每页显示条数"
// @Success 200 {object} models.Student
// @Failure 403
// @router / [get]
func (c *StudentController) Student() {
	kindergarten_id, _ := c.GetInt("kindergarten_id")
	v, err := models.ClassStudent(kindergarten_id)
	if err != nil {
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
func (c *StudentController) GetStudentClass() {
	prepage, _ := c.GetInt("per_page", 20)
	page, _ := c.GetInt("page")
	kindergarten_id, _ := c.GetInt("kindergarten_id")
	class_type, _ := c.GetInt("class_type")
	v := models.GetStudentClass(kindergarten_id, class_type, page, prepage)
	if v == nil {
		c.Data["json"] = JSONStruct{"error", 1005, nil, "获取失败"}
	} else {
		c.Data["json"] = JSONStruct{"success", 0, v, "获取成功"}
	}
	c.ServeJSON()
}

// RemoveStudent ...
// @Title RemoveStudent
// @Description 移除学生
// @Param	student_id		path 	    int	true		"学生ID"
// @Param	class_id		    path 	    int	true		"班级ID"
// @Success 200 {string} delete success!
// @Failure 403 student_id is empty
// @router /remove [delete]
func (c *StudentController) RemoveStudent() {
	student_id, _ := c.GetInt("student_id")
	class_id, _ := c.GetInt("class_id")
	err := models.RemoveStudent(class_id, student_id)
	if err != nil {
		c.Data["json"] = JSONStruct{"error", 1004, nil, err.Error()}
	} else {
		c.Data["json"] = JSONStruct{"success", 0, nil, "移除成功"}
	}
	c.ServeJSON()
}

// GetStudentInfo ...
// @Title Get Student Info
// @Description 学生详情
// @Param	student_id       query	int	 true		"学生编号"
// @Success 200 {object} models.Student
// @Failure 403 :学生编号为空
// @router /:id [get]
func (c *StudentController) GetOne() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	v, err := models.GetStudentInfo(id)
	if err != nil {
		c.Data["json"] = JSONStruct{"error", 1005, nil, err.Error()}
	} else {
		c.Data["json"] = JSONStruct{"success", 0, v, "获取成功"}
	}
	c.ServeJSON()
}

// UpdateStudent ...
// @Title 编辑学生
// @Description 编辑学生
// @Param	id		    path 	int	               true		    "学生编号"
// @Param	body		body 	models.Student	       true		"param(json)"
// @Success 200 {object} models.Student
// @Failure 403 :id is not int
// @router /:id [put]
func (c *StudentController) UpdateStudent() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	student := c.GetString("student")
	kinship := c.GetString("kinship")
	valid := validation.Validation{}
	valid.Required(student, "student").Message("学生信息不能为空")
	valid.Required(kinship, "kinship").Message("亲属信息不能为空")
	if valid.HasErrors() {
		c.Data["json"] = JSONStruct{"error", 1001, nil, valid.Errors[0].Message}
		c.ServeJSON()
	} else {
		_, err := models.UpdateStudent(id, student, kinship)
		if err != nil {
			c.Data["json"] = JSONStruct{"error", 1003, nil, err.Error()}
		} else {
			c.Data["json"] = JSONStruct{"success", 0, nil, "编辑成功"}
		}
		c.ServeJSON()
	}
}

// Post ...
// @Title 学生-录入信息
// @Description 学生-录入信息
// @Param	body		body 	models.Animation	true		"json"
// @Success 201 {int} models.Student
// @Failure 403 body is empty
// @router / [post]
func (c *StudentController) Post() {
	student := c.GetString("student")
	kinship := c.GetString("kinship")
	valid := validation.Validation{}
	valid.Required(student, "student").Message("学生信息不能为空")
	if valid.HasErrors() {
		c.Data["json"] = JSONStruct{"error", 1001, nil, valid.Errors[0].Message}
		c.ServeJSON()
	} else {
		_, err := models.AddStudent(student, kinship)
		if err != nil {
			c.Data["json"] = JSONStruct{"error", 1003, nil, err.Error()}
		} else {
			c.Data["json"] = JSONStruct{"success", 0, nil, "保存成功"}
		}
		c.ServeJSON()
	}
}

// Invite ...
// @Title 邀请学生/批量邀请
// @Description 邀请学生/批量邀请
// @Param	name		        body 	int   	true		"学生姓名(json)"
// @Param	baby_id		        body 	int   	true		"宝宝id(json)"
// @Param	kindergarten_id		body 	int   	true		"幼儿园id(json)"
// @Success 201 {int} models.Student
// @Failure 403 body is empty
// @router /invite [post]
func (c *StudentController) Invite() {
	student := c.GetString("student")
	valid := validation.Validation{}
	valid.Required(student, "student").Message("学生信息不能为空")
	if valid.HasErrors() {
		c.Data["json"] = JSONStruct{"error", 1001, nil, valid.Errors[0].Message}
		c.ServeJSON()
	} else {
		err := models.Invite(student)
		if err != nil {
			c.Data["json"] = JSONStruct{"error", 1003, nil, err.Error()}
		} else {
			c.Data["json"] = JSONStruct{"success", 0, nil, "邀请成功"}
		}
		c.ServeJSON()
	}
}

// DeleteStudent ...
// @Title DeleteStudent
// @Description 删除学生
// @Param	student_id		path 	int	true		"学生ID"
// @Param	status		    path 	int	true		"状态(status 0:未分班 2:离园)"
// @Param	type		        path 	int	true		"删除类型（type 0:学生离园 1:删除档案）"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (c *StudentController) DeleteStudent() {
	class_type, _ := c.GetInt("class_type")
	status, _ := c.GetInt("status")
	ty, _ := c.GetInt("type")
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	err := models.DeleteStudent(id, status, ty, class_type)
	if err != nil {
		c.Data["json"] = JSONStruct{"error", 1004, nil, "删除失败"}
	} else {
		c.Data["json"] = JSONStruct{"success", 0, nil, "删除成功"}
	}
	c.ServeJSON()
}

// GetBaby ...
// @Title GetBaby
// @Description 未激活baby
// @Param	kindergarten_id       query	 int	 true		"幼儿园id"
// @Success 200 {object} models.Student
// @Failure 403 :幼儿园id不能为空
// @router /baby [get]
func (c *StudentController) GetBaby() {
	kindergarten_id, _ := c.GetInt("kindergarten_id")
	v, err := models.GetBabyInfo(kindergarten_id)
	if err != nil {
		c.Data["json"] = JSONStruct{"error", 1005, nil, err.Error()}
	} else {
		c.Data["json"] = JSONStruct{"success", 0, v, "获取成功"}
	}
	c.ServeJSON()
}

// GetNameClass ...
// @Title GetNameClass
// @Description 学生名字获取班级
// @Param	name       query	 int	 true		"学生姓名"
// @Success 200 {object} models.Student
// @Failure 403 :幼儿园id不能为空
// @router /get_class [get]
func (c *StudentController) GetNameClass() {
	name := c.GetString("name")
	kindergarten_id, _ := c.GetInt("kindergarten_id")
	valid := validation.Validation{}
	valid.Required(name, "name").Message("学生姓名不能为空")
	valid.Required(kindergarten_id, "kindergarten_id").Message("幼儿园ID不能为空")
	if valid.HasErrors() {
		c.Data["json"] = JSONStruct{"error", 1001, nil, valid.Errors[0].Message}
		c.ServeJSON()
	} else {
		v, err := models.GetNameClass(name, kindergarten_id)
		if err != nil {
			c.Data["json"] = JSONStruct{"error", 1005, nil, err.Error()}
		} else {
			c.Data["json"] = JSONStruct{"success", 0, v, "获取成功"}
		}
		c.ServeJSON()
	}
}
