package controllers

import (
	"encoding/json"
	"kindergarten-service-go/models"
	"log"
	"strconv"

	"github.com/astaxie/beego/validation"

	"github.com/astaxie/beego"
)

//园内生活
type KindergartenLifeController struct {
	beego.Controller
}

// Store ...
// @Title 保存园内生活
// @Description Web-保存园内生活
// @Param	content		            string   	boby	true		"标题"
// @Param	template		        int 	    boby	true		"公告内容"
// @Param	kindergarten_id		    int 	    boby	true		"幼儿园ID"
// @Success 201 {int} models.KindergartenLife
// @Failure 403 body is empty
// @router / [post]
func (c *KindergartenLifeController) Store() {
	var v models.KindergartenLife
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		valid := validation.Validation{}
		valid.Required(v.Content, "Content").Message("内容不能为空")
		valid.Required(v.KindergartenId, "KindergartenId").Message("幼儿园编号不能为空")
		valid.Required(v.Template, "Template").Message("模板不能为空")
		if valid.HasErrors() {
			log.Println(valid.Errors)
			c.Data["json"] = JSONStruct{"error", 1001, nil, valid.Errors[0].Message}
			c.ServeJSON()
		} else {
			v := models.AddKindergartenLife(&v)
			if v == nil {
				c.Data["json"] = JSONStruct{"error", 1003, err.Error(), "保存失败"}
			} else {
				c.Data["json"] = JSONStruct{"success", 0, nil, "保存成功"}
			}
			c.ServeJSON()
		}
	} else {
		c.Data["json"] = JSONStruct{"error", 1001, err.Error(), "字段必须为json格式"}
		c.ServeJSON()
	}
}

// GetKindergartenLifeInfo ...
// @Title 园内生活详情
// @Description Web-园内生活详情
// @Param	id       query	int	 true		"主键ID"
// @Param	page     query	int	 false		"页数"
// @Param	per_page query	int	 false		"每页显示条数"
// @Success 200 {object} models.KindergartenLife
// @Failure 403 :编号为空
// @router /:id [get]
func (c *KindergartenLifeController) GetKindergartenLifeInfo() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	v := models.GetKindergartenLifeInfo(id)
	if v == nil {
		c.Data["json"] = JSONStruct{"error", 1005, v, "获取失败"}
	} else {
		c.Data["json"] = JSONStruct{"success", 0, v, "获取成功"}
	}
	c.ServeJSON()
}

// GetKindergartenLifeList ...
// @Title Web-园内生活列表
// @Description Web-园内生活列表
// @Param	page     query	int	 false		"页数"
// @Param	per_page query	int	 false		"每页显示条数"
// @Success 200 {object} models.KindergartenLife
// @Failure 403
// @router / [get]
func (c *KindergartenLifeController) GetKindergartenLifeList() {
	prepage, _ := c.GetInt("per_page", 20)
	page, _ := c.GetInt("page")
	v := models.GetKindergartenLifeList(page, prepage)
	if v == nil {
		c.Data["json"] = JSONStruct{"error", 1005, v, "获取失败"}
	} else {
		c.Data["json"] = JSONStruct{"success", 0, v, "获取成功"}
	}
	c.ServeJSON()
}

// Delete ...
// @Title Web-删除园内生活
// @Description Web-删除园内生活
// @Param	id		path 	string	true		"自增id"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (c *KindergartenLifeController) Delete() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	v := models.DeleteKindergartenLife(id)
	if v == nil {
		c.Data["json"] = JSONStruct{"error", 1004, nil, "删除失败"}
	} else {
		c.Data["json"] = JSONStruct{"success", 0, nil, "删除成功"}
	}
	c.ServeJSON()
}
