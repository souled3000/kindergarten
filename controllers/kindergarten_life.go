package controllers

import (
	"kindergarten-service-go/models"
	"log"
	"strconv"

	"github.com/astaxie/beego/validation"
)

//园内生活
type KindergartenLifeController struct {
	BaseController
}

// Store ...
// @Title 保存园内生活
// @Description Web-保存园内生活
// @Param	content		        query    string   		true		"标题"
// @Param	template		    query    int 	    	true		"公告内容"
// @Param	kindergarten_id		query    int 	    	true		"幼儿园ID"
// @Success 201 {int} models.KindergartenLife
// @Failure 403 body is empty
// @router / [post]
func (c *KindergartenLifeController) Store() {
	number, _ := c.GetInt("number")
	content := c.GetString("content")
	template, _ := c.GetInt("template")
	picture := c.GetString("picture")
	kindergarten_id, _ := c.GetInt("kindergarten_id")
	valid := validation.Validation{}
	valid.Required(picture, "picture").Message("图片不能为空")
	valid.Required(content, "content").Message("内容不能为空")
	valid.Required(template, "template").Message("模板不能为空")
	valid.Required(kindergarten_id, "kindergarten_id").Message("幼儿园编号不能为空")
	if valid.HasErrors() {
		log.Println(valid.Errors)
		c.Data["json"] = JSONStruct{"error", 1001, nil, valid.Errors[0].Message}
		c.ServeJSON()
	} else {
		err := models.AddKindergartenLife(content, template, kindergarten_id, picture, number)
		if err != nil {
			c.Data["json"] = JSONStruct{"error", 1003, err.Error(), "保存失败"}
		} else {
			c.Data["json"] = JSONStruct{"success", 0, nil, "保存成功"}
		}
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
	v, err := models.GetKindergartenLifeInfo(id)
	if err != nil {
		c.Data["json"] = JSONStruct{"error", 1005, err.Error(), "获取失败"}
	} else {
		c.Data["json"] = JSONStruct{"success", 0, v, "获取成功"}
	}
	c.ServeJSON()
}

// GetKindergartenLifeList ...
// @Title 园内生活列表
// @Description Web-园内生活列表
// @Param	page     query	int	 false		"页数"
// @Param	per_page query	int	 false		"每页显示条数"
// @Param	kindergarten_id query	int	 false		"幼儿园id"
// @Success 200 {object} models.KindergartenLife
// @Failure 403
// @router / [get]
func (c *KindergartenLifeController) GetKindergartenLifeList() {
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
		v, err := models.GetKindergartenLifeList(page, prepage, kindergarten_id)
		if err != nil {
			c.Data["json"] = JSONStruct{"error", 1005, err.Error(), "获取失败"}
		} else {
			c.Data["json"] = JSONStruct{"success", 0, v, "获取成功"}
		}
		c.ServeJSON()
	}
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
	err := models.DeleteKindergartenLife(id)
	if err != nil {
		c.Data["json"] = JSONStruct{"error", 1004, nil, "删除失败"}
	} else {
		c.Data["json"] = JSONStruct{"success", 0, nil, "删除成功"}
	}
	c.ServeJSON()
}

// Put ...
// @Title 编辑园内生活
// @Description 编辑园内生活
// @Param	id		             path 	    int	            true		"编号"
// @Param	content		         query    string   		true		"标题"
// @Param	template		     query    int 	    	true		"公告内容"
// @Param	kindergarten_id		 query    int 	    	true		"幼儿园ID"
// @Success 200 {object} models.KindergartenLife
// @Failure 403 :id is not int
// @router /:id [put]
func (c *KindergartenLifeController) Put() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	content := c.GetString("content")
	template := c.GetString("template")
	kindergarten_id, _ := c.GetInt("kindergarten_id")
	err := models.UpdateKL(id, content, template, kindergarten_id)
	if err != nil {
		c.Data["json"] = JSONStruct{"error", 1003, err.Error(), "编辑失败"}
	} else {
		c.Data["json"] = JSONStruct{"success", 0, nil, "编辑成功"}
	}
	c.ServeJSON()
}
