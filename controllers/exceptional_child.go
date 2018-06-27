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
// @Param	page	  	 query	int	 	 false	"当前页，默认为1"
// @Param   per_page     query  int      false  "每页显示条数，默认为10"
// @Success 0 {json} 	JSONStruct
// @Failure 1005 获取失败
// @router / [get]
func(c *ExceptionalChildController) Get() {
	// 当前页
	page, _ := c.GetInt("page")
	// 每页显示条数
	per_page, _ := c.GetInt("per_page")
	// 关键字
	keyword := c.GetString("keyword")


	// 获取访客列表
	if info, err := models.GetExceptionalChild(page, per_page, keyword); err == nil {
		c.Data["json"] = JSONStruct{"success", 0, info, "获取成功"}
	} else {
		c.Data["json"] = JSONStruct{"error", 1005, err, "获取失败"}
	}

	c.ServeJSON()
}
