package controllers

import (
	"kindergarten-service-go/models"

	"github.com/astaxie/beego/validation"

	"github.com/astaxie/beego"
)

// OrganizationalMemberController operations for OrganizationalMember
type OrganizationalMemberController struct {
	beego.Controller
}

// URLMapping ...
func (c *OrganizationalMemberController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
}

// Post ...
// @Title 班级添加成员
// @Description 班级添加成员
// @Param	organizational_id		body 	int	    true		"班级ID"
// @Param	member_ids		        body 	string	true		"教师ID(,分割)"
// @Param	is_principal		        body 	int   	true		"是不是负责人（0不是，1是"）
// @Param	type		                body 	int  	true		"身份（0教师，1学生"）
// @Success 201 {int} models.OrganizationalMember
// @Failure 403 body is empty
// @router / [post]
func (c *OrganizationalMemberController) Post() {
	ty, _ := c.GetInt("type")
	member_ids := c.GetString("member_ids")
	is_principal, _ := c.GetInt("is_principal")
	organizational_id, _ := c.GetInt("organizational_id")
	valid := validation.Validation{}
	valid.Required(member_ids, "member_ids").Message("成员id不能为空")
	valid.Required(organizational_id, "organizational_id").Message("班级id不能为空")
	if valid.HasErrors() {
		c.Data["json"] = JSONStruct{"error", 1001, nil, valid.Errors[0].Message}
		c.ServeJSON()
	} else {
		_, err := models.AddMembers(ty, member_ids, organizational_id, is_principal)
		if err != nil {
			c.Data["json"] = JSONStruct{"error", 1006, nil, err.Error()}
		} else {
			c.Data["json"] = JSONStruct{"success", 0, nil, "保存成功"}
		}
		c.ServeJSON()
	}
}
