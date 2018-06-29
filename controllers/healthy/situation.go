package healthy

import (
	"github.com/astaxie/beego"
	"kindergarten-service-go/models/healthy"
	"strconv"
	"fmt"
	"github.com/astaxie/beego/validation"
)

// 病因
type SituationController struct {
	beego.Controller
}

// URLMapping ...
func (c *SituationController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
}

// @Title 添加病因
// @Description 添加病因
// @Param   name     				formData    int  	true        "记录名称"
// @Param   types     				formData    string  true        "类型"
// @Success 0 {int} models.Drug.Id
// @Failure 1001 补全信息
// @Failure 1003 保存失败
// @router / [post]
func (c *SituationController) Post() {
	name := c.GetString("name")
	types, _:= c.GetInt("types")

	fmt.Print(name)
	fmt.Println(types)

	valid := validation.Validation{}
	valid.Required(name,"name").Message("名字1111不能为空")
	valid.Required(types,"types").Message("类型ID不能为空")
	if valid.HasErrors(){
		c.Data["json"] = JSONStruct{"error", 1001, struct {}{}, valid.Errors[0].Message}
		c.ServeJSON()
		c.StopRun()
	}

	if valid.HasErrors() {
		c.Data["json"] = JSONStruct{"error", 1001, "", valid.Errors[0].Message}

		c.ServeJSON()
		c.StopRun()
	}
	w := healthy.Situation{
		Name:name,
		Type:types,
	}
	if err := w.Post(); err == nil {

		c.Data["json"] = JSONStruct{"success", 0, "", "创建成功"}
	} else {

		c.Data["json"] = JSONStruct{"error", 1003, "", "创建失败"}
	}

	c.ServeJSON()

}

// Delete ...
// @Title Delete
// @Description 删除
// @Param	id		path 	string	true		"自增ID"
// @Success 0 {string} delete success!
// @Failure 1003 id is empty
// @router /:id [delete]
func (c *SituationController) Delete() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	v := healthy.DeleteSituation(id)
	if v == nil {
		c.Data["json"] = JSONStruct{"error", 1003, nil, "删除失败"}
	} else {
		c.Data["json"] = JSONStruct{"success", 0, nil, "删除成功"}
	}
	c.ServeJSON()
}
