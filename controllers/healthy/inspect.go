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
	class_name := c.GetString("class_name")
	class_id, _:= c.GetInt("class_id")
	student_id, _:= c.GetInt("student_id")
	types, _:= c.GetInt("types")
	abnormal := c.GetString("abnormal")
	handel := c.GetString("handel")
	url := c.GetString("url")
	infect, _:= c.GetInt("infect")
	drug_id, _:= c.GetInt("drug_id")
	teacher_id, _:= c.GetInt("teacher_id")
	kindergarten_id, _:= c.GetInt("kindergarten_id")
	date := c.GetString("date")

	valid := validation.Validation{}
	valid.Required(class_name,"class_name").Message("班级名称不能为空")
	valid.Required(student_id, "student_id").Message("学生ID不能为空")
	valid.Required(class_id,"class_id").Message("班级ID不能为空")
	valid.Required(teacher_id,"teacher_id").Message("教师ID不能为空")
	valid.Required(kindergarten_id,"kindergarten_id").Message("幼儿园ID不能为空")

	if valid.HasErrors() {
		c.Data["json"] = JSONStruct{"error", 1001, "", valid.Errors[0].Message}

		c.ServeJSON()
		c.StopRun()
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
		TeacherId:teacher_id,
		KindergartenId:kindergarten_id,
		ClassName:class_name,
		Date:date,
	}
	if err := w.Save(); err == nil {
		c.Data["json"] = JSONStruct{"success", 0, "", "记录成功"}
	} else {
		fmt.Println(err)
		c.Data["json"] = JSONStruct{"error", 1003, "", "记录失败"}
	}

	c.ServeJSON()
}

// GetAll ...
// @Title GetAll
// @Description 餐检列表
// @Param	page			query	int		false		"第几页"
// @Param	kindergarten_id	query	int		true		"幼儿园ID"
// @Param	per_page		query	int		true		"页数"
// @Param	class_id		query	int		false		"班级ID"
// @Param	role			query	int		true		"身份类型"
// @Param	date			query	string	true		"餐检时间"
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
	role, _:= c.GetInt("role")
	date := c.GetString("time")
	bady_id, _:= c.GetInt("baby_id")

	//验证参数是否为空
	valid := validation.Validation{}
	valid.Required(kindergarten_id,"kindergarten_id").Message("幼儿园ID不能为空")
	valid.Required(kindergarten_id,"role").Message("用户身份不能为空")
	if valid.HasErrors(){
		c.Data["json"] = JSONStruct{"error", 1001, struct {}{}, valid.Errors[0].Message}
		c.ServeJSON()
		c.StopRun()
	}
	if works, err := f.GetAll(page, perPage, kindergarten_id, class_id, types, role, bady_id, date ); err == nil {
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

// GetAll ...
// @Title GetAll
// @Description 统计
// @Param	kindergarten_id	query	int	true		"幼儿园ID"
// @Success 0 {object} 	healthy.Counts
// @Failure 1001 		参数不能为空
// @Failure 1003 		获取失败
// @router /counts/ [get]
func (c *InspectController) Counts() {
	kindergarten_id, _:= c.GetInt("kindergarten_id")

	valid := validation.Validation{}
	valid.Required(kindergarten_id,"kindergarten_id").Message("幼儿园ID不能为空")
	if valid.HasErrors() {
		c.Data["json"] = JSONStruct{"error",1001, struct {}{},valid.Errors[0].Message}

		c.ServeJSON()
		c.StopRun()
	}
	v := healthy.Counts(kindergarten_id)

	if v == nil {
		c.Data["json"] = JSONStruct{"error", 1003, nil, "获取失败"}
	} else {
		c.Data["json"] = JSONStruct{"success", 0, v, "获取成功"}
	}

	c.ServeJSON()
}

// GetAll ...
// @Title GetAll
// @Description 详情
// @Param	kindergarten_id	query	int	true		"幼儿园ID"
// @Success 0 {object} 	healthy.Counts
// @Failure 1001 		参数不能为空
// @Failure 1003 		获取失败
// @router /:id
func (c *InspectController) Inspect() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	v := healthy.InspectInfo(id)
	if v == nil {
		c.Data["json"] = JSONStruct{"error", 1005, nil, "获取失败"}
	} else {
		c.Data["json"] = JSONStruct{"success", 0, v, "获取成功"}
	}
	c.ServeJSON()
}

