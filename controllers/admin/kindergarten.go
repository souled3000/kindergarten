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
// @Title 添加园长/教师/学生到幼儿园 未激活状态
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
	//	valid.Required(role, "role").Message("身份不能为空")
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
