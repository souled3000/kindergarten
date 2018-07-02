package controllers

import (
	"kindergarten-service-go/models"
	"github.com/astaxie/beego/validation"
)

type AllergenChildController struct {
	BaseController
}



// GetAll ...
// @Title 根据过敏源获取特殊儿童
// @Description 根据过敏源获取特殊儿童
// @Param	allergen	query	string	true	"过敏源信息，多个过敏源以','分隔"
// @Success 0 			{json} 	JSONStruct
// @Failure 1005 获取失败
// @router / [get]
func (c *AllergenChildController) GetAllergenChild() {
	allergen := c.GetString("allergen")
	valid := validation.Validation{}
	valid.Required(allergen, "allergen").Message("过敏源信息不能为空")

	if valid.HasErrors() {
		c.Data["json"] = JSONStruct{"error", 1001, nil, valid.Errors[0].Message}
	} else {
		if allergenChild, err := models.GetAllergenChild(allergen); err == nil {
			c.Data["json"] = JSONStruct{"success", 0, allergenChild, "获取成功"}
		} else {
			c.Data["json"] = JSONStruct{"error", 1005, err, "获取失败"}
		}
	}
	c.ServeJSON()
}
