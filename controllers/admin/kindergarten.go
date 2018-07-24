package admin

import (
	"kindergarten-service-go/models"
	"strconv"

	"github.com/astaxie/beego/validation"
)

//幼儿园
type KindergartenController struct {
	BaseController
}

// GetIntroduceInfo ...
// @Title 幼儿园介绍详情
// @Description 幼儿园介绍详情
// @Param	id		path 	string	true		"幼儿园ID"
// @Success 200 {object} models.Kindergarten
// @Failure 403 :id is empty
// @router /:id [get]
func (c *KindergartenController) GetIntroduceInfo() {
	prepage, _ := c.GetInt("per_page", 20)
	page, _ := c.GetInt("page")
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	v := models.GetKindergartenById(id, page, prepage)
	if v == nil {
		c.Data["json"] = JSONStruct{"error", 1005, v, "获取失败"}
	} else {
		c.Data["json"] = JSONStruct{"success", 0, v, "获取成功"}
	}
	c.ServeJSON()
}

// SetPrincipal ...
// @Title 添加园长/教师到幼儿园 未激活状态
// @Description 设置园长
// @Param	user_id		        path 	int	true		"用户ID"
// @Param	kindergarten_id		path 	int	true		"幼儿园ID"
// @Param	role		        path 	int	true		"身份（1 园长 5 ）"
// @Success 200 {object} models.Kindergarten
// @Failure 403 :id is empty
// @router / [post]
func (c *KindergartenController) SetPrincipal() {
	user_id, _ := c.GetInt("user_id")
	kindergarten_id, _ := c.GetInt("kindergarten_id")
	role, _ := c.GetInt("role")
	valid := validation.Validation{}
	valid.Required(user_id, "user_id").Message("用户ID不能为空")
	valid.Required(kindergarten_id, "kindergarten_id").Message("幼儿园ID不能为空")
	if valid.HasErrors() {
		c.Data["json"] = JSONStruct{"error", 1001, nil, valid.Errors[0].Message}
		c.ServeJSON()
	} else {
		err := models.AddPrincipal(user_id, kindergarten_id, role)
		if err != nil {
			c.Data["json"] = JSONStruct{"error", 1003, err, "保存失败"}
		} else {
			c.Data["json"] = JSONStruct{"success", 0, nil, "保存成功"}
		}
		c.ServeJSON()
	}
}

// GetAll ...
// @Title 幼儿园列表
// @Description 幼儿园列表
// @Param	page                  query	int	     false		"页数"
// @Param	per_page              query	int	     false		"每页显示条数"
// @Success 200 {object} models.Kindergarten
// @Failure 403
// @router / [get]
func (c *KindergartenController) GetAll() {
	search := c.GetString("search")
	prepage, _ := c.GetInt("per_page", 20)
	page, _ := c.GetInt("page")
	v := models.GetAll(page, prepage, search)
	if v == nil {
		c.Data["json"] = JSONStruct{"error", 1005, nil, "获取失败"}
	} else {
		c.Data["json"] = JSONStruct{"success", 0, v, "获取成功"}
	}
	c.ServeJSON()
}

// StudentClass ...
// @Title 学生姓名搜索班级
// @Description 学生姓名搜索班级
// @Param	name                  query	 string	     true		"姓名"
// @Param	page                  query	   int	     false		"页数"
// @Param	per_page              query	   int	     false		"每页显示条数"
// @Success 200 {object} models.Kindergarten
// @Failure 403
// @router /student_class [get]
func (c *KindergartenController) StudentClass() {
	name := c.GetString("name")
	prepage, _ := c.GetInt("per_page", 20)
	page, _ := c.GetInt("page")
	valid := validation.Validation{}
	valid.Required(name, "name").Message("学生姓名不能为空")
	if valid.HasErrors() {
		c.Data["json"] = JSONStruct{"error", 1001, nil, valid.Errors[0].Message}
		c.ServeJSON()
	} else {
		v := models.StudentClass(page, prepage, name)
		if v == nil {
			c.Data["json"] = JSONStruct{"error", 1005, nil, "获取失败"}
		} else {
			c.Data["json"] = JSONStruct{"success", 0, v, "获取成功"}
		}
		c.ServeJSON()
	}
}

