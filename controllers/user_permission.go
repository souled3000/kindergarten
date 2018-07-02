package controllers

import (
	"kindergarten-service-go/models"
	"strconv"

	"github.com/astaxie/beego/validation"
)

//用户分配权限
type UserPermissionController struct {
	BaseController
}

// Post ...
// @Title 设置权限
// @Description 设置权限
// @Param	user_id		    body 	int	  true		"用户ID"
// @Param	role		    body 	int	  true		"角色ID(json)"
// @Param	permission		body 	int	  true		"权限(json)"
// @Param	group		    body 	int	  true		"圈子(班级类型)(json)"
// @Success 201 {int} models.UserPermission
// @Failure 403 body is empty
// @router / [post]
func (c *UserPermissionController) Post() {
	user_id, _ := c.GetInt("user_id")
	role := c.GetString("role")
	permission := c.GetString("permission")
	group := c.GetString("group")
	valid := validation.Validation{}
	valid.Required(user_id, "user_id").Message("用户ID不能为空")
	if valid.HasErrors() {
		c.Data["json"] = JSONStruct{"error", 1001, nil, valid.Errors[0].Message}
		c.ServeJSON()
	} else {
		_, err := models.AddUserPermission(user_id, role, permission, group)
		if err != nil {
			c.Data["json"] = JSONStruct{"error", 1003, nil, err.Error()}
		} else {
			c.Data["json"] = JSONStruct{"success", 0, nil, "保存成功"}
		}
		c.ServeJSON()
	}
}

// GetOne ...
// @Title 查看用户权限
// @Description 查看用户权限
// @Param	user_id		path 	string	true		"用户ID"
// @Success 200 {object} models.UserPermission
// @Failure 403 :id is empty
// @router /:id [get]
func (c *UserPermissionController) GetOne() {
	idStr := c.Ctx.Input.Param(":id")
	user_id, _ := strconv.Atoi(idStr)
	kindergarten_id, _ := c.GetInt("kindergarten_id")
	v, err := models.GetUserPermissionById(user_id, kindergarten_id)
	if err != nil {
		c.Data["json"] = JSONStruct{"error", 1003, nil, err.Error()}
	} else {
		c.Data["json"] = JSONStruct{"success", 0, v, "获取成功"}
	}
	c.ServeJSON()
}

// GetUserPermission ...
// @Title 查看用户权限标识
// @Description 查看用户权限标识
// @Param	user_id		path 	string	true		"用户ID"
// @Success 200 {object} models.UserPermission
// @Failure 403 :id is empty
// @router /user/:id [get]
func (c *UserPermissionController) GetUserPermission() {
	idStr := c.Ctx.Input.Param(":id")
	user_id, _ := strconv.Atoi(idStr)
	v, err := models.GetUserIdentificationById(user_id)
	if err != nil {
		c.Data["json"] = JSONStruct{"error", 1003, nil, err.Error()}
	} else {
		c.Data["json"] = JSONStruct{"success", 0, v, "获取成功"}
	}
	c.ServeJSON()
}

// GetGroupPermission ...
// @Title 查看用户圈子权限标识
// @Description 查看用户圈子权限标识
// @Param	user_id		path 	string	true		"用户ID"
// @Success 200 {object} models.UserPermission
// @Failure 403 :id is empty
// @router /group/:id [get]
func (c *UserPermissionController) GetGroupPermission() {
	idStr := c.Ctx.Input.Param(":id")
	user_id, _ := strconv.Atoi(idStr)
	v, err := models.GetGroupIdentificationById(user_id)
	if err != nil {
		c.Data["json"] = JSONStruct{"error", 1003, nil, err.Error()}
	} else {
		c.Data["json"] = JSONStruct{"success", 0, v, "获取成功"}
	}
	c.ServeJSON()
}

// Put ...
// @Title 更新用户权限
// @Description 更新用户权限
// @Param	id		    formData 	int	    true		"用户ID"
// @Param	role		    formData 	string	true		"角色ID(json)"
// @Param	permission		formData 	string	true		"权限ID(json)"
// @Param	group		    formData 	string	true		"圈子ID(json)"
// @Success 200 {object} models.UserPermission
// @Failure 403 :id is not int
// @router /:id [put]
func (c *UserPermissionController) Put() {
	idStr := c.Ctx.Input.Param(":id")
	user_id, _ := strconv.Atoi(idStr)
	role := c.GetString("role")
	permission := c.GetString("permission")
	group := c.GetString("group")
	valid := validation.Validation{}
	valid.Required(user_id, "user_id").Message("用户ID不能为空")
	valid.Required(permission, "permission").Message("权限不能为空")
	if valid.HasErrors() {
		c.Data["json"] = JSONStruct{"error", 1001, nil, valid.Errors[0].Message}
		c.ServeJSON()
	} else {
		_, err := models.UpdateUserPermissionById(user_id, role, permission, group)
		if err != nil {
			c.Data["json"] = JSONStruct{"error", 1003, nil, err.Error()}
		} else {
			c.Data["json"] = JSONStruct{"success", 0, nil, "更新成功"}
		}
		c.ServeJSON()
	}
}

// GroupAll ...
// @Title 筛选圈子
// @Description 筛选圈子
// @Param	user_id		        path 	int	true		"用户ID"
// @Param	class_type		    path 	int	true		"班级类型"
// @Param	role		        path 	int	true		"身份(1 kindergarten_id,role,class_type  5 user_id,role,class_type)"
// @Param	kindergarten_id		path 	int	true		"幼儿园id"
// @Success 200 {object} models.UserPermission
// @Failure 403 :id is empty
// @router /group_all [get]
func (c *UserPermissionController) GroupAll() {
	role, _ := c.GetInt("role")
	user_id, _ := c.GetInt("user_id")
	class_type, _ := c.GetInt("class_type")
	kindergarten_id, _ := c.GetInt("kindergarten_id")
	v, err := models.GetGroupAll(user_id, class_type, role, kindergarten_id)
	if err != nil {
		c.Data["json"] = JSONStruct{"error", 1003, nil, err.Error()}
	} else {
		c.Data["json"] = JSONStruct{"success", 0, v, "获取成功"}
	}
	c.ServeJSON()
}
