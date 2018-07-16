package app

import (
	"github.com/astaxie/beego"
)

type BaseController struct {
	beego.Controller
}

type JSONStruct struct {
	Status string      `json:"status"`
	Code   int         `json:"code"`
	Result interface{} `json:"result"`
	Msg    string      `json:"msg"`
}

func (c *BaseController) StopJson() {
	c.ServeJSON()
	c.Finish()
	c.StopRun()
}
