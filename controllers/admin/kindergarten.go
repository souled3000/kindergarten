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
