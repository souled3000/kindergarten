package controllers

import (
	"kindergarten-service-go/models"
)

//幼儿园 - 访客
type KindergartenVisitorsController struct {
	BaseController
}

// @Title 幼儿园访客列表
// @Description 幼儿园访客列表
// @Param	page	  query	int	 false	"当前页"
// @Param	per_page  query	int	 false	"每页显示条数"
// @Success 0 {string} success
// @Failure 1005 获取失败
// @router / [get]
func (c *KindergartenVisitorsController) GetAll() {
	// 当前页
	page, _ := c.GetInt("page")
	// 每页显示条数
	per_page, _ := c.GetInt("per_page")

	// 获取访客列表
	if info, err := models.GetVisitors(page, per_page); err == nil {
		c.Data["json"] = JSONStruct{"success", 0, info, "获取成功"}
	} else {
		c.Data["json"] = JSONStruct{"error", 1005, err, "获取失败"}
	}

	c.ServeJSON()
}
