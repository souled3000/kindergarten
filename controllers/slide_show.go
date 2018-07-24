package controllers

import (
	"kindergarten-service-go/models"
	"log"
	"strconv"

	"github.com/astaxie/beego/validation"
)

//轮播图
type SideShowController struct {
	BaseController
}

// Store ...
// @Title 保存轮播图
// @Description 保存轮播图
// @Param	title		        query    string   		    true		"标题"
// @Param	content		        query    string 	    	true		"内容"
// @Param	picture		        query    string 	    	true		"图片"
// @Param	kindergarten_id		query    int 	    	    true		"幼儿园ID"
// @Success 201 {int} models.SideShow
// @Failure 403 body is empty
// @router / [post]
func (c *SideShowController) Store() {
	title := c.GetString("title")
	content := c.GetString("content")
	picture := c.GetString("picture")
	kindergarten_id, _ := c.GetInt("kindergarten_id")
	valid := validation.Validation{}
	valid.Required(title, "title").Message("标题不能为空")
	valid.Required(content, "content").Message("内容不能为空")
	valid.Required(picture, "picture").Message("图片不能为空")
	valid.Required(kindergarten_id, "kindergarten_id").Message("幼儿园编号不能为空")
	if valid.HasErrors() {
		log.Println(valid.Errors)
		c.Data["json"] = JSONStruct{"error", 1001, nil, valid.Errors[0].Message}
		c.ServeJSON()
	} else {
		err := models.AddSlideShow(title, content, kindergarten_id, picture)
		if err != nil {
			c.Data["json"] = JSONStruct{"error", 1003, err.Error(), "保存失败"}
		} else {
			c.Data["json"] = JSONStruct{"success", 0, nil, "保存成功"}
		}
		c.ServeJSON()
	}
}

// GetSideShow ...
// @Title 轮播图详情
// @Description 轮播图详情
// @Param	id       query	int	 true		"主键ID"
// @Success 200 {object} models.SideShow
// @Failure 403 :编号为空
// @router /:id [get]
func (c *SideShowController) GetSideShow() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	v, err := models.GetSlideShow(id)
	if err != nil {
		c.Data["json"] = JSONStruct{"error", 1005, err.Error(), "获取失败"}
	} else {
		c.Data["json"] = JSONStruct{"success", 0, v, "获取成功"}
	}
	c.ServeJSON()
}

// GetSideShowList ...
// @Title 轮播图列表
// @Description 轮播图列表
// @Param	page     query	int	 false		"页数"
// @Param	per_page query	int	 false		"每页显示条数"
// @Param	kindergarten_id query	int	 false		"幼儿园id"
// @Success 200 {object} models.SideShow
// @Failure 403
// @router / [get]
func (c *SideShowController) GetSideShowList() {
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
		v, err := models.GetSlideShowList(page, prepage, kindergarten_id)
		if err != nil {
			c.Data["json"] = JSONStruct{"error", 1005, err.Error(), "获取失败"}
		} else {
			c.Data["json"] = JSONStruct{"success", 0, v, "获取成功"}
		}
		c.ServeJSON()
	}
}

// DeleteSideShow ...
// @Title 删除轮播图
// @Description 删除轮播图
// @Param	id		path 	string	true		"自增id"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (c *SideShowController) DeleteSideShow() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	err := models.DeleteSlideShow(id)
	if err != nil {
		c.Data["json"] = JSONStruct{"error", 1004, nil, "删除失败"}
	} else {
		c.Data["json"] = JSONStruct{"success", 0, nil, "删除成功"}
	}
	c.ServeJSON()
}

// Put ...
// @Title 编辑轮播图
// @Description 编辑轮播图
// @Param	id		             path 	    int	            true		"编号"
// @Param	title		        query    string   		true		"标题"
// @Param	content		        query    string 	    	true		"内容"
// @Param	picture		        query    string 	    	true		"图片"
// @Param	kindergarten_id		query    int 	    	true		"幼儿园ID"
// @Success 200 {object} models.SideShow
// @Failure 403 :id is not int
// @router /:id [put]
func (c *SideShowController) Put() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	title := c.GetString("title")
	content := c.GetString("content")
	picture := c.GetString("picture")
	kindergarten_id, _ := c.GetInt("kindergarten_id")
	valid := validation.Validation{}
	valid.Required(title, "title").Message("标题不能为空")
	valid.Required(content, "content").Message("内容不能为空")
	valid.Required(picture, "picture").Message("图片不能为空")
	valid.Required(kindergarten_id, "kindergarten_id").Message("幼儿园编号不能为空")
	if valid.HasErrors() {
		log.Println(valid.Errors)
		c.Data["json"] = JSONStruct{"error", 1001, nil, valid.Errors[0].Message}
		c.ServeJSON()
	} else {
		err := models.UpdateSlideShow(id, title, content, kindergarten_id, picture)
		if err != nil {
			c.Data["json"] = JSONStruct{"error", 1003, err.Error(), "编辑失败"}
		} else {
			c.Data["json"] = JSONStruct{"success", 0, nil, "编辑成功"}
		}
		c.ServeJSON()
	}
}