// SetKindergarten ...
// @Title 添加幼儿园
// @Description 添加幼儿园
// @Param	name		                path 	int	true		"幼儿园名称"
// @Param	license_no   		        path 	int	true		"执照号"
// @Param	kinder_grade		        path 	int	true		"幼儿园级别"
// @Param	kinder_child_no		        path 	int	true		"分校数"
// @Param	address      		        path 	int	true		"地址"
// @Param	tenant_id    		        path 	int	true		"租户，企业编号"
// @Success 200 {object} models.Kindergarten
// @Failure 403 :id is empty
// @router /set_kindergarten [post]
func (c *KindergartenController) SetKindergarten() {
	name := c.GetString("name")
	license_no, _ := c.GetInt("license_no")
	kinder_grade := c.GetString("kinder_grade")
	kinder_child_no, _ := c.GetInt("kinder_child_no")
	address := c.GetString("address")
	tenant_id, _ := c.GetInt("tenant_id")
	valid := validation.Validation{}
	valid.Required(name, "name").Message("幼儿园名称不能为空")
	valid.Required(license_no, "license_no").Message("执照号不能为空")
	valid.Required(kinder_grade, "kinder_grade").Message("幼儿园级别不能为空")
	valid.Required(kinder_child_no, "kinder_child_no").Message("分校数")
	valid.Required(address, "address").Message("地址不能为空")
	valid.Required(tenant_id, "tenant_id").Message("企业编号不能为空")
	if valid.HasErrors() {
		c.Data["json"] = JSONStruct{"error", 1001, nil, valid.Errors[0].Message}
		c.ServeJSON()
	} else {
		err := models.AddKindergarten(name, license_no, kinder_grade, kinder_child_no, address, tenant_id)
		if err != nil {
			c.Data["json"] = JSONStruct{"error", 1003, err, "保存失败"}
		} else {
			c.Data["json"] = JSONStruct{"success", 0, nil, "保存成功"}
		}
		c.ServeJSON()
	}
}

// delete ...
// @Title 删除幼儿园
// @Description 删除幼儿园
// @Param	kindergarten_id		path 	int	true		"幼儿园ID"
// @Success 200 {object} models.Kindergarten
// @Failure 403 :id is empty
// @router / [delete]
func (c *KindergartenController) Delete() {
	kindergarten_id, _ := c.GetInt("kindergarten_id")
	valid := validation.Validation{}
	valid.Required(kindergarten_id, "kindergarten_id").Message("幼儿园ID不能为空")
	if valid.HasErrors() {
		c.Data["json"] = JSONStruct{"error", 1001, nil, valid.Errors[0].Message}
		c.ServeJSON()
	} else {
		err := models.DeleteKinder(kindergarten_id)
		if err != nil {
			c.Data["json"] = JSONStruct{"error", 1003, err.Error(), "删除失败"}
		} else {
			c.Data["json"] = JSONStruct{"success", 0, nil, "删除成功"}
		}
		c.ServeJSON()
	}
}

// updata ...
// @Title 编辑幼儿园
// @Description 编辑幼儿园
// @Param	kindergarten_id		path 	int	true		"幼儿园ID"
// @Success 200 {object} models.Kindergarten
// @Failure 403 :id is empty
// @router /:id [put]
func (c *KindergartenController) Update() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	name := c.GetString("name")
	license_no, _ := c.GetInt("license_no")
	kinder_grade := c.GetString("kinder_grade")
	kinder_child_no, _ := c.GetInt("kinder_child_no")
	address := c.GetString("address")
	tenant_id, _ := c.GetInt("tenant_id")
	valid := validation.Validation{}
	valid.Required(name, "name").Message("幼儿园名称不能为空")
	valid.Required(license_no, "license_no").Message("执照号不能为空")
	valid.Required(kinder_grade, "kinder_grade").Message("幼儿园级别不能为空")
	valid.Required(kinder_child_no, "kinder_child_no").Message("分校数")
	valid.Required(address, "address").Message("地址不能为空")
	valid.Required(tenant_id, "tenant_id").Message("企业编号不能为空")
	if valid.HasErrors() {
		c.Data["json"] = JSONStruct{"error", 1001, nil, valid.Errors[0].Message}
		c.ServeJSON()
	} else {
		err := models.UpdataKinder(id, name, license_no, kinder_grade, kinder_child_no, address, tenant_id)
		if err != nil {
			c.Data["json"] = JSONStruct{"error", 1003, err.Error(), "编辑失败"}
		} else {
			c.Data["json"] = JSONStruct{"success", 0, nil, "编辑成功"}
		}
		c.ServeJSON()
	}
}

