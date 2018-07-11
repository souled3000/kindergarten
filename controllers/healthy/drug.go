package healthy

import (
	"kindergarten-service-go/models/healthy"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/validation"
)

// 喂药申请
type DrugController struct {
	beego.Controller
}

// URLMapping ...
func (c *DrugController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
}

// @Title 添加喂药申请
// @Description 添加喂药申请
// @Param   student_id     			formData    int  	true        "学生ID"
// @Param   drug     				formData    string  true        "药品"
// @Param   explain     			formData    string  true        "用量说明"
// @Param   symptom     			formData    string  true        "症状"
// @Param   user_id     			formData    int  	true        "用户ID"
// @Success 0 {int} models.Drug.Id
// @Failure 1001 补全信息
// @Failure 1003 保存失败
// @router / [post]
func (c *DrugController) Post() {
	baby_id, _ := c.GetInt("baby_id")
	drug := c.GetString("drug")
	explain := c.GetString("explain")
	symptom := c.GetString("symptom")
	user_id, _ := c.GetInt("user_id")
	url := c.GetString("url")

	valid := validation.Validation{}
	valid.Required(baby_id, "student_id").Message("宝宝ID不能为空")
	valid.Required(drug, "drug").Message("药品不能为空")
	valid.Required(explain, "explain").Message("用量说明不能为空")
	valid.Required(symptom, "symptom").Message("症状不能为空")
	valid.Required(user_id, "user_id").Message("用户ID不能为空")

	if valid.HasErrors() {
		c.Data["json"] = JSONStruct{"error", 1001, "", valid.Errors[0].Message}

		c.ServeJSON()
		c.StopRun()
	}
	w := healthy.Drug{
		Drug:    drug,
		Explain: explain,
		Symptom: symptom,
		UserId:  user_id,
		Url:     url,
	}
	if err := w.Save(baby_id); err == nil {

		c.Data["json"] = JSONStruct{"success", 0, "", "申请成功"}
	} else if err == orm.ErrNoRows {

		c.Data["json"] = JSONStruct{"success", 0, "", "宝宝不存在班级，无法申请"}
	} else {

		c.Data["json"] = JSONStruct{"error", 1003, "", "申请失败"}
	}

	c.ServeJSON()

}

// GetAll ...
// @Title GetAll
// @Description 喂药申请列表
// @Param	page			query	int	false		"第几页"
// @Param	kindergarten_id	query	int	true		"幼儿园ID"
// @Param	per_page		query	int	true		"页数"
// @Param	class_id		query	int	true		"班级ID"
// @Success 0 {object} 		healthy.Drug
// @Failure 1001 		参数不能为空
// @Failure 1005 		获取失败
// @router / [get]
func (c *DrugController) GetAll() {
	var f *healthy.Drug

	kindergarten_id, _ := c.GetInt("kindergarten_id")
	role, _ := c.GetInt("role")
	class_id, _ := c.GetInt("class_id")
	types, _ := c.GetInt("types")
	perPage, _ := c.GetInt("per_page")
	page, _ := c.GetInt("page")

	//验证参数是否为空
	valid := validation.Validation{}
	valid.Required(kindergarten_id, "kindergarten_id").Message("幼儿园ID不能为空")
	valid.Required(role, "role").Message("用户身份不能为空")
	if valid.HasErrors() {
		c.Data["json"] = JSONStruct{"error", 1001, struct{}{}, valid.Errors[0].Message}
		c.ServeJSON()
		c.StopRun()
	}
	if works, err := f.GetAll(page, perPage, kindergarten_id, class_id, role, types); err == nil {
		c.Data["json"] = JSONStruct{"success", 0, works, "获取成功"}
	} else {
		c.Data["json"] = JSONStruct{"error", 1005, err, "获取失败"}
	}

	c.ServeJSON()

}

// GetAll ...
// @Title GetAll
// @Description 详情
// @Param	kindergarten_id	query	int	true		"幼儿园ID"
// @Success 0 {object} 	healthy.Drug
// @Failure 1001 		参数不能为空
// @Failure 1003 		获取失败
// @router /:id
func (c *DrugController) DrugInfo() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	v := healthy.DrugInfo(id)
	if v == nil {
		c.Data["json"] = JSONStruct{"error", 1005, nil, "获取失败"}
	} else {
		c.Data["json"] = JSONStruct{"success", 0, v, "获取成功"}
	}
	c.ServeJSON()
}
