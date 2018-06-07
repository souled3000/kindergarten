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
	prepage, _ := o.GetInt("per_page", 20)
	page, _ := o.GetInt("page")
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
	prepage, _ := o.GetInt("per_page", 20)
	page, _ := o.GetInt("page")
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
// @Param	class_id        query	int	     false		"班级id"
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
		_, err := models.Destroy(class_id)
		if err != nil {
			o.Data["json"] = JSONStruct{"error", 1003, nil, err.Error()}
		} else {
			o.Data["json"] = JSONStruct{"success", 0, nil, "删除成功"}
		}
		o.ServeJSON()
	}
}

// Store ...
// @Title 创建班级
// @Description 创建班级
// @Param	kindergarten_id      query	int	     true		"幼儿园id"
// @Param	class_type           query	int	     true		"班级类型"
// @Success 200 {object} models.Organizational
// @Failure 403
// @router / [post]
func (o *OrganizationalController) Store() {
	class_type, _ := o.GetInt("class_type")
	kindergarten_id, _ := o.GetInt("kindergarten_id")
	valid := validation.Validation{}
	valid.Required(class_type, "class_type").Message("班级类型不能为空")
	valid.Required(kindergarten_id, "kindergarten_id").Message("幼儿园ID不能为空")
	if valid.HasErrors() {
		o.Data["json"] = JSONStruct{"error", 1001, nil, valid.Errors[0].Message}
		o.ServeJSON()
	} else {
		v, err := models.StoreClass(class_type, kindergarten_id)
		if err != nil {
			o.Data["json"] = JSONStruct{"error", 1003, nil, err.Error()}
		} else {
			o.Data["json"] = JSONStruct{"success", 0, v, "保存成功"}
		}
		o.ServeJSON()
	}
}

// GetOrganization ...
// @Title 组织架构列表
// @Description 组织架构列表
// @Param	kindergarten_id           query	int	     true		"幼儿园ID"
// @Param	page                      query	int	     false		"页数"
// @Param	per_page                  query	int	     false		"每页显示条数"
// @Success 200 {object} models.Organizational
// @Failure 403
// @router / [get]
func (o *OrganizationalController) GetOrganization() {
	prepage, _ := o.GetInt("per_page", 20)
	page, _ := o.GetInt("page")
	kindergarten_id, _ := o.GetInt("kindergarten_id")
	valid := validation.Validation{}
	valid.Required(kindergarten_id, "kindergarten_id").Message("幼儿园编号不能为空")
	if valid.HasErrors() {
		o.Data["json"] = JSONStruct{"error", 1001, nil, valid.Errors[0].Message}
		o.ServeJSON()
	} else {
		v := models.GetOrganization(kindergarten_id, page, prepage)
		if v == nil {
			o.Data["json"] = JSONStruct{"error", 1005, nil, "获取组织架构失败"}
		} else {
			o.Data["json"] = JSONStruct{"success", 0, v, "获取组织架构成功"}
		}
		o.ServeJSON()
	}
}

// AddOrganization ...
// @Title 添加组织架构
// @Description 添加组织架构
// @Param	parent_id           query	int	         true		"父级ID"
// @Param	kindergarten_id     query	int	         true		"幼儿园ID"
// @Param	name                query	string	     true		"组织架构名字"
// @Param	type                query	int	         true		"类型"
// @Success 200 {object} models.Organizational
// @Failure 403
// @router /addorgan [post]
func (o *OrganizationalController) AddOrganization() {
	name := o.GetString("name")
	ty, _ := o.GetInt("type")
	parent_id, _ := o.GetInt("parent_id")
	kindergarten_id, _ := o.GetInt("kindergarten_id")
	valid := validation.Validation{}
	valid.Required(name, "name").Message("组织架构名称不能为空")
	if valid.HasErrors() {
		o.Data["json"] = JSONStruct{"error", 1001, nil, valid.Errors[0].Message}
		o.ServeJSON()
	} else {
		_, err := models.AddOrganization(name, ty, parent_id, kindergarten_id)
		if err != nil {
			o.Data["json"] = JSONStruct{"error", 1005, nil, err.Error()}
		} else {
			o.Data["json"] = JSONStruct{"success", 0, nil, "添加组织架构成功"}
		}
		o.ServeJSON()
	}

}

// DelOrganization ...
// @Title 删除组织架构
// @Description 删除组织架构
// @Param	organization_id           query	int	         true		"组织架构ID"
// @Success 200 {object} models.Organizational
// @Failure 403
// @router /delorgan [delete]
func (o *OrganizationalController) DelOrganization() {
	organization_id, _ := o.GetInt("organization_id")
	valid := validation.Validation{}
	valid.Required(organization_id, "organization_id").Message("组织架构id不能为空")
	if valid.HasErrors() {
		o.Data["json"] = JSONStruct{"error", 1001, nil, valid.Errors[0].Message}
		o.ServeJSON()
	} else {
		_, err := models.DelOrganization(organization_id)
		if err != nil {
			o.Data["json"] = JSONStruct{"error", 1005, nil, err.Error()}
		} else {
			o.Data["json"] = JSONStruct{"success", 0, nil, "删除组织架构成功"}
		}
		o.ServeJSON()
	}

}

// UpOrganization ...
// @Title 编辑组织架构
// @Description 编辑组织架构
// @Param	organization_id           query	int	         true		"组织架构ID"
// @Success 200 {object} models.Organizational
// @Failure 403
// @router / [put]
func (o *OrganizationalController) UpOrganization() {
	organization_id, _ := o.GetInt("organizational_id")
	name := o.GetString("name")
	valid := validation.Validation{}
	valid.Required(organization_id, "organization_id").Message("组织架构id不能为空")
	valid.Required(name, "name").Message("组织架构名称不能为空")
	if valid.HasErrors() {
		o.Data["json"] = JSONStruct{"error", 1001, nil, valid.Errors[0].Message}
		o.ServeJSON()
	} else {
		_, err := models.UpOrganization(organization_id, name)
		if err != nil {
			o.Data["json"] = JSONStruct{"error", 1005, nil, err.Error()}
		} else {
			o.Data["json"] = JSONStruct{"success", 0, nil, "编辑组织架构成功"}
		}
		o.ServeJSON()
	}
}

// Principal ...
// @Title 组织架构成员和负责人
// @Description 组织架构成员和负责人
// @Param	kindergarten_id           query	int	     true		"幼儿园ID"
// @Param	class_id                  query	int	     false		"班级id"
// @Param	page                      query	int	     false		"页数"
// @Param	per_page                  query	int	     false		"每页显示条数"
// @Success 200 {object} models.Organizational
// @Failure 403
// @router /principal [get]
func (o *OrganizationalController) Principal() {
	prepage, _ := o.GetInt("per_page", 20)
	page, _ := o.GetInt("page")
	class_id, _ := o.GetInt("class_id")
	valid := validation.Validation{}
	valid.Required(class_id, "class_id").Message("班级id不能为空")
	if valid.HasErrors() {
		o.Data["json"] = JSONStruct{"error", 1001, nil, valid.Errors[0].Message}
		o.ServeJSON()
	} else {
		v := models.Principal(class_id, page, prepage)
		if v == nil {
			o.Data["json"] = JSONStruct{"error", 1005, nil, "获取失败"}
		} else {
			o.Data["json"] = JSONStruct{"success", 0, v, "获取成功"}
		}
		o.ServeJSON()
	}
}