// @Title 编辑检查
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
// @router /:id [put]
func (c *InspectController) Put() {
	var id int
	c.Ctx.Input.Bind(&id, ":id")

	class_name := c.GetString("class_name")
	class_id, _:= c.GetInt("class_id")
	student_id, _:= c.GetInt("student_id")
	types, _:= c.GetInt("types")
	abnormal := c.GetString("abnormal")
	handel := c.GetString("handel")
	url := c.GetString("url")
	infect, _:= c.GetInt("infect")
	drug_id, _:= c.GetInt("drug_id")
	teacher_id, _:= c.GetInt("teacher_id")
	kindergarten_id, _:= c.GetInt("kindergarten_id")

	valid := validation.Validation{}
	valid.Required(class_name,"class_name").Message("班级名称不能为空")
	valid.Required(student_id, "student_id").Message("学生ID不能为空")
	valid.Required(class_id,"class_id").Message("班级ID不能为空")
	valid.Required(types,"types").Message("检查类型不能为空")
	valid.Required(abnormal,"abnormal").Message("异常情况不能为空")
	valid.Required(teacher_id,"teacher_id").Message("教师ID不能为空")
	valid.Required(handel,"handel").Message("异常情况不能为空")
	valid.Required(infect,"infect").Message("传染情况不能为空")
	valid.Required(kindergarten_id,"kindergarten_id").Message("幼儿园ID不能为空")

	if valid.HasErrors() {
		c.Data["json"] = JSONStruct{"error", 1001, "", valid.Errors[0].Message}

		c.ServeJSON()
		c.StopRun()
	}

	w := healthy.Inspect{
		Id:id,
		StudentId:student_id,
		ClassId:class_id,
		Types:types,
		Abnormal:abnormal,
		Handel:handel,
		Url:url,
		Infect:infect,
		DrugId:drug_id,
		TeacherId:teacher_id,
		KindergartenId:kindergarten_id,
		ClassName:class_name,
	}
	if err := w.SaveInspect(); err == nil {
		c.Data["json"] = JSONStruct{"success", 0, "", "申请成功"}
	} else {
		fmt.Println(err)
		c.Data["json"] = JSONStruct{"error", 1003, "", "申请失败"}
	}

	c.ServeJSON()
}

// GetAll ...
// @Title GetAll
// @Description 异常档案列表
// @Param	page			query	int		false		"第几页"
// @Param	kindergarten_id	query	int		true		"幼儿园ID"
// @Param	per_page		query	int		true		"页数"
// @Param	class_id		query	int		false		"班级ID"
// @Param	role			query	int		true		"身份类型"
// @Param	date			query	string	true		"餐检时间"
// @Success 0 {object} 		shanxi.SxWorks
// @Failure 1001 		参数不能为空
// @Failure 1005 		获取失败
// @router /archives/ [get]
func (c *InspectController) Abnormal() {
	var f *healthy.Inspect
	page, _ := c.GetInt("page")
	kindergarten_id, _:= c.GetInt("kindergarten_id")
	class_id, _:= c.GetInt("class_id")
	perPage, _ := c.GetInt("per_page")
	date := c.GetString("time")
	search := c.GetString("search")

	//验证参数是否为空
	valid := validation.Validation{}
	valid.Required(kindergarten_id,"kindergarten_id").Message("幼儿园ID不能为空")
	if valid.HasErrors(){
		c.Data["json"] = JSONStruct{"error", 1001, struct {}{}, valid.Errors[0].Message}
		c.ServeJSON()
		c.StopRun()
	}
	if works, err := f.Abnormals(page, perPage, kindergarten_id, class_id, date, search ); err == nil {
		c.Data["json"] = JSONStruct{"success", 0, works, "获取成功"}
	} else {
		c.Data["json"] = JSONStruct{"error", 1005, err, "获取失败"}
	}

	c.ServeJSON()
}

