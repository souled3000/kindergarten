package controllers

import (
	"encoding/json"
	"kindergarten-service-go/models"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/validation"
)

//考勤
type AttCtl struct {
	BaseController
}

// @Title 获取教师的学生
// @Description 获取教师管理的学生(教师考勤首页使用)
// @Param	tid		query int	true	"教师ID"
// @Success 200		success
// @Failure 403
// @router /stus [get]
func (this *AttCtl) GotStdsByTeaID() {
	defer this.ServeJSON()
	tid, _ := this.GetInt("tid")
	beego.Info("tid:", tid)
	r := models.GotStdsByTeaID(tid)
	beego.Info(r)
	if len(r) > 0 {
		this.Data["json"] = JSONStruct{"success", 0, r, "获取学生成功"}
	} else {
		this.Data["json"] = JSONStruct{"success", 0, nil, "教师下无学生"}
	}
}

// @Title 异常考勤
// @Description 异常考勤内容
func (this *AttCtl) GotAbnDetail() {
	defer this.ServeJSON()

}

// @Title 执行考勤
// @Description 教师对学生进行考勤操作（教师考勤按钮使用）
// @Param	sid			query int	true	"学生ID"
// @Param	status		query int	false	"考勤状态；0：有效：1：无效，缺省值：0"
// @Success 200			success
// @Failure 403
// @router /toatt [get]
func (this *AttCtl) ToAtt() {
	defer this.ServeJSON()
	sid, _ := this.GetInt("sid")
	beego.Info("sid:", sid)
	status, _ := this.GetInt("status", 0)
	beego.Info("status:", status)
	e := models.Att(sid, status)
	if e == nil {
		this.Data["json"] = JSONStruct{"success", 0, nil, "考勤成功"}
	} else {
		this.Data["json"] = JSONStruct{"failure", 0, e, "考勤失败"}
	}
}

// @Title 请假
// @Description 请假（用户使用）{"Sid":15,"Applicant":"lchj","Type":1,"Reason":"xxxxx","Beg":"2018-01-01T13:13:13Z","End":"2018-04-05T07:07:07Z"}
// @Param	body	body	models.Attendance	true	"json"
// @Success 200	success
// @Failure 403
// @router /askleave [post]
func (this *AttCtl) AskForLeave() {
	defer this.ServeJSON()
	var o models.Leave
	if err := json.Unmarshal(this.Ctx.Input.RequestBody, &o); err == nil {
		beego.Debug("INPUT:", o)
		valid := validation.Validation{}
		valid.Required(o.Applicant, "Applicant").Message("申请人不能为空")
		valid.Required(o.Reason, "Reason").Message("请假原因不能为空")
		valid.Required(o.Beg, "Beg").Message("请假开始时间不能为空")
		valid.Required(o.End, "End").Message("请假结束时间不能为空")
		valid.Required(o.Type, "Type").Message("类型不能为空")
		if o.Beg.Unix() > o.End.Unix() {
			valid.AddError("BGTE", "开始时间大于结束时间")
		}
		if valid.HasErrors() {
			this.Data["json"] = JSONStruct{"error", 1001, nil, valid.Errors[0].Message}
			return
		}
		e := models.AskForLeave(o)
		beego.Debug(e)
		if e == nil {
			this.Data["json"] = JSONStruct{"success", 0, nil, "请假成功"}
		} else {
			this.Data["json"] = JSONStruct{"failure", 0, e, "请假失败"}
		}
	} else {
		this.Data["json"] = JSONStruct{"error", 1001, err.Error(), "字段必须为json格式"}
	}
}

// @Title  记录考勤规则
// @Description 写入考勤规则(园长使用)
// @Param	body	body	models.Attendance	true	"json"
// @Success 200		success
// @Failure 403
// @router /stus [get]
func (this *AttCtl) AttRule() {
	defer this.ServeJSON()
	var o models.Leave
	if err := json.Unmarshal(this.Ctx.Input.RequestBody, &o); err == nil {
		beego.Debug("INPUT:", o)
		valid := validation.Validation{}
		valid.Required(o.Applicant, "Applicant").Message("申请人不能为空")
		valid.Required(o.Reason, "Reason").Message("请假原因不能为空")
		valid.Required(o.Beg, "Beg").Message("请假开始时间不能为空")
		valid.Required(o.End, "End").Message("请假结束时间不能为空")
		valid.Required(o.Type, "Type").Message("类型不能为空")
		if o.Beg.Unix() > o.End.Unix() {
			valid.AddError("BGTE", "开始时间大于结束时间")
		}
		if valid.HasErrors() {
			this.Data["json"] = JSONStruct{"error", 1001, nil, valid.Errors[0].Message}
			return
		}
		e := models.AskForLeave(o)
		beego.Debug(e)
		if e == nil {
			this.Data["json"] = JSONStruct{"success", 0, nil, "请假成功"}
		} else {
			this.Data["json"] = JSONStruct{"failure", 0, e, "请假失败"}
		}
	} else {
		this.Data["json"] = JSONStruct{"error", 1001, err.Error(), "字段必须为json格式"}
	}
}
