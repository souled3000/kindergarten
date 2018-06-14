package admin

import (
	"encoding/json"
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

// URLMapping ...
func (c *NoticeController) URLMapping() {
	c.Mapping("Store", c.Store)
	c.Mapping("GetNoticeList", c.GetNoticeList)
	c.Mapping("Delete", c.Delete)
}

// Store ...
// @Title 保存公告
// @Description Web-保存公告
// @Param	title		        string   	boby	true		"标题"
// @Param	content		        string 	    boby	true		"公告内容"
// @Param	kindergarten_id		int 	    boby	true		"幼儿园ID"
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
// @Title 公告列表
// @Description Web-公告列表
// @Param	page     query	int	 false		"页数"
// @Param	per_page query	int	 false		"每页显示条数"
// @Success 200 {object} models.Notice
// @Failure 403
// @router / [get]
func (c *NoticeController) GetNoticeList() {
	prepage, _ := c.GetInt("per_page", 20)
	page, _ := c.GetInt("page")
	v := models.GetNoticeList(page, prepage)
	if v == nil {
		c.Data["json"] = JSONStruct{"error", 1005, v, "获取失败"}
	} else {
		c.Data["json"] = JSONStruct{"success", 0, v, "获取成功"}
	}
	c.ServeJSON()
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
	v := models.GetNoticeInfo(id)
	if v == nil {
		c.Data["json"] = JSONStruct{"error", 1005, v, "获取失败"}
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
	v := models.DeleteNotice(id)
	if v == nil {
		c.Data["json"] = JSONStruct{"error", 1004, nil, "删除失败"}
	} else {
		c.Data["json"] = JSONStruct{"success", 0, nil, "删除成功"}
	}
	c.ServeJSON()
}
