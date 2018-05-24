package controllers

import (
	"kindergarten-service-go/models"

	"github.com/astaxie/beego/validation"

	"github.com/astaxie/beego"
)

//组织架构
type OrganizationalController struct {
	beego.Controller
}

// URLMapping ...
func (c *OrganizationalController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
}

// GetClass ...
// @Title 班级
// @Description 班级列表
// @Param	kindergarten_id           query	int	     true		"幼儿园ID"
// @Param	class_type                query	int	     false		"班级类型"
// @Param	page                      query	int	     false		"页数"
// @Param	per_page                  query	int	     false		"每页显示条数"
// @Success 200 {object} models.Organizational
// @Failure 403
// @router /class [get]
func (o *OrganizationalController) GetClass() {
	var prepage int = 20
	var page int
	if v, err := o.GetInt("per_page"); err == nil {
		prepage = v
	}
	if v, err := o.GetInt("page"); err == nil {
		page = v
	}
	kindergarten_id, _ := o.GetInt("kindergarten_id")
	class_type, _ := o.GetInt("class_type")
	valid := validation.Validation{}
	valid.Required(kindergarten_id, "kindergarten_id").Message("幼儿园编号不能为空")
	valid.Required(class_type, "class_type").Message("班级类型不能为空")
	if valid.HasErrors() {
		o.Data["json"] = JSONStruct{"error", 1001, nil, valid.Errors[0].Message}
		o.ServeJSON()
	} else {
		v := models.GetClassAll(kindergarten_id, class_type, page, prepage)
		if v == nil {
			o.Data["json"] = JSONStruct{"error", 1005, nil, "获取失败"}
		} else {
			o.Data["json"] = JSONStruct{"success", 0, v, "获取成功"}
		}
		o.ServeJSON()
	}

}

// Member ...
// @Title 班级成员
// @Description 班级成员
// @Param	kindergarten_id           query	int	     true		"幼儿园ID"
// @Param	class_type                query	int	     false		"班级类型"
// @Param	class_id                  query	int	     false		"班级id"
// @Param	page                      query	int	     false		"页数"
// @Param	per_page                  query	int	     false		"每页显示条数"
// @Success 200 {object} models.Organizational
// @Failure 403
// @router /member [get]
func (o *OrganizationalController) Member() {
	var prepage int = 20
	var page int
	if v, err := o.GetInt("per_page"); err == nil {
		prepage = v
	}
	if v, err := o.GetInt("page"); err == nil {
		page = v
	}
	kindergarten_id, _ := o.GetInt("kindergarten_id")
	class_type, _ := o.GetInt("class_type")
	class_id, _ := o.GetInt("class_id")
	valid := validation.Validation{}
	valid.Required(kindergarten_id, "kindergarten_id").Message("幼儿园编号不能为空")
	valid.Required(class_type, "class_type").Message("班级类型不能为空")
	valid.Required(class_id, "class_id").Message("班级id不能为空")
	if valid.HasErrors() {
		o.Data["json"] = JSONStruct{"error", 1001, nil, valid.Errors[0].Message}
		o.ServeJSON()
	} else {
		v := models.ClassMember(kindergarten_id, class_type, class_id, page, prepage)
		if v == nil {
			o.Data["json"] = JSONStruct{"error", 1005, nil, "获取失败"}
		} else {
			o.Data["json"] = JSONStruct{"success", 0, v, "获取成功"}
		}
		o.ServeJSON()
	}
}

// Destroy ...
// @Title 删除班级
// @Description 删除班级
// @Param	class_id                  query	int	     false		"班级id"
// @Success 200 {object} models.Organizational
// @Failure 403
// @router / [delete]
func (o *OrganizationalController) Destroy() {
	class_id, _ := o.GetInt("class_id")
	valid := validation.Validation{}
	valid.Required(class_id, "class_id").Message("班级id不能为空")
	if valid.HasErrors() {
		o.Data["json"] = JSONStruct{"error", 1001, nil, valid.Errors[0].Message}
		o.ServeJSON()
	} else {
		v := models.Destroy(class_id)
		if v == nil {
			o.Data["json"] = JSONStruct{"error", 1003, nil, "删除失败"}
		} else {
			o.Data["json"] = JSONStruct{"success", 0, v, "删除成功"}
		}
		o.ServeJSON()
	}
}
