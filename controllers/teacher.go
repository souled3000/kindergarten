package controllers

import (
	"encoding/json"
	"kindergarten-service-go/models"
	"strconv"

	"github.com/astaxie/beego/validation"
)

//教师
type TeacherController struct {
	BaseController
}

// GetTeacher ...
// @Title 全部教师列表
// @Description 全部教师列表
// @Param	kindergarten_id       query	int	     true		"幼儿园ID"
// @Param	status                query	int	     false		"状态"
// @Param	search                query	int	     false		"搜索条件"
// @Param	page                  query	int	     false		"页数"
// @Param	per_page              query	int	     false		"每页显示条数"
// @Success 200 {object} models.Teacher
// @Failure 403
// @router / [get]
func (c *TeacherController) GetTeacher() {
	search := c.GetString("search")
	prepage, _ := c.GetInt("per_page", 20)
	page, _ := c.GetInt("page")
	kindergarten_id, _ := c.GetInt("kindergarten_id")
	status, _ := c.GetInt("status", -1)
	valid := validation.Validation{}
	valid.Required(kindergarten_id, "kindergarten_id").Message("幼儿园编号不能为空")
	if valid.HasErrors() {
		c.Data["json"] = JSONStruct{"error", 1001, nil, valid.Errors[0].Message}
		c.ServeJSON()
	} else {
		v := models.GetTeacher(kindergarten_id, status, search, page, prepage)
		if v == nil {
			c.Data["json"] = JSONStruct{"error", 1005, nil, "获取失败"}
		} else {
			c.Data["json"] = JSONStruct{"success", 0, v, "获取成功"}
		}
		c.ServeJSON()
	}
}

// GetClass ...
// @Title 班级列表
// @Description 班级列表
// @Param	kindergarten_id       query	int	     true		"幼儿园ID"
// @Param	class_type            query	int	     true		"班级类型"
// @Param	page                  query	int	     false		"页数"
// @Param	per_page              query	int	     false		"每页显示条数"
// @Success 200 {object} models.Teacher
// @Failure 403
// @router /class [get]
func (c *TeacherController) GetClass() {
	prepage, _ := c.GetInt("per_page", 20)
	page, _ := c.GetInt("page")
	kindergarten_id, _ := c.GetInt("kindergarten_id")
	class_type, _ := c.GetInt("class_type")
	v := models.GetClass(kindergarten_id, class_type, page, prepage)
	if v == nil {
		c.Data["json"] = JSONStruct{"error", 1005, nil, "获取失败"}
	} else {
		c.Data["json"] = JSONStruct{"success", 0, v, "获取成功"}
	}
	c.ServeJSON()
}

