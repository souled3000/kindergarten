package controllers

import (
	//	"fmt"
	//	"kindergarten-service-go/models"

	"github.com/astaxie/beego"
)

//基类
type BaseController struct {
	beego.Controller
}

//func (c *BaseController) Prepare() {
//	u_id := 2726
//	//	u_id, _ := c.GetInt("u_id")
//	Url, err := models.GetPermissionRoute(u_id)
//	fmt.Println(Url)
//	controller, action := c.GetControllerAndAction()
//	route := "" + controller + "." + action + ""
//	if err == nil {
//		tmp := false
//		for _, v := range Url {
//			if v["route"] != nil {
//				if v["route"].(string) == route {
//					tmp = true
//					break
//				}
//			}
//		}
//		if !tmp {
//			c.Data["json"] = JSONStruct{"error", 1006, nil, "没有权限"}
//			c.ServeJSON()
//		}
//	}
//}
