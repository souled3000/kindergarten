package task

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/validation"
	"time"
	"kindergarten-service-go/models/task"
	"github.com/astaxie/beego/orm"
)

type WorkPlanController struct {
	beego.Controller
}

func (c *WorkPlanController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("Get", c.Get)
}

// @Title 创建工作计划
// @Description 创建工作计划
// @Param   content     		formData    string  true        "内容"
// @Param   plan_time     		formData    string  true        "计划时间"
// @Param   creator     		formData    int  	true        "创建人"
// @Success 0 {json} JSONStruct
// @Failure 1001 参数验证
// @Failure 1003 创建失败
// @router / [post]
func (c *WorkPlanController) Post() {
	content := c.GetString("content")
	planTimeS := c.GetString("plan_time")
	creator, _ := c.GetInt("creator")

	valid := validation.Validation{}
	valid.Required(content, "content").Message("内容不能为空")
	valid.Required(planTimeS, "plan_time").Message("计划时间不能为空")
	valid.Required(creator, "creator").Message("创建人不能为空")
	planTime, err := time.Parse("2006-01-02 15:04:05", planTimeS)
	if err != nil {
		valid.SetError("plan_time", "计划时间格式不正确")
	}
	if valid.HasErrors() {
		c.Data["json"] = JSONStruct{"error", 1001, "", valid.Errors[0].Message}

		c.ServeJSON()
		c.StopRun()
	}

	wp := task.WorkPlan{Content:content, PlanTime:planTime, Creator:creator}

	if _, err := wp.Save(); err == nil {
		c.Data["json"] = JSONStruct{"success", 0, "", "创建成功"}
	} else {
		c.Data["json"] = JSONStruct{"error", 1003, "", "创建失败"}
	}

	c.ServeJSON()
}

// @Title 获取工作计划
// @Description 获取工作计划
// @Param   u_id     		query    int  	true        "创建人"
// @Success 0 {json} JSONStruct
// @Failure 1005 获取失败
// @router / [get]
func (c *WorkPlanController) Get() {
	uId, _ := c.GetInt("u_id")

	wp := task.WorkPlan{Creator:uId}

	if res, err := wp.Get(); err == nil {
		c.Data["json"] = JSONStruct{"success", 0, res, "获取成功"}
	} else {
		c.Data["json"] = JSONStruct{"error", 1005, "", "获取失败"}
	}

	c.ServeJSON()
}

// @Title 删除工作计划
// @Description 删除工作计划
// @Param   id     		path    int  	true        "工作计划ID"
// @Success 0 {json} JSONStruct
// @Failure 1002 工作计划不存在
// @Failure 1004 删除失败
// @router /:id [delete]
func (c *WorkPlanController) Delete() {
	var id int
	c.Ctx.Input.Bind(&id, ":id")

	wp := task.WorkPlan{Id:id}

	if err := wp.Delete(); err == nil {
		c.Data["json"] = JSONStruct{"success", 0, "", "删除成功"}
	} else if err == orm.ErrNoRows {
		c.Data["json"] = JSONStruct{"error", 1002, "", "工作计划不存在"}
	} else {
		c.Data["json"] = JSONStruct{"error", 1004, "", "删除失败"}
	}

	c.ServeJSON()
}