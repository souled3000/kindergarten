package admin

import (
	"kindergarten-service-go/models"
	"strconv"
	"github.com/astaxie/beego/validation"
)

type ExceptionalChildController struct {
	BaseController
}

func (c *ExceptionalChildController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
	c.Mapping("Get", c.Get)
}



// GetAll ...
// @Title 特殊儿童列表/搜索特殊儿童
// @Description 特殊儿童列表/搜索特殊儿童
// @Param	child_name	query	string	false	"特殊儿童姓名"
// @Param	somatotype	query	int		false	"体质类型"
// @Param	page		query	string	false	"当前页，默认为1"
// @Param	per_page	query	string	false	"每页显示条数，默认为10"
// @Param	keyword		query	string	false	"关键字(特殊儿童姓名/特殊儿童过敏源)"
// @Success 0 			{object} 	models.ExceptionalChild
// @Failure 1005 获取失败
// @router / [get]
func (c *ExceptionalChildController) GetAll() {
	child_name := c.GetString("child_name")
	somatotype, _ := c.GetInt8("somatotype")
	// 关键字
	keyword := c.GetString("keyword")
	// page_num
	page, _ := c.GetInt64("page")

	// limit
	limit, _ := c.GetInt64("per_page")
	if info, err := models.GetAllExceptionalChild(child_name, somatotype, page, limit, keyword); err == nil {
		c.Data["json"] = JSONStruct{"success", 0, info, "获取成功"}

	} else {
		c.Data["json"] = JSONStruct{"error", 1005, nil, "获取失败"}
	}
	c.ServeJSON()
}




// Post ...
// @Title 						新增特殊儿童
// @Description 				新增特殊儿童
// @Param	child_name			formData  string	true		"特殊儿童姓名"
// @Param	class				formData  int 		true		"特殊儿童班级"
// @Param	somatotype			formData  int		true		"体质类型"
// @Param	allergen			formData  string	true		"过敏源"
// @Param	source				formData  int		true		"信息来源"
// @Param	kindergarten_id		formData  int		true		"幼儿园ID"
// @Param	creator				formData  int		true		"创建人"
// @Param	student_id			formData  int		true		"学生ID"
// @Success 0					{json} JSONSturct
// @Failure 1003 				新增失败
// @router / [post]
func (c *ExceptionalChildController) Post() {
	// 特殊儿童姓名
	child_name := c.GetString("child_name")
	// 特殊儿童班级
	class, _ := c.GetInt("class")
	// 体质类型
	somatotype, _ := c.GetInt8("somatotype")
	// 过敏源
	allergen := c.GetString("allergen")
	// 信息来源
	source, _ := c.GetInt8("source")
	// 幼儿园ID
	kindergarten_id, _ := c.GetInt("kindergarten_id")
	// 创建人
	creator, _ := c.GetInt("creator")
	// 学生ID
	student_id, _ := c.GetInt("student_id")

	valid := validation.Validation{}

	valid.Required(child_name, "child_name").Message("儿童姓名不能为空")
	valid.Required(class, "class").Message("所在班级不能为空")
	valid.Required(somatotype, "somatotype").Message("体质类型不能为空")
	valid.Required(allergen, "allergen").Message("信息来源不能为空")
	valid.Required(kindergarten_id, "kindergarten_id").Message("幼儿园ID不能为空")
	valid.Required(creator, "creator").Message("创建人不能为空")
	valid.Required(student_id, "student_id").Message("学生ID不能为空")

	if valid.HasErrors() {
		c.Data["json"] = JSONStruct{"error", 1001, nil, valid.Errors[0].Message}
	} else {
		if _, err := models.AddExceptionalChild(child_name, class, somatotype, allergen, source, kindergarten_id, creator, student_id); err == nil {
			c.Data["json"] = JSONStruct{"success", 0, err, "新增成功"}
		} else {
			c.Data["json"] = JSONStruct{"error", 1003, nil, "新增失败"}
		}
	}
	c.ServeJSON()
}

// GetOne ...
// @Title 			按ID查询特殊儿童
// @Description 	按ID查询特殊儿童
// @Param	id		path 	string	true		"特殊儿童ID"
// @Success 0		{object}  models.ExceptionalChild
// @Failure 1005 	获取失败
// @router /:id [get]
func (c *ExceptionalChildController) GetOne() {
	// 主键ID
	idStr := c.Ctx.Input.Param(":id")
	v, err := models.GetExceptionalChildById(idStr)
	if err == nil {
		if v != nil {
			c.Data["json"] = JSONStruct{"success", 0, v, "获取成功"}
		} else {
			c.Data["json"] = JSONStruct{"error", 1002, nil, "没有相关数据"}
		}
	} else {
			c.Data["json"] = JSONStruct{"error", 1005, nil, "获取失败"}
	}
	c.ServeJSON()
}



// Put ...
// @Title 					更新特殊儿童
// @Description 			更新特殊儿童
// @Param	id				path 	string	true		"特殊儿童主键自增ID"
// @Param	child_name		body 	string	true		"特殊儿童姓名"
// @Param	class			body 	int		true		"特殊儿童班级"
// @Param	somatotype		body 	int		true		"体质类型"
// @Param	allergen		body 	string	true		"过敏源"
// @Param	source			body 	int		true		"来源信息"
// @Param	kindergarten_id	body 	int		true		"幼儿园ID"
// @Param	creator			body 	int		true		"创建人"
// @Param	student_id		body 	int		true		"学生ID"
// @Success 0 				{json} 	JSONStruct
// @Failure 1003			更新失败
// @router /:id [put]
func (c *ExceptionalChildController) Put() {
	// 主键ID
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)

	child_name := c.GetString("child_name")
	class, _ := c.GetInt("class")

	somatotype, _ := c.GetInt8("somatotype")

	allergen := c.GetString("allergen")
	source, _ := c.GetInt8("source")
	kindergarten_id, _ := c.GetInt("kindergarten_id")
	creator, _ := c.GetInt("creator")
	student_id, _ := c.GetInt("student_id")

	valid := validation.Validation{}

	valid.Required(child_name, "child_name").Message("儿童姓名不能为空")
	valid.Required(class, "class").Message("所在班级不能为空")
	valid.Required(somatotype, "somatotype").Message("体质类型不能为空")
	valid.Required(allergen, "allergen").Message("过敏源不能为空")
	valid.Required(source, "source").Message("信息来源不能为空")
	valid.Required(kindergarten_id, "kindergarten_id").Message("幼儿园ID不能为空")
	valid.Required(creator, "creator").Message("创建人不能为空")
	valid.Required(student_id, "student_id").Message("学生ID不能为空")

	if valid.HasErrors() {
		c.Data["json"] = JSONStruct{"error", 1001, nil, valid.Errors[0].Message}
	} else {
		if err := models.UpdateExceptionalChildById(id, child_name, class, somatotype, allergen, source, kindergarten_id, creator, student_id); err == nil {
			c.Data["json"] = JSONStruct{"success", 0, nil, "更新成功"}
		} else {
			c.Data["json"] = JSONStruct{"error", 1003, nil, "更新失败"}
		}
	}
	c.ServeJSON()
}

// Delete ...
// @Title 			删除特殊儿童
// @Description 	删除特殊儿童
// @Param	id		path 	string	true		"特殊儿童ID"
// @Success 0 		{json} 	JSONStruct
// @Failure 1004	删除失败
// @router /:id [delete]
func (c *ExceptionalChildController) Delete() {
	// 主键ID
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	if err := models.DeleteExceptionalChild(id); err == nil {
		c.Data["json"] = JSONStruct{"success", 0, err, "删除成功"}
	} else {
		c.Data["json"] = JSONStruct{"error", 1004, nil, "删除失败"}
	}
	c.ServeJSON()
}