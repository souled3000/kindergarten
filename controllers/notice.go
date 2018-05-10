package controllers

import (
	"encoding/json"
	"kindergarten-service-go/models"
	"log"
	"strconv"

	"github.com/astaxie/beego/validation"

	"github.com/astaxie/beego"
)

// NoticeController operations for Notice
type NoticeController struct {
	beego.Controller
}

type JSONStruct struct {
	Status string      `json:"status";`
	Code   int         `json:"code";`
	Result interface{} `json:"result";`
	Msg    string      `json:"msg";`
}

// URLMapping ...
func (c *NoticeController) URLMapping() {
	c.Mapping("Store", c.Store)
	c.Mapping("GetNoticeList", c.GetNoticeList)
	c.Mapping("Delete", c.Delete)
}

// Store ...
// @Title Store
// @Description create Notice
// @Param	body		body 	models.Notice	true		"array"
// @Success 201 {int} models.Notice
// @Failure 403 body is empty
// @router / [post]
func (c *NoticeController) Store() {
	var v models.Notice
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		valid := validation.Validation{}
		valid.Required(v.Title, "Title").Message("标题不能为空")
		valid.Required(v.Content, "Content").Message("发布内容不能为空")
		valid.Required(v.KindergartenId, "KindergartenId").Message("幼儿园编号不能为空")
		if valid.HasErrors() {
			log.Println(valid.Errors)
			c.Data["json"] = JSONStruct{"error", 1001, nil, valid.Errors[0].Message}
			c.ServeJSON()
		} else {
			v := models.AddNotice(&v)
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

// GetNoticeList ...
// @Title Get Notice List
// @Description get Notice
// @Param	page     query	int	 false		"页数"
// @Param	per_page query	int	 false		"每页显示条数"
// @Success 200 {object} models.Notice
// @Failure 403
// @router / [get]
func (c *NoticeController) GetNoticeList() {
	var prepage int = 20
	var page int
	if v, err := c.GetInt("per_page"); err == nil {
		prepage = v
	}
	if v, err := c.GetInt("page"); err == nil {
		page = v
	}
	v := models.GetNoticeList(page, prepage)
	if v == nil {
		c.Data["json"] = JSONStruct{"error", 1005, v, "获取失败"}
	} else {
		c.Data["json"] = JSONStruct{"success", 0, v, "获取成功"}
	}
	c.ServeJSON()
}

// Delete ...
// @Title Delete
// @Description delete the Notice
// @Param	id		path 	string	true		"用户编号"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (c *NoticeController) Delete() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	v := models.DeleteNotice(id)
	if v == nil {
		c.Data["json"] = JSONStruct{"error", 1004, nil, "删除失败"}
	} else {
		c.Data["json"] = JSONStruct{"success", 0, nil, "删除成功"}
	}
	c.ServeJSON()
}
