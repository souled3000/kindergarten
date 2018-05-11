package controllers

import (
	"github.com/astaxie/beego"
)

type ErrorController struct {
	beego.Controller
}

func (c *ErrorController) Error404() {
	c.Data["json"] = JSONStruct{"error", 404, nil, "获取失败"}
	c.ServeJSON()
}
