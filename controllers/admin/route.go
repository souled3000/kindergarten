package admin

import (
	"kindergarten-service-go/models"
	"strconv"

	"github.com/astaxie/beego/validation"
)

//路由
type RouteController struct {
	BaseController
}

// Post ...
// @Title 保存路由
// @Description 保存路由
// @Param	name		    body 	string		"路由名称"
// @Param	route		body 	string	    "路由"
// @Success 201 {int} models.Route
// @Failure 403 body is empty
// @router / [post]
func (c *RouteController) Post() {
	name := c.GetString("name")
	route := c.GetString("route")
	valid := validation.Validation{}
	valid.Required(name, "name").Message("路由名称不能为空")
	valid.Required(route, "route").Message("路由不能为空")
	if valid.HasErrors() {
		c.Data["json"] = JSONStruct{"error", 1001, nil, valid.Errors[0].Message}
		c.ServeJSON()
	} else {
		l := models.AddRoute(name, route)
		if l == nil {
			c.Data["json"] = JSONStruct{"error", 1003, nil, "保存失败"}
		} else {
			c.Data["json"] = JSONStruct{"success", 0, nil, "保存成功"}
		}
		c.ServeJSON()
	}
}

// GetOne ...
// @Title 路由详情
// @Description 路由详情
// @Param	id		path 	string	true		"路由ID"
// @Success 200 {object} models.Route
// @Failure 403 :id is empty
// @router /:id [get]
func (c *RouteController) GetOne() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	r := models.GetRouteById(id)
	if r == nil {
		c.Data["json"] = JSONStruct{"error", 1003, nil, "获取失败"}
	} else {
		c.Data["json"] = JSONStruct{"success", 0, r, "获取成功"}
	}
	c.ServeJSON()
}

// GetAll ...
// @Title 路由列表
// @Description 路由列表
// @Param	page                  query	int	     false		"页数"
// @Param	per_page              query	int	     false		"每页显示条数"
// @Success 200 {object} models.Route
// @Failure 403
// @router / [get]
func (c *RouteController) GetAll() {
	var prepage int = 20
	var page int
	if v, err := c.GetInt("per_page"); err == nil {
		prepage = v
	}
	if v, err := c.GetInt("page"); err == nil {
		page = v
	}
	v := models.GetAllRoute(page, prepage)
	if v == nil {
		c.Data["json"] = JSONStruct{"error", 1005, nil, "获取失败"}
	} else {
		c.Data["json"] = JSONStruct{"success", 0, v, "获取成功"}
	}
	c.ServeJSON()
}

// Put ...
// @Title Put
// @Description 编辑路由
// @Param	id		    path 	int       	true		"路由ID"
// @Param	name		    body 	string	    false	"路由名称"
// @Param	route		body 	string	    false	"路由"
// @Success 200 {object} models.Route
// @Failure 403 :id is not int
// @router /:id [put]
func (c *RouteController) Put() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	name := c.GetString("name")
	route := c.GetString("route")
	v := models.UpdateRouteById(id, name, route)
	if v == nil {
		c.Data["json"] = JSONStruct{"error", 1003, nil, "编辑失败"}
	} else {
		c.Data["json"] = JSONStruct{"success", 0, nil, "编辑成功"}
	}
	c.ServeJSON()
}

// Delete ...
// @Title 删除路由
// @Description 删除路由
// @Param	id		path 	string	true		"路由ID"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (c *RouteController) Delete() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	v := models.DeleteRoute(id)
	if v == nil {
		c.Data["json"] = JSONStruct{"error", 1003, nil, "删除失败"}
	} else {
		c.Data["json"] = JSONStruct{"success", 0, nil, "删除成功"}
	}
	c.ServeJSON()
}