// Delete ...
// @Title Delete
// @Description 删除教师
// @Param	teacher_id		path 	int	true		"教师ID"
// @Param	status		    path 	int	true		"状态(status 0:未分班 2:离职)"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (c *TeacherController) Delete() {
	class_type, _ := c.GetInt("class_type")
	status, _ := c.GetInt("status")
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	err := models.DeleteTeacher(id, status, class_type)
	if err != nil {
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
// @Success 200 {object} models.Teacher
// @Failure 403 :教师编号为空
// @router /:id [get]
func (c *TeacherController) GetTeacherInfo() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	v, err := models.GetTeacherInfo(id)
	if err != nil {
		c.Data["json"] = JSONStruct{"error", 1005, nil, "获取失败"}
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
		err := models.UpdateTeacher(&v)
		if err != nil {
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

// Post ...
// @Title 教师-录入信息
// @Description 教师-录入信息
// @Param	body		body 	models.Animation	true		"json"
// @Success 201 {int} models.Animation
// @Failure 403 body is empty
// @router / [post]
func (c *TeacherController) Post() {
	var v models.Teacher
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		valid := validation.Validation{}
		valid.Required(v.KindergartenId, "KindergartenId").Message("幼儿园编号不能为空")
		valid.Required(v.UserId, "UserId").Message("用户编号不能为空")
		valid.Required(v.Birthday, "Birthday").Message("出生年月日不能为空")
		valid.Required(v.Name, "Name").Message("用户名不能为空")
		valid.Required(v.Number, "Number").Message("教职工编号不能为空")
		valid.Required(v.NationOrReligion, "NationOrReligion").Message("民族或宗教不能为空")
		valid.Required(v.NativePlace, "NativePlace").Message("籍贯不能为空")
		valid.Required(v.EnterJobTime, "EnterJobTime").Message("参加工作时间不能为空")
		valid.Required(v.Address, "Address").Message("住址不能为空")
		valid.Required(v.EmergencyContact, "EmergencyContact").Message("紧急联系人不能为空")
		valid.Required(v.EmergencyContactPhone, "EmergencyContactPhone").Message("紧急联系人电话不能为空")
		valid.Required(v.Source, "Source").Message("来源不能为空")
		valid.Required(v.TeacherCertificationNumber, "TeacherCertificationNumber").Message("教师认证编号不能为空")
		valid.Required(v.Phone, "Phone").Message("手机号不能为空")
		valid.Required(v.EnterGardenTime, "EnterGardenTime").Message("进入本园时间不能为空")
		if valid.HasErrors() {
			c.Data["json"] = JSONStruct{"error", 1001, nil, valid.Errors[0].Message}
			c.ServeJSON()
		} else {
			err := models.AddTeacher(&v)
			if err != nil {
				c.Data["json"] = JSONStruct{"error", 1003, err.Error(), "保存失败"}
			} else {
				c.Data["json"] = JSONStruct{"success", 0, nil, "保存成功"}
			}
			c.ServeJSON()
		}
	} else {
		c.Data["json"] = JSONStruct{"error", 1001, err.Error(), "字段必须为json格式"}
		c.ServeJSON()
	}
}

// RemoveTeacher ...
// @Title 移除教师
// @Description 移除教师
// @Param	teacher_id		    path 	    int	true		    "教师ID"
// @Param	class_id		    path 	    int	true		    "班级ID"
// @Success 200 {string} delete success!
// @Failure 403 teacher_id is empty
// @router /remove [delete]
func (c *TeacherController) RemoveTeacher() {
	teacher_id, _ := c.GetInt("teacher_id")
	class_id, _ := c.GetInt("class_id")
	valid := validation.Validation{}
	valid.Required(teacher_id, "teacher_id").Message("教师ID不能为空")
	valid.Required(class_id, "class_id").Message("班级ID不能为空")
	if valid.HasErrors() {
		c.Data["json"] = JSONStruct{"error", 1001, nil, valid.Errors[0].Message}
		c.ServeJSON()
	} else {
		err := models.RemoveTeacher(teacher_id, class_id)
		if err != nil {
			c.Data["json"] = JSONStruct{"error", 1004, err.Error(), "移除失败"}
		} else {
			c.Data["json"] = JSONStruct{"success", 0, nil, "移除成功"}
		}
		c.ServeJSON()
	}
}

// OrganizationalTeacher ...
// @Title 组织框架教师列表
// @Description 组织框架教师列表
// @Param	kindergarten_id       query	int	     true		"幼儿园ID"
// @Param	type                  query	int	     true		"年级组标识(1 年级组)"
// @Param	person                query	int	     true		"是否为负责人(1 负责 2 不是负责人)"
// @Success 200 {object} models.Teacher
// @Failure 403
// @router /organizational_teacher [get]
func (c *TeacherController) OrganizationalTeacher() {
	ty, _ := c.GetInt("type")
	person, _ := c.GetInt("person")
	class_id, _ := c.GetInt("class_id")
	kindergarten_id, _ := c.GetInt("kindergarten_id")
	valid := validation.Validation{}
	valid.Required(person, "person").Message("身份不能为空")
	valid.Required(kindergarten_id, "kindergarten_id").Message("幼儿园编号不能为空")
	if valid.HasErrors() {
		c.Data["json"] = JSONStruct{"error", 1001, nil, valid.Errors[0].Message}
		c.ServeJSON()
	} else {
		v, err := models.OrganizationalTeacher(kindergarten_id, ty, person, class_id)
		if err != nil {
			c.Data["json"] = JSONStruct{"error", 1005, nil, "获取失败"}
		} else {
			c.Data["json"] = JSONStruct{"success", 0, v, "获取成功"}
		}
		c.ServeJSON()
	}
}

// Teacher ...
// @Title 教师筛选列表-前台
// @Description 教师筛选列表-前台
// @Param	class_id       query	int	     true		"班级ID"
// @Success 200 {object} models.Student
// @Failure 403
// @router /filter_teacher [get]
func (c *TeacherController) Teacher() {
	class_id, _ := c.GetInt("class_id")
	v, err := models.FilterTeacher(class_id)
	if err != nil {
		c.Data["json"] = JSONStruct{"error", 1005, nil, "获取失败"}
	} else {
		c.Data["json"] = JSONStruct{"success", 0, v, "获取成功"}
	}
	c.ServeJSON()
}