// GetKg ...
// @Title 登陆幼儿园信息
// @Description 登陆幼儿园信息
// @Param	kindergarten_id		path 	int	true		"幼儿园ID"
// @Param	user_id		        path 	int	true		"用户ID"
// @Success 200 {object} models.Kindergarten
// @Failure 403 :id is empty
// @router /getkg [get]
func (c *KindergartenController) GetKg() {
	user_id, _ := c.GetInt("user_id")
	kindergarten_id, _ := c.GetInt("kindergarten_id")
	valid := validation.Validation{}
	valid.Required(user_id, "user_id").Message("用户id为空")
	valid.Required(kindergarten_id, "kindergarten_id").Message("幼儿园id不能为空")
	if valid.HasErrors() {
		c.Data["json"] = JSONStruct{"error", 1001, nil, valid.Errors[0].Message}
		c.ServeJSON()
	} else {
		v, err := models.GetKg(user_id, kindergarten_id)
		if err != nil {
			c.Data["json"] = JSONStruct{"error", 1003, err.Error(), "获取失败"}
		} else {
			c.Data["json"] = JSONStruct{"success", 0, v, "获取成功"}
		}
		c.ServeJSON()
	}
}

// GetKinderMbmber ...
// @Title oms-幼儿园所有成员
// @Description oms-幼儿园所有成员
// @Param	kindergarten_id		path 	int	true		"幼儿园ID"
// @Success 200 {object} models.Kindergarten
// @Failure 403 :id is empty
// @router /get_member [get]
func (c *KindergartenController) GetKinderMbmber() {
	kindergarten_id, _ := c.GetInt("kindergarten_id")
	prepage, _ := c.GetInt("per_page", 20)
	page, _ := c.GetInt("page")
	valid := validation.Validation{}
	valid.Required(kindergarten_id, "kindergarten_id").Message("幼儿园id不能为空")
	if valid.HasErrors() {
		c.Data["json"] = JSONStruct{"error", 1001, nil, valid.Errors[0].Message}
		c.ServeJSON()
	} else {
		v, err := models.GetKinderMbmber(kindergarten_id, page, prepage)
		if err != nil {
			c.Data["json"] = JSONStruct{"error", 1003, err.Error(), "获取失败"}
		} else {
			c.Data["json"] = JSONStruct{"success", 0, v, "获取成功"}
		}
		c.ServeJSON()
	}
}

// FoodClass ...
// @Title 饮食班级
// @Description 饮食班级
// @Param	kindergarten_id		path 	int	true		"幼儿园ID"
// @Success 200 {object} models.Kindergarten
// @Failure 403 :id is empty
// @router /food_class [get]
func (c *KindergartenController) FoodClass() {
	kindergarten_id, _ := c.GetInt("kindergarten_id")
	valid := validation.Validation{}
	valid.Required(kindergarten_id, "kindergarten_id").Message("幼儿园id不能为空")
	if valid.HasErrors() {
		c.Data["json"] = JSONStruct{"error", 1001, nil, valid.Errors[0].Message}
		c.ServeJSON()
	} else {
		v, err := models.FoodClass(kindergarten_id)
		if err != nil {
			c.Data["json"] = JSONStruct{"error", 1003, err.Error(), "获取失败"}
		} else {
			c.Data["json"] = JSONStruct{"success", 0, v, "获取成功"}
		}
		c.ServeJSON()
	}
}

// FoodScale ...
// @Title 饮食比例
// @Description 饮食比例
// @Param	class_type		    path 	int	true		"班级类型"
// @Param	is_muslim		    path 	int	true		"是否清真"
// @Success 200 {object} models.Kindergarten
// @Failure 403 :id is empty
// @router /food_scale [get]
func (c *KindergartenController) FoodScale() {
	is_muslim, _ := c.GetInt("is_muslim")
	class_type := c.GetString("class_type")
	kindergarten_id, _ := c.GetInt("kindergarten_id")
	valid := validation.Validation{}
	valid.Required(kindergarten_id, "kindergarten_id").Message("幼儿园id不能为空")
	if valid.HasErrors() {
		c.Data["json"] = JSONStruct{"error", 1001, nil, valid.Errors[0].Message}
		c.ServeJSON()
	} else {
		v, err := models.FoodScale(is_muslim, kindergarten_id, class_type)
		if err != nil {
			c.Data["json"] = JSONStruct{"error", 1003, err.Error(), "获取失败"}
		} else {
			c.Data["json"] = JSONStruct{"success", 0, v, "获取成功"}
		}
		c.ServeJSON()
	}
}
