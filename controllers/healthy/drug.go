package healthy

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/validation"
	"fmt"
	"kindergarten-service-go/models/healthy"
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
	student_id, _:= c.GetInt("student_id")
	drug := c.GetString("drug")
	explain := c.GetString("explain")
	symptom := c.GetString("symptom")
	user_id, _:= c.GetInt("user_id")
	url := c.GetString("url")

	valid := validation.Validation{}
	valid.Required(student_id, "student_id").Message("学生ID不能为空")
	valid.Required(drug,"drug").Message("药品不能为空")
	valid.Required(explain,"explain").Message("用量说明不能为空")
	valid.Required(symptom,"symptom").Message("症状不能为空")
	valid.Required(user_id,"user_id").Message("用户ID不能为空")
	valid.Required(url,"url").Message("图片不能为空")
	if valid.HasErrors() {
		c.Data["json"] = JSONStruct{"error", 1001, "", valid.Errors[0].Message}

		c.ServeJSON()
	}
	w := healthy.Drug{
		StudentId:student_id,
		Drug:drug,
		Explain:explain,
		Symptom:symptom,
		UserId:user_id,
		Url:url,
	}
	if err := w.Save(); err == nil {
		c.Data["json"] = JSONStruct{"success", 0, "", "申请成功"}
	} else {
		fmt.Println(err)
		c.Data["json"] = JSONStruct{"error", 1003, "", "申请失败"}
	}

	c.ServeJSON()

}
