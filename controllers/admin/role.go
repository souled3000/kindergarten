package admin

import (
	"kindergarten-service-go/models"
	"strconv"

	"github.com/astaxie/beego/validation"
)

//角色
type RoleController struct {
	BaseController
}

// Post ...
// @Title Post
// @Description 添加角色
// @Param	body		body 	string  true		"角色名称"
// @Success 201 {int} models.Role
// @Failure 403 body is empty
// @router / [post]
func (c *RoleController) Post() {
	name := c.GetString("name")
	valid := validation.Validation{}
	valid.Required(name, "name").Message("角色名称不能为空")
	if valid.HasErrors() {
		c.Data["json"] = JSONStruct{"error", 1001, nil, valid.Errors[0].Message}
		c.ServeJSON()
	} else {
		v := models.AddRole(name)
		if v == nil {
			c.Data["json"] = JSONStruct{"error", 1005, nil, "保存失败"}
		} else {
			c.Data["json"] = JSONStruct{"success", 0, nil, "保存成功"}
		}
		c.ServeJSON()
	}
}

// GetOne ...
// @Title 角色详情
// @Description 角色详情
// @Param	id		path 	string	true		"角色id"
// @Success 200 {object} models.Role
// @Failure 403 :id is empty
// @router /:id [get]
func (c *RoleController) GetOne() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	v := models.GetRoleById(id)
	if v == nil {
		c.Data["json"] = JSONStruct{"error", 1005, nil, "获取失败"}
	} else {
		c.Data["json"] = JSONStruct{"success", 0, v, "获取成功"}
	}
	c.ServeJSON()
}

// GetAll ...
// @Title Get All
// @Description 角色列表
// @Param	page                  query	int	     false		"页数"
// @Param	per_page              query	int	     false		"每页显示条数"
// @Success 200 {object} models.Role
// @Failure 403
// @router / [get]
func (c *RoleController) GetAll() {
	var prepage int = 20
	var page int
	if v, err := c.GetInt("per_page"); err == nil {
		prepage = v
	}
	if v, err := c.GetInt("page"); err == nil {
		page = v
	}
	v := models.GetAllRole(page, prepage)
	if v == nil {
		c.Data["json"] = JSONStruct{"error", 1005, nil, "获取失败"}
	} else {
		c.Data["json"] = JSONStruct{"success", 0, v, "获取成功"}
	}
	c.ServeJSON()
}

// Put ...
// @Title Put
// @Description  编辑角色
// @Param	id		    path 	string	true		"角色ID"
// @Param	name		path 	string	true		"角色名称"
// @Success 200 {object} models.Role
// @Failure 403 :id is not int
// @router /:id [put]
func (c *RoleController) Put() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	name := c.GetString("name")
	valid := validation.Validation{}
	valid.Required(id, "id").Message("角色ID不能为空")
	valid.Required(name, "name").Message("角色名称不能为空")
	if valid.HasErrors() {
		c.Data["json"] = JSONStruct{"error", 1001, nil, valid.Errors[0].Message}
		c.ServeJSON()
	} else {
		v := models.UpdateRoleById(id, name)
		if v == nil {
			c.Data["json"] = JSONStruct{"error", 1005, nil, "编辑失败"}
		} else {
			c.Data["json"] = JSONStruct{"success", 0, nil, "编辑成功"}
		}
		c.ServeJSON()
	}
}
