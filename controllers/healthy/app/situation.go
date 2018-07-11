package app

import (
	"kindergarten-service-go/models/healthy"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/validation"
)

// SituationController operations for Situation
type SituationController struct {
	beego.Controller
}

// URLMapping ...
func (c *SituationController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
}

// GetAll ...
// @Title GetAll
// @Description 病例列表
// @Param	types			query	int	true		"类型"
// @Success 0 {object} 		healthy.Situation
// @Failure 1001 		参数不能为空
// @Failure 1005 		获取失败
// @router / [get]
func (c *SituationController) GetAll() {
	var f *healthy.Situation

	types, _ := c.GetInt("types")

	//验证参数是否为空
	valid := validation.Validation{}
	valid.Required(types, "types").Message("类型ID不能为空")
	if valid.HasErrors() {
		c.Data["json"] = JSONStruct{"error", 1001, struct{}{}, valid.Errors[0].Message}
		c.ServeJSON()
		c.StopRun()
	}
	if works, err := f.GetAll(types); err == nil {
		c.Data["json"] = JSONStruct{"success", 0, works, "获取成功"}
	} else {
		c.Data["json"] = JSONStruct{"error", 1005, err, "获取失败"}
	}

	c.ServeJSON()

}
