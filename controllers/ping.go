package controllers

import (
	"github.com/astaxie/beego"
)

type PingController struct {
	beego.Controller
}

// Ping ...
// @Title Ping
// @Description Ping
// @Success 200 {object} OK
// @Failure 403 :没有该服务
// @router / [get]
func (c *PingController) Ping() {
	var appname = beego.AppConfig.String("appname")
	c.Data["json"] = JSONStruct{"success", 0, appname, "获取成功"}
	c.ServeJSON()
}
