package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/validation"
	"kindergarten-service-go/models"
	"time"
)

//考勤
type AttCtl struct {
	BaseController
}

// @Title 某班待考勤学生列表
// @Description 某班待考勤学生列表(教师考勤首页使用)
// @Param	cid		query int	true	"班ID"
// @Success 200		success
// @Failure 403
// @router /stus [get]
func (this *AttCtl) GotStdsByTeaID() {
	defer this.ServeJSON()
	cid, _ := this.GetInt("cid")
	beego.Info("cid:", cid)
	r := models.GotStdsByTeaID(cid)
	beego.Info(r)
	if len(r) > 0 {
		this.Data["json"] = JSONStruct{"success", 0, r, "成功"}
	} else {
		this.Data["json"] = JSONStruct{"success", 0, nil, "无结果"}
	}
}

// @Title 异常考勤
// @Description 异常考勤内容
// @Param	cid		query int	true	"班ID"
// @Success 200		success
// @Failure 403
// @router /abn [get]
func (this *AttCtl) GotAbnDetail() {
	defer this.ServeJSON()
	cid, _ := this.GetInt("cid")
	beego.Info("cid:", cid)
	r := models.GotAbnDtl(cid)
	beego.Info(r)
	if len(r) > 0 {
		this.Data["json"] = JSONStruct{"success", 0, r, "成功"}
	} else {
		this.Data["json"] = JSONStruct{"success", 0, nil, "无结果"}
	}
}

// @Title 执行考勤
// @Description 教师对学生进行考勤操作（教师考勤按钮使用）
// @Param	sid			query int	true	"学生ID"
// @Param	status		query int	false	"考勤状态；0：有效：1：无效，缺省值：0"
// @Success 200			success
// @Failure 403
// @router /toone [get]
func (this *AttCtl) ToAtt() {
	defer this.ServeJSON()
	sid, _ := this.GetInt("sid")
	beego.Info("sid:", sid)
	status, _ := this.GetInt("status", 0)
	beego.Info("status:", status)
	e := models.Att(sid, status)
	if e == nil {
		this.Data["json"] = JSONStruct{"success", 0, nil, "成功"}
	} else {
		this.Data["json"] = JSONStruct{"failure", 0, e, "失败"}
	}
}

// @Title 一键入园
// @Description 批量考勤（教师考勤按钮使用）
// @Param	cid			query int	true	"班ID"
// @Success 200			success
// @Failure 403
// @router /toall [get]
func (this *AttCtl) ToAll() {
	defer this.ServeJSON()
	cid, _ := this.GetInt("cid")
	e := models.ToAll(cid)
	if e == nil {
		this.Data["json"] = JSONStruct{"success", 0, nil, "成功"}
	} else {
		this.Data["json"] = JSONStruct{"failure", 0, nil, e.Error()}
	}
}

// @Title 请假
// @Description 请假（用户使用）{"Sid":33,"Applicant":"lchj","Type":1,"Reason":"xxxxx","Beg":"2018-01-01T13:13:13Z","End":"2018-12-05T07:07:07Z"}
// @Param	body	body	models.Leave	true	"json"
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

// @Title  考勤规则
// @Description 写入考勤规则(园长使用)
// @Param	body	body	models.AttendanceRule	true	"json"
// @Success 200	success
// @Failure 403
// @router /rule [post]
func (this *AttCtl) Rule() {
	defer this.ServeJSON()
	var o models.AttendanceRule
	if err := json.Unmarshal(this.Ctx.Input.RequestBody, &o); err == nil {
		beego.Debug("INPUT:", o)
		valid := validation.Validation{}
		valid.Required(o.Kid, "kid").Message("幼儿园ID不可空")
		valid.Required(o.Mbeg, "mbeg").Message("上午考勤开始时间不可空")
		valid.Required(o.Mend, "mend").Message("上午考勤结束时间不可空")
		valid.Required(o.Abeg, "abeg").Message("下午考勤开始时间不可空")
		valid.Required(o.Aend, "aend").Message("下午考勤结束时间不可空")
		valid.Required(o.Aend, "days").Message("作用日")
		if valid.HasErrors() {
			this.Data["json"] = JSONStruct{"error", 1001, nil, valid.Errors[0].Message}
			return
		}
		e := models.AttRule(o)
		beego.Debug(e)
		if e == nil {
			this.Data["json"] = JSONStruct{"success", 0, nil, "成功"}
		} else {
			this.Data["json"] = JSONStruct{"failure", 0, e, "失败"}
		}
	} else {
		this.Data["json"] = JSONStruct{"error", 1001, err.Error(), "字段必须为json格式"}
	}
}

// @Title 考勤统计
// @Description 按日年级对班统计考勤（园长使用）
// @Param	gid			query int	true	"年级ID"
// @Param	day			query string	true	"2016-01-02"
// @Success 200			success
// @Failure 403
// @router /count [get]
func (this *AttCtl) Count() {
	defer this.ServeJSON()
	gid, _ := this.GetInt("gid")
	day := this.GetString("day", time.Now().Format("2006-01-02"))
	r := models.CountByGrade(day, gid)
	beego.Info(r)
	if len(r) > 0 {
		this.Data["json"] = JSONStruct{"success", 0, r, "成功"}
	} else {
		this.Data["json"] = JSONStruct{"success", 0, nil, "无考勤"}
	}
}

// @Title 考勤详情
// @Description 某日某班的考勤详情（园长使用）
// @Param	cid			query int	true	"班级ID"
// @Param	day			query string	true	"2016-01-02"
// @Success 200			success
// @Failure 403
// @router /dtl [get]
func (this *AttCtl) AttDtl() {
	defer this.ServeJSON()
	cid, _ := this.GetInt("cid")
	day := this.GetString("day", time.Now().Format("2006-01-02"))
	r := models.GotAttsByDayAndCls(day, cid)
	beego.Info(r)
	if len(r) > 0 {
		this.Data["json"] = JSONStruct{"success", 0, r, "成功"}
	} else {
		this.Data["json"] = JSONStruct{"success", 0, nil, "无考勤"}
	}
}
