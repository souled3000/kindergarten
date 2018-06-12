package controllers

import (
	"kindergarten-service-go/models"
	"strconv"

	"github.com/astaxie/beego"
)

//幼儿园
type KindergartenController struct {
	beego.Controller
}

// GetIntroduceInfo ...
// @Title 幼儿园介绍详情
// @Description 幼儿园介绍详情
// @Param	id		path 	string	true		"幼儿园ID"
// @Success 200 {object} models.Kindergarten
// @Failure 403 :id is empty
// @router /:id [get]
func (c *KindergartenController) GetIntroduceInfo() {
	prepage, _ := c.GetInt("per_page", 20)
	page, _ := c.GetInt("page")
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	v := models.GetKindergartenById(id, page, prepage)
	if v == nil {
		c.Data["json"] = JSONStruct{"error", 1005, v, "获取失败"}
	} else {
		c.Data["json"] = JSONStruct{"success", 0, v, "获取成功"}
	}
	c.ServeJSON()
}