// GetAll ...
// @Title GetAll
// @Description 项目详情
// @Param	page			query	int		false		"第几页"
// @Param	per_page		query	int		true		"页数"
// @Param	class_id		query	int		false		"班级ID"
// @Param	role			query	int		true		"身份类型"
// @Param	date			query	string	true		"餐检时间"
// @Success 0 {object} 		shanxi.SxWorks
// @Failure 1001 		参数不能为空
// @Failure 1005 		获取失败
// @router /project/ [get]
func (c *InspectController) Project() {
	var f *healthy.Inspect
	page, _ := c.GetInt("page")
	kindergarten_id, _:= c.GetInt("kindergarten_id")
	class_id, _:= c.GetInt("class_id")
	perPage, _ := c.GetInt("per_page")
	body_id, _:= c.GetInt("body_id")
	baby_id, _:= c.GetInt("baby_id")
	column := c.GetString("column")

	if works, err := f.Project(page, perPage, kindergarten_id, class_id, body_id,baby_id, column ); err == nil {
		c.Data["json"] = JSONStruct{"success", 0, works, "获取成功"}
	} else {
		c.Data["json"] = JSONStruct{"error", 1005, err, "获取失败"}
	}

	c.ServeJSON()
}

// GetAll ...
// @Title GetAll
// @Description 项目详情
// @Param	page			query	int		false		"第几页"
// @Param	per_page		query	int		true		"页数"
// @Param	class_id		query	int		false		"班级ID"
// @Param	role			query	int		true		"身份类型"
// @Param	date			query	string	true		"餐检时间"
// @Success 0 {object} 		shanxi.SxWorks
// @Failure 1001 		参数不能为空
// @Failure 1005 		获取失败
// @router /projectNew/ [get]
func (c *InspectController) ProjectNew() {
	var f *healthy.Inspect
	page, _ := c.GetInt("page")
	kindergarten_id, _:= c.GetInt("kindergarten_id")
	class_id, _:= c.GetInt("class_id")
	perPage, _ := c.GetInt("per_page")
	body_id, _:= c.GetInt("body_id")
	baby_id, _:= c.GetInt("baby_id")
	column := c.GetString("column")


	if works, err := f.ProjectNew(page, perPage, kindergarten_id, class_id, body_id,baby_id, column); err == nil {
		c.Data["json"] = JSONStruct{"success", 0, works, "获取成功"}
	} else {
		c.Data["json"] = JSONStruct{"error", 1005, err, "获取失败"}
	}

	c.ServeJSON()
}

// GetAll ...
// @Title GetAll
// @Description 体重健康统计
// @Param	kindergarten_id		query	int		true		"幼儿园ID"
// @Param	types				query	string	true		"类型"
// @Success 0 {object} 			shanxi.SxWorks
// @Failure 1001 		参数不能为空
// @Failure 1005 		获取失败
// @router /chart/ [get]
func (c *InspectController) Chart() {
	var f *healthy.Inspect

	kindergarten_id, _:= c.GetInt("kindergarten_id")
	types, _:= c.GetInt("types")

	//验证参数是否为空
	valid := validation.Validation{}
	valid.Required(kindergarten_id,"kindergarten_id").Message("幼儿园ID不能为空")
	if valid.HasErrors(){
		c.Data["json"] = JSONStruct{"error", 1001, struct {}{}, valid.Errors[0].Message}
		c.ServeJSON()
		c.StopRun()
	}
	if works, err :=f.Chart(kindergarten_id, types); err == nil {
		c.Data["json"] = JSONStruct{"success", 0, works, "获取成功"}
	} else {
		c.Data["json"] = JSONStruct{"error", 1005, err, "获取失败"}
	}

	c.ServeJSON()

}
