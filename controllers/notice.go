package controllers

import (
	"kindergarten-service-go/models"
	"log"
	"strconv"

	"github.com/astaxie/beego/validation"
)

//发布公告
type NoticeController struct {
	BaseController
}

type JSONStruct struct {
	Status string      `json:"status";`
	Code   int         `json:"code";`
	Result interface{} `json:"result";`
	Msg    string      `json:"msg";`
}

// Store ...
// @Title 保存公告
// @Description Web-保存公告
// @Param	title		     query      string   		true		"标题"
// @Param	content		     query      string 	    	true		"公告内容"
// @Param	kindergarten_id	 query	   int 	    	    true		"幼儿园ID"
// @Success 201 {int} models.Notice
// @Failure 403 body is empty
// @router / [post]
func (c *NoticeController) Store() {
	title := c.GetString("title")
	content := c.GetString("content")
	kindergarten_id, _ := c.GetInt("kindergarten_id")
	valid := validation.Validation{}
	valid.Required(title, "title").Message("标题不能为空")
	valid.Required(content, "content").Message("发布内容不能为空")
	valid.Required(kindergarten_id, "kindergarten_id").Message("幼儿园编号不能为空")
	if valid.HasErrors() {
		log.Println(valid.Errors)
		c.Data["json"] = JSONStruct{"error", 1001, nil, valid.Errors[0].Message}
		c.ServeJSON()
	} else {
		err := models.AddNotice(title, content, kindergarten_id)
		if err != nil {
			c.Data["json"] = JSONStruct{"error", 1003, err.Error(), "保存失败"}
		} else {
			c.Data["json"] = JSONStruct{"success", 0, nil, "保存成功"}
		}
		c.ServeJSON()
	}
}

// GetNoticeList ...
// @Title 公告列表
// @Description 公告列表
// @Param	page     query	int	 false		"页数"
// @Param	per_page query	int	 false		"每页显示条数"
// @Success 200 {object} models.Notice
// @Failure 403
// @router / [get]
func (c *NoticeController) GetNoticeList() {
	prepage, _ := c.GetInt("per_page", 20)
	page, _ := c.GetInt("page")
	kindergarten_id, _ := c.GetInt("kindergarten_id")
	valid := validation.Validation{}
	valid.Required(kindergarten_id, "kindergarten_id").Message("幼儿园编号不能为空")
	if valid.HasErrors() {
		log.Println(valid.Errors)
		c.Data["json"] = JSONStruct{"error", 1001, nil, valid.Errors[0].Message}
		c.ServeJSON()
	} else {
		v, err := models.GetNoticeList(page, prepage, kindergarten_id)
		if err != nil {
			c.Data["json"] = JSONStruct{"error", 1005, err.Error(), "获取失败"}
		} else {
			c.Data["json"] = JSONStruct{"success", 0, v, "获取成功"}
		}
		c.ServeJSON()
	}
}

// GetNoticeInfo ...
// @Title 公告详情
// @Description Web-公告详情
// @Param	id       query	int	 true		"主键ID"
// @Param	page     query	int	 false		"页数"
// @Param	per_page query	int	 false		"每页显示条数"
// @Success 200 {object} models.Notice
// @Failure 403 :编号为空
// @router /:id [get]
func (c *NoticeController) GetNoticeInfo() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	v, err := models.GetNoticeInfo(id)
	if err != nil {
		c.Data["json"] = JSONStruct{"error", 1005, err.Error(), "获取失败"}
	} else {
		c.Data["json"] = JSONStruct{"success", 0, v, "获取成功"}
	}
	c.ServeJSON()
}

// Delete ...
// @Title 删除公告
// @Description Web-删除公告
// @Param	id		path 	string	true		"用户编号"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (c *NoticeController) Delete() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	err := models.DeleteNotice(id)
	if err != nil {
		c.Data["json"] = JSONStruct{"error", 1004, err.Error(), "删除失败"}
	} else {
		c.Data["json"] = JSONStruct{"success", 0, nil, "删除成功"}
	}
	c.ServeJSON()
}

// Put ...
// @Title 编辑公告
// @Description 编辑公告
// @Param	id		             path 	    int	            true		"编号"
// @Param	title		         query      string   		true		"标题"
// @Param	content		         query      string 	    	true		"公告内容"
// @Param	kindergarten_id	     query	   int 	    	    true		"幼儿园ID"
// @Success 200 {object} models.Notice
// @Failure 403 :id is not int
// @router /:id [put]
func (c *NoticeController) Put() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	title := c.GetString("title")
	content := c.GetString("content")
	kindergarten_id, _ := c.GetInt("kindergarten_id")
	err := models.UpdateNotice(id, title, content, kindergarten_id)
	if err != nil {
		c.Data["json"] = JSONStruct{"error", 1003, err.Error(), "编辑失败"}
	} else {
		c.Data["json"] = JSONStruct{"success", 0, nil, "编辑成功"}
	}
	c.ServeJSON()
}
