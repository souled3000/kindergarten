package healthy

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/validation"
	"kindergarten-service-go/models/healthy"
	"fmt"
	"strconv"
)

//餐检
type InspectController struct {
	beego.Controller
}

// URLMapping ...
func (c *InspectController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
}

// @Title 添加检查
// @Description 添加喂药申请
// @Param   class_id     			formData    int  	true        "班级ID"
// @Param   student_id     			formData    int  	true        "学生ID"
// @Param   types     				formData    int	    true        "检查类型"
// @Param   abnormal     			formData    string  true        "异常情况"
// @Param   handel     				formData    string  true        "处理方式"
// @Param   url     				formData    string  true        "照片留档"
// @Param   infect     				formData    string  true        "是否传染（1，否2，是）"
// @Success 0 {int} models.Drug.Id
// @Failure 1001 补全信息
// @Failure 1003 保存失败
// @router / [post]
func (c *InspectController) Post() {
	class_id, _:= c.GetInt("class_id")
	student_id, _:= c.GetInt("student_id")
	types, _:= c.GetInt("types")
	abnormal := c.GetString("abnormal")
	handel := c.GetString("handel")
	url := c.GetString("url")
	infect, _:= c.GetInt("infect")
	drug_id, _:= c.GetInt("drug_id")
	user_id, _:= c.GetInt("user_id")
	kindergarten_id, _:= c.GetInt("kindergarten_id")

	valid := validation.Validation{}
	valid.Required(student_id, "student_id").Message("学生ID不能为空")
	valid.Required(class_id,"class_id").Message("班级ID不能为空")
	valid.Required(types,"types").Message("检查类型不能为空")
	valid.Required(abnormal,"abnormal").Message("异常情况不能为空")
	valid.Required(user_id,"user_id").Message("用户ID不能为空")
	valid.Required(handel,"handel").Message("异常情况不能为空")
	valid.Required(infect,"infect").Message("传染情况不能为空")
	valid.Required(kindergarten_id,"kindergarten_id").Message("幼儿园ID不能为空")

	if valid.HasErrors() {
		c.Data["json"] = JSONStruct{"error", 1001, "", valid.Errors[0].Message}

		c.ServeJSON()
	}

	w := healthy.Inspect{
		StudentId:student_id,
		ClassId:class_id,
		Types:types,
		Abnormal:abnormal,
		Handel:handel,
		Url:url,
		Infect:infect,
		DrugId:drug_id,
		NoteTaker:user_id,
		KindergartenId:kindergarten_id,
	}
	if err := w.Save(); err == nil {
		c.Data["json"] = JSONStruct{"success", 0, "", "申请成功"}
	} else {
		fmt.Println(err)
		c.Data["json"] = JSONStruct{"error", 1003, "", "申请失败"}
	}

	c.ServeJSON()
}

// GetAll ...
// @Title GetAll
// @Description 餐检列表
// @Param	page			query	int	false		"第几页"
// @Param	kindergarten_id	query	int	true		"幼儿园ID"
// @Param	per_page		query	int	true		"页数"
// @Param	class_id		query	int	true		"班级ID"
// @Success 0 {object} 		shanxi.SxWorks
// @Failure 1001 		参数不能为空
// @Failure 1005 		获取失败
// @router / [get]
func (c *InspectController) GetAll() {
	var f *healthy.Inspect
	page, _ := c.GetInt("page")
	kindergarten_id, _:= c.GetInt("kindergarten_id")
	class_id, _:= c.GetInt("class_id")
	types, _:= c.GetInt("types")
	perPage, _ := c.GetInt("per_page")

	//验证参数是否为空
	valid := validation.Validation{}
	valid.Required(kindergarten_id,"kindergarten_id").Message("幼儿园ID不能为空")
	if valid.HasErrors(){
		c.Data["json"] = JSONStruct{"error", 1001, struct {}{}, valid.Errors[0].Message}
		c.ServeJSON()

	}
	if works, err := f.GetAll(page, perPage, kindergarten_id, class_id, types ); err == nil {
		c.Data["json"] = JSONStruct{"success", 0, works, "获取成功"}
	} else {
		c.Data["json"] = JSONStruct{"error", 1005, err, "获取失败"}
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
func (c *InspectController) Delete() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	v := healthy.DeleteInspect(id)
	if v == nil {
		c.Data["json"] = JSONStruct{"error", 1003, nil, "删除失败"}
	} else {
		c.Data["json"] = JSONStruct{"success", 0, nil, "删除成功"}
	}
	c.ServeJSON()
}