package controllers

import (
	"kindergarten-service-go/models"
)

type ExceptionalChildController struct {
	BaseController
}

func (c *ExceptionalChildController) URLMapping() {
	c.Mapping("Get", c.Get)
}



// @Title 特殊儿童列表/搜索特殊儿童
// @Description 特殊儿童列表/搜索特殊儿童
// @Param	keyword	   	 query	string	 false	"关键字(特殊儿童姓名/特殊儿童过敏源)"
// @Param	page	  	 query	int	 	 false	"当前页"
// @Success 0 			 {object}  models.ExceptionalChild
// @Failure 1005 获取失败
// @router / [get]
func (c *ExceptionalChildController) GetSearch() {
	// 关键字
	keyword := c.GetString("keyword")
	// page_num
	page, _ := c.GetInt64("page")

	// limit
	limit, _ := c.GetInt64("per_page")

	if info, err := models.GetAllExceptionalChild("", 0, page, limit, keyword); err == nil {
		c.Data["json"] = JSONStruct{"success", 0, info, "获取成功"}
	} else {
		c.Data["json"] = JSONStruct{"error", 1005, err, "获取失败"}
	}
	c.ServeJSON()
}
