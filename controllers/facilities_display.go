package controllers

import (
	"kindergarten-service-go/models"
	"log"
	"strconv"

	"github.com/astaxie/beego/validation"
)

//设施展示
type FacilitiesDisplayController struct {
	BaseController
}

// Store ...
// @Title 保存设施
// @Description 保存设施
// @Param	picture		     query      string   		    true		"图片"
// @Param	order		     query      int 	    	    true		"排序"
// @Param	kindergarten_id	 query	    int 	    	    true		"幼儿园ID"
// @Success 201 {int} models.FacilitiesDisplay
// @Failure 403 body is empty
// @router / [post]
func (c *FacilitiesDisplayController) Store() {
	order, _ := c.GetInt("order")
	picture := c.GetString("picture")
	kindergarten_id, _ := c.GetInt("kindergarten_id")
	valid := validation.Validation{}
	valid.Required(order, "order").Message("排序不能为空")
	valid.Required(picture, "picture").Message("图片不能为空")
	valid.Required(kindergarten_id, "kindergarten_id").Message("幼儿园编号不能为空")
	if valid.HasErrors() {
		log.Println(valid.Errors)
		c.Data["json"] = JSONStruct{"error", 1001, nil, valid.Errors[0].Message}
		c.ServeJSON()
	} else {
		err := models.Store(order, picture, kindergarten_id)
		if err != nil {
			c.Data["json"] = JSONStruct{"error", 1003, err.Error(), "保存失败"}
		} else {
			c.Data["json"] = JSONStruct{"success", 0, nil, "保存成功"}
		}
		c.ServeJSON()
	}
}

// GetNoticeList ...
// @Title 设施列表
// @Description 设施列表
// @Param	page            query	int	 false		"页数"
// @Param	per_page        query	int	 false		"每页显示条数"
// @Param	kindergarten_id query	int	 false		"幼儿园id"
// @Success 200 {object} models.FacilitiesDisplay
// @Failure 403
// @router / [get]
func (c *FacilitiesDisplayController) GetNoticeList() {
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
		v, err := models.GetList(page, prepage, kindergarten_id)
		if err != nil {
			c.Data["json"] = JSONStruct{"error", 1005, err.Error(), "获取失败"}
		} else {
			c.Data["json"] = JSONStruct{"success", 0, v, "获取成功"}
		}
		c.ServeJSON()
	}
}

// GetOne ...
// @Title 设施详情
// @Description 设施详情
// @Param	id       query	int	 true		"主键ID"
// @Success 200 {object} models.FacilitiesDisplay
// @Failure 403 :编号为空
// @router /:id [get]
func (c *FacilitiesDisplayController) GetOne() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	v, err := models.GetOne(id)
	if err != nil {
		c.Data["json"] = JSONStruct{"error", 1005, err.Error(), "获取失败"}
	} else {
		c.Data["json"] = JSONStruct{"success", 0, v, "获取成功"}
	}
	c.ServeJSON()
}

// Delete ...
// @Title 删除设施
// @Description 删除设施
// @Param	id		path 	string	true		"用户编号"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (c *FacilitiesDisplayController) Delete() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	err := models.Delete(id)
	if err != nil {
		c.Data["json"] = JSONStruct{"error", 1004, err.Error(), "删除失败"}
	} else {
		c.Data["json"] = JSONStruct{"success", 0, nil, "删除成功"}
	}
	c.ServeJSON()
}

// Put ...
// @Title 编辑设施
// @Description 编辑设施
// @Param	id		             path 	    int	            true		"编号"
// @Param	order		         query      int   		true		"排序"
// @Param	picture		         query      string 	    	true		"图片"
// @Param	kindergarten_id	     query	   int 	    	    true		"幼儿园ID"
// @Success 200 {object} models.FacilitiesDisplay
// @Failure 403 :id is not int
// @router /:id [put]
func (c *FacilitiesDisplayController) Put() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	picture := c.GetString("picture")
	order, _ := c.GetInt("order")
	kindergarten_id, _ := c.GetInt("kindergarten_id")
	err := models.Update(id, picture, order, kindergarten_id)
	if err != nil {
		c.Data["json"] = JSONStruct{"error", 1003, err.Error(), "编辑失败"}
	} else {
		c.Data["json"] = JSONStruct{"success", 0, nil, "编辑成功"}
	}
	c.ServeJSON()
}
