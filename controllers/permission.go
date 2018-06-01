package controllers

import (
	"encoding/json"
	"kindergarten-service-go/models"
	"strconv"

	"github.com/astaxie/beego"
)

//权限
type PermissionController struct {
	beego.Controller
}

// Post ...
// @Title 保存权限
// @Description 保存权限
// @Param	body		body 	models.Permission	true		"body for Permission content"
// @Success 201 {int} models.Permission
// @Failure 403 body is empty
// @router / [post]
func (c *PermissionController) Post() {
	var v models.Permission
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		if _, err := models.AddPermission(&v); err == nil {
			c.Ctx.Output.SetStatus(201)
			c.Data["json"] = v
		} else {
			c.Data["json"] = err.Error()
		}
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}

// GetOne ...
// @Title Get One
// @Description get Permission by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Permission
// @Failure 403 :id is empty
// @router /:id [get]
func (c *PermissionController) GetOne() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	v, err := models.GetPermissionById(id)
	if err != nil {
		c.Data["json"] = err.Error()
	} else {
		c.Data["json"] = v
	}
	c.ServeJSON()
}

// GetAll ...
// @Title 权限列表
// @Description 权限列表
// @Param	id                    query	int	     false		"自增ID"
// @Param	page                  query	int	     false		"页数"
// @Param	per_page              query	int	     false		"每页显示条数"
// @Success 200 {object} models.Permission
// @Failure 403
// @router / [get]
func (c *PermissionController) GetAll() {
	var prepage int = 20
	var page int
	if v, err := c.GetInt("per_page"); err == nil {
		prepage = v
	}
	if v, err := c.GetInt("page"); err == nil {
		page = v
	}
	id, _ := c.GetInt("id")
	v := models.GetAllPermission(id, page, prepage)
	if v == nil {
		c.Data["json"] = JSONStruct{"error", 1005, nil, "获取失败"}
	} else {
		c.Data["json"] = JSONStruct{"success", 0, v, "获取成功"}
	}
	c.ServeJSON()
}

// Put ...
// @Title Put
// @Description update the Permission
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.Permission	true		"body for Permission content"
// @Success 200 {object} models.Permission
// @Failure 403 :id is not int
// @router /:id [put]
func (c *PermissionController) Put() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	v := models.Permission{Id: id}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		if err := models.UpdatePermissionById(&v); err == nil {
			c.Data["json"] = "OK"
		} else {
			c.Data["json"] = err.Error()
		}
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}

// Delete ...
// @Title Delete
// @Description delete the Permission
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (c *PermissionController) Delete() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	if err := models.DeletePermission(id); err == nil {
		c.Data["json"] = "OK"
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}
