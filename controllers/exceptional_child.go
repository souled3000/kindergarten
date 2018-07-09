package controllers

import (
	"kindergarten-service-go/models"
	"strconv"

	"github.com/astaxie/beego/validation"
)

type ExceptionalChildController struct {
	BaseController
}

func (c *ExceptionalChildController) URLMapping() {
	c.Mapping("Get", c.Get)
}

// @Title 特殊儿童列表/搜索特殊儿童
// @Description 特殊儿童列表/搜索特殊儿童
// @Param	keyword	   	 query	string	 false	"关键字(特殊儿童姓名/特殊儿童过敏源)"
// @Param	page	  	 query	int	 	 false	"当前页"
// @Success 0 			 {object}  models.ExceptionalChild
// @Failure 1005 获取失败
// @router / [get]
func (c *ExceptionalChildController) GetSearch() {
	// 关键字
	keyword := c.GetString("keyword")
	// page_num
	page, _ := c.GetInt64("page")

	// limit
	limit, _ := c.GetInt64("per_page")

	if info, err := models.GetAllExceptionalChild("", 0, page, limit, keyword); err == nil {
		c.Data["json"] = JSONStruct{"success", 0, info, "获取成功"}
	} else {
		c.Data["json"] = JSONStruct{"error", 1005, nil, "获取失败"}
	}
	c.ServeJSON()
}

// GetAll ...
// @Title 根据过敏源获取特殊儿童
// @Description 根据过敏源获取特殊儿童
// @Param	allergen	query	string	true	"过敏源信息，多个过敏源以','分隔"
// @Success 0 			{json} 	JSONStruct
// @Failure 1005 获取失败
// @router / [get]
func (c *ExceptionalChildController) GetAllergenChild() {
	allergen := c.GetString("allergen")
	valid := validation.Validation{}
	valid.Required(allergen, "allergen").Message("过敏源信息不能为空")

	if valid.HasErrors() {
		c.Data["json"] = JSONStruct{"error", 1001, nil, valid.Errors[0].Message}
	} else {
		if allergenChild, err := models.GetAllergenChild(allergen); err == nil {
			c.Data["json"] = JSONStruct{"success", 0, allergenChild, "获取成功"}
		} else {
			c.Data["json"] = JSONStruct{"error", 1005, nil, "获取失败"}
		}
	}
	c.ServeJSON()
}

// @Title 过敏食物报备
// @Description 过敏食物报备
// @param 		class				query  	int    	true		"班级ID"
// @param 		kindergarten_id		query  	int    	true		"幼儿园ID"
// @param 		creator				query  	int    	true		"创建人ID"
// @param 		student_id			query  	int    	true		"宝宝ID"
// @param 		source				query  	int    	true		"来源信息"
// @param 		somatotype			query  	int    	true		"体质类型"
// @param		child_name			query	string	true		"特殊儿童姓名"
// @param 		allergen			query  	string 	true		"过敏源，多个过敏源以','分隔"
// @router / [post]
func (c *ExceptionalChildController) AllergenPreparation() {
	// 幼儿园ID
	kindergarten_id, _ := c.GetInt("kindergarten_id")
	// 创建人ID
	creator, _ := c.GetInt("creator")
	// 宝宝ID
	baby_id, _ := c.GetInt("baby_id")
	// 来源信息
	source, _ := c.GetInt8("source")
	// 体质类型
	somatotype, _ := c.GetInt8("somatotype")
	// 特殊儿童姓名
	child_name := c.GetString("child_name")
	// 过敏源
	allergen := c.GetString("allergen")

	valid := validation.Validation{}
	valid.Required(kindergarten_id, "kindergarten_id").Message("幼儿园ID不能为空")
	valid.Required(creator, "creator").Message("创建人ID不能为空")
	valid.Required(baby_id, "baby_id").Message("宝宝ID不能为空")
	valid.Required(source, "source").Message("来源信息不能为空")
	valid.Required(somatotype, "somatotype").Message("体质类型不能为空")
	valid.Required(allergen, "allergen").Message("过敏源不能为空")

	if valid.HasErrors() {
		c.Data["json"] = JSONStruct{"error", 1001, nil, valid.Errors[0].Message}
	} else {
		if _, err := models.AllergenPreparation(child_name, somatotype, allergen, source, kindergarten_id, creator, baby_id); err == nil {
			c.Data["json"] = JSONStruct{"success", 0, err, "保存成功"}
		} else {
			c.Data["json"] = JSONStruct{"error", 1003, nil, "保存失败"}
		}
	}
	c.ServeJSON()
}

// GetAllergen ...
// @Title 			根据宝宝ID过敏源
// @Description 		根据宝宝ID过敏源
// @Param	id		query 	int		true		"宝宝IDxs"
// @Success 0 		{json} 	JSONStruct
// @Failure 1004		删除失败
func (c *ExceptionalChildController) GetAllergen() {
	// 宝宝ID
	id, _ := c.GetInt("baby_id")
	if info, err := models.GetAllergen(id); err == nil {
		c.Data["json"] = JSONStruct{"success", 0, info, "获取成功"}
	} else {
		c.Data["json"] = JSONStruct{"error", 1005, nil, "获取失败"}
	}
	c.ServeJSON()
}

// DelAllergen ...
// @Title 			删除过敏源
// @Description 		删除过敏源
// @Param	id		path 	string	true		"特殊儿童ID"
// @Success 0 		{json} 	JSONStruct
// @Failure 1004		删除失败
// @router /:id [delete]
func (c *ExceptionalChildController) DelAllergen() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	if err := models.DeleteExceptionalChild(id); err == nil {
		c.Data["json"] = JSONStruct{"success", 0, err, "删除成功"}
	} else {
		c.Data["json"] = JSONStruct{"error", 1004, nil, "删除失败"}
	}
	c.ServeJSON()
}
