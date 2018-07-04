package task

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/validation"
	"time"
	"kindergarten-service-go/models/task"
	"encoding/json"
)

type WorkTaskController struct {
	beego.Controller
}

type JSONStruct struct {
	Status string      `json:"status"`
	Code   int         `json:"code"`
	Result interface{} `json:"result"`
	Msg    string      `json:"msg"`
}

func (c *WorkTaskController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("Get", c.Get)
}

// @Title 发布任务
// @Description 发布任务
// @Param   title     				formData    string  true        "标题"
// @Param   describe     			formData    string  true        "描述"
// @Param   deadline     			formData    time    true        "截止日期"
// @Param   save_folder_id     		formData    int  	false       "存入文件夹ID"
// @Param   save_folder_name     	formData    string  false       "存入文件夹名称"
// @Param   publisher     			formData    int  	true        "发布人ID"
// @Param   publisher_name     		formData    string  true        "发布人名称"
// @Param   operator     			formData    string  true        "执行人"
// @Param   cc     					formData    string  false       "抄送人"
// @Success 0 {int} models.Feedback.Id
// @Failure 1001 参数验证
// @Failure 1003 发布失败
// @router / [post]
func (c *WorkTaskController) Post() {
	title := c.GetString("title")
	describe := c.GetString("describe")
	deadlineS := c.GetString("deadline")
	saveFolderId, _ := c.GetInt("save_folder_id")
	saveFolderName := c.GetString("save_folder_name")
	publisher, _ := c.GetInt("publisher")
	publisherName := c.GetString("publisher_name")
	operatorS := c.GetString("operator")
	ccS := c.GetString("cc")

	valid := validation.Validation{}
	valid.Required(title, "title").Message("标题不能为空")
	valid.Required(describe, "describe").Message("描述不能为空")
	valid.Required(deadlineS, "deadline").Message("截止时间不能为空")
	valid.Required(publisher, "publisher").Message("发布人不能为空")
	valid.Required(operatorS, "operator").Message("执行人不能为空")
	deadline, err := time.Parse("2006-01-02 15:04:05", deadlineS)
	if err != nil {
		valid.SetError("deadline", "截止时间格式不正确")
	}
	var operator []map[string]interface{}
	if err := json.Unmarshal([]byte(operatorS), &operator); err != nil {
		valid.SetError("operator", "执行人格式不正确")
	}
	var cc []map[string]interface{}
	if ccS != "" {
		if err := json.Unmarshal([]byte(ccS), &cc); err != nil {
			valid.SetError("cc", "抄送人格式不正确")
		}
	}
	if valid.HasErrors() {
		c.Data["json"] = JSONStruct{"error", 1001, "", valid.Errors[0].Message}

		c.StopRun()
		c.ServeJSON()
	}

	wt := task.WorkTasks{
		Title:title,
		Describe:describe,
		SaveFolderId:saveFolderId,
		SaveFolderName:saveFolderName,
		Publisher:publisher,
		PublisherName:publisherName,
		Deadline:deadline,
		TaskNum:len(operator),
	}

	if err := wt.Save(operator, cc); err == nil {
		c.Data["json"] = JSONStruct{"success", 0, "", "发布成功"}
	} else {
		c.Data["json"] = JSONStruct{"error", 1003, "", "发布失败"}
	}

	c.ServeJSON()
}

// @Title 获取任务列表
// @Description 获取任务列表
// @Success 0 {int} models.Feedback.Id
// @Failure 1005 获取失败
// @router / [get]
func (c *WorkTaskController) Get() {
	var wt task.WorkTasks

	if res, err := wt.Get(); err == nil {
		c.Data["json"] = JSONStruct{"success", 0, res, "获取成功"}
	} else {
		c.Data["json"] = JSONStruct{"error", 1005, "", "获取失败"}
	}

	c.ServeJSON()
}

// @Title 获取任务详情
// @Description 获取任务详情
// @Param   id     			path    int  	true        "任务ID"
// @Success 0 {int} models.Feedback.Id
// @Failure 1005 获取失败
// @router /:id [get]
func (c *WorkTaskController) GetInfo() {
	var id int
	c.Ctx.Input.Bind(&id, ":id")
	wt := task.WorkTasks{Id:id}

	if res, err := wt.GetInfoById(); err == nil {
		c.Data["json"] = JSONStruct{"success", 0, res, "获取成功"}
	} else {
		c.Data["json"] = JSONStruct{"error", 1005, "", "获取失败"}
	}

	c.ServeJSON()
}

func (c *WorkTaskController) Complete() {

}