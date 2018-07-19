package admin

import (
	"kindergarten-service-go/models"
	"strconv"

	"github.com/astaxie/beego/validation"
)

//权限
type PermissionController struct {
	BaseController
}

// Post ...
// @Title 保存权限
// @Description 保存权限
// @Param	name		        body 	string	true		"权限名称"
// @Param	identification		body 	string	true		"权限标识"
// @Param	parent_id		    body 	int 	true		"父级ID"
// @Param	route		        body 	int 	false		"路由功能"
// @Success 201 {int} models.Permission
// @Failure 403 body is empty
// @router / [post]
func (c *PermissionController) Post() {
	name := c.GetString("name")
	identification := c.GetString("identification")
	parent_id, _ := c.GetInt("parent_id", -2)
	route := c.GetString("route")
	valid := validation.Validation{}
	valid.Required(name, "name").Message("权限名称不能为空")
	valid.Required(identification, "identification").Message("权限标识不能为空")
	if valid.HasErrors() {
		c.Data["json"] = JSONStruct{"error", 1001, nil, valid.Errors[0].Message}
		c.ServeJSON()
	} else {
		err := models.AddPermission(name, identification, parent_id, route)
		if err != nil {
			c.Data["json"] = JSONStruct{"error", 1003, nil, err.Error()}
		} else {
			c.Data["json"] = JSONStruct{"success", 0, nil, "保存成功"}
		}
		c.ServeJSON()
	}
}

// GetOne ...
// @Title 权限详情
// @Description 权限详情
// @Param	id		path 	string	true		"权限ID"
// @Success 200 {object} models.Permission
// @Failure 403 :id is empty
// @router /:id [get]
func (c *PermissionController) GetOne() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	v, err := models.GetPermissionById(id)
	if err != nil {
		c.Data["json"] = JSONStruct{"error", 1005, nil, "获取失败"}
	} else {
		c.Data["json"] = JSONStruct{"success", 0, v, "获取成功"}
	}
	c.ServeJSON()
}

// GetAll ...
// @Title 权限列表
// @Description 权限列表
// @Param	page                  query	int	     false		"页数"
// @Param	per_page              query	int	     false		"每页显示条数"
// @Success 200 {object} models.Permission
// @Failure 403
// @router / [get]
func (c *PermissionController) GetAll() {
	prepage, _ := c.GetInt("per_page", 20)
	page, _ := c.GetInt("page")
	v, err := models.GetAllPermission(page, prepage)
	if err != nil {
		c.Data["json"] = JSONStruct{"error", 1005, nil, "获取失败"}
	} else {
		c.Data["json"] = JSONStruct{"success", 0, v, "获取成功"}
	}
	c.ServeJSON()
}

// Option ...
// @Title 权限选项
// @Description 权限选项
// @Param	id		path 	string	true		"id"
// @Success 200 {object} models.Permission
// @Failure 403
// @router /option [get]
func (c *PermissionController) Option() {
	v := models.PermissionOption()
	if v == nil {
		c.Data["json"] = JSONStruct{"error", 1005, nil, "获取失败"}
	} else {
		c.Data["json"] = JSONStruct{"success", 0, v, "获取成功"}
	}
	c.ServeJSON()
}

// Update ...
// @Title 编辑权限
// @Description 编辑权限
// @Param	route		        body 	int 	false		"路由功能"
// @Success 201 {int} models.Permission
// @Failure 403 body is empty
// @router /:id [put]
func (c *PermissionController) Update() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	routeId := c.GetString("routeId")
	valid := validation.Validation{}
	valid.Required(routeId, "routeId").Message("路由不能空")
	if valid.HasErrors() {
		c.Data["json"] = JSONStruct{"error", 1001, nil, valid.Errors[0].Message}
		c.ServeJSON()
	} else {
		err := models.UpdatePermission(id, routeId)
		if err != nil {
			c.Data["json"] = JSONStruct{"error", 1003, err.Error(), "编辑失败"}
		} else {
			c.Data["json"] = JSONStruct{"success", 0, nil, "编辑成功"}
		}
		c.ServeJSON()
	}
}

// Delete ...
// @Title 删除权限
// @Description 删除权限
// @Param	id		        body 	int 	false		"权限id"
// @Success 201 {int} models.Permission
// @Failure 403 body is empty
// @router /:id [delete]
func (c *PermissionController) Delete() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	valid := validation.Validation{}
	valid.Required(id, "id").Message("权限id不能空")
	if valid.HasErrors() {
		c.Data["json"] = JSONStruct{"error", 1001, nil, valid.Errors[0].Message}
		c.ServeJSON()
	} else {
		err := models.DeletePermission(id)
		if err != nil {
			c.Data["json"] = JSONStruct{"error", 1003, err.Error(), "删除失败"}
		} else {
			c.Data["json"] = JSONStruct{"success", 0, nil, "删除成功"}
		}
		c.ServeJSON()
	}
}
