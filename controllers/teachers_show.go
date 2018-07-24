package controllers

import (
	"kindergarten-service-go/models"
	"log"
	"strconv"

	"github.com/astaxie/beego/validation"
)

//教师展示
type TeachersShowController struct {
	BaseController
}

// Store ...
// @Title 保存教师展示
// @Description 保存教师展示
// @Param	teacher_id		     query      int   		   true		    "教师id"
// @Param	introduction		 query      string 	    	true		"介绍"
// @Param	kindergarten_id	     query	    int 	    	true		"幼儿园ID"
// @Success 201 {int} models.TeachersShow
// @Failure 403 body is empty
// @router / [post]
func (c *TeachersShowController) Store() {
	teacher_id, _ := c.GetInt("teacher_id")
	introduction := c.GetString("introduction")
	kindergarten_id, _ := c.GetInt("kindergarten_id")
	valid := validation.Validation{}
	valid.Required(teacher_id, "teacher_id").Message("教师id不能为空")
	valid.Required(introduction, "introduction").Message("介绍不能为空")
	valid.Required(kindergarten_id, "kindergarten_id").Message("幼儿园编号不能为空")
	if valid.HasErrors() {
		log.Println(valid.Errors)
		c.Data["json"] = JSONStruct{"error", 1001, nil, valid.Errors[0].Message}
		c.ServeJSON()
	} else {
		err := models.AddTeachersShow(teacher_id, introduction, kindergarten_id)
		if err != nil {
			c.Data["json"] = JSONStruct{"error", 1003, err.Error(), "保存失败"}
		} else {
			c.Data["json"] = JSONStruct{"success", 0, nil, "保存成功"}
		}
		c.ServeJSON()
	}
}

// TeachersShowAll ...
// @Title 教师展示列表
// @Description 教师展示列表
// @Param	page     query	int	 false		"页数"
// @Param	per_page query	int	 false		"每页显示条数"
// @Success 200 {object} models.TeachersShow
// @Failure 403
// @router / [get]
func (c *TeachersShowController) TeachersShowAll() {
	page, _ := c.GetInt("page")
	prepage, _ := c.GetInt("per_page", 20)
	kindergarten_id, _ := c.GetInt("kindergarten_id")
	valid := validation.Validation{}
	valid.Required(kindergarten_id, "kindergarten_id").Message("幼儿园编号不能为空")
	if valid.HasErrors() {
		log.Println(valid.Errors)
		c.Data["json"] = JSONStruct{"error", 1001, nil, valid.Errors[0].Message}
		c.ServeJSON()
	} else {
		v, err := models.TeachersShowAll(page, prepage, kindergarten_id)
		if err != nil {
			c.Data["json"] = JSONStruct{"error", 1005, nil, err.Error()}
		} else {
			c.Data["json"] = JSONStruct{"success", 0, v, "获取成功"}
		}
		c.ServeJSON()
	}
}

// TeachersShowOne ...
// @Title 公告详情
// @Description Web-公告详情
// @Param	id       query	int	 true		"主键ID"
// @Param	page     query	int	 false		"页数"
// @Param	per_page query	int	 false		"每页显示条数"
// @Success 200 {object} models.TeachersShow
// @Failure 403 :编号为空
// @router /:id [get]
func (c *TeachersShowController) TeachersShowOne() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	v, err := models.TeachersShowOne(id)
	if err != nil {
		c.Data["json"] = JSONStruct{"error", 1005, err.Error(), "获取失败"}
	} else {
		c.Data["json"] = JSONStruct{"success", 0, v, "获取成功"}
	}
	c.ServeJSON()
}

// Delete ...
// @Title 删除教师展示
// @Description 删除教师展示
// @Param	id		path 	string	true		"用户编号"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (c *TeachersShowController) Delete() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	err := models.DeleteTeachersShow(id)
	if err != nil {
		c.Data["json"] = JSONStruct{"error", 1004, err.Error(), "删除失败"}
	} else {
		c.Data["json"] = JSONStruct{"success", 0, nil, "删除成功"}
	}
	c.ServeJSON()
}

// Put ...
// @Title 编辑教师展示
// @Description 编辑教师展示
// @Param	id		             path 	    int	            true		    "教师编号"
// @Param	teacher_id		     query      int   		    true		    "教师id"
// @Param	introduction		 query      string 	    	true		"介绍"
// @Param	kindergarten_id	     query	    int 	    	true		"幼儿园ID"
// @Success 200 {object} models.Teacher
// @Failure 403 :id is not int
// @router /:id [put]
func (c *TeachersShowController) Put() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	teacher_id, _ := c.GetInt("teacher_id")
	introduction := c.GetString("introduction")
	kindergarten_id, _ := c.GetInt("kindergarten_id")
	err := models.UpdateTeachersShow(id, teacher_id, introduction, kindergarten_id)
	if err != nil {
		c.Data["json"] = JSONStruct{"error", 1003, err.Error(), "编辑失败"}
	} else {
		c.Data["json"] = JSONStruct{"success", 0, nil, "编辑成功"}
	}
	c.ServeJSON()
}
