package admin

import (
	"encoding/json"
	"fmt"
	"kindergarten-service-go/models"
	"math/rand"
	"strconv"
	"time"

	"github.com/hprose/hprose-golang/rpc"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/validation"
)

//教师
type TeacherController struct {
	BaseController
}

type UserService struct {
	GetOne   func(string) (int, error)
	GetUK    func(string) error
	Encrypt  func(string) string
	Test     func() string
	CreateUK func(userId int, kindergartenId int, role int) (int64, error)
	Create   func(phone string, name string, password string, kindergartenId int, role int) (interface{}, error)
}

type inviteTeacher struct {
	Name           string `json:"name"`
	Phone          string `json:"phone"`
	Role           int    `json:"role"`
	KindergartenId int    `json:"kindergarten_id"`
}

type OnemoreService struct {
	Test func() string
	Send func(phone string, text string) (map[string]interface{}, error)
}

// GetTeacher ...
// @Title 全部教师列表
// @Description 全部教师列表
// @Param	kindergarten_id       query	int	     true		"幼儿园ID"
// @Param	status                query	int	     false		"状态"
// @Param	search                query	int	     false		"搜索条件"
// @Param	page                  query	int	     false		"页数"
// @Param	per_page              query	int	     false		"每页显示条数"
// @Success 200 {object} models.Teacher
// @Failure 403
// @router / [get]
func (c *TeacherController) GetTeacher() {
	search := c.GetString("search")
	prepage, _ := c.GetInt("per_page", 20)
	page, _ := c.GetInt("page")
	kindergarten_id, _ := c.GetInt("kindergarten_id")
	status, _ := c.GetInt("status", -1)
	valid := validation.Validation{}
	valid.Required(kindergarten_id, "kindergarten_id").Message("幼儿园编号不能为空")
	if valid.HasErrors() {
		c.Data["json"] = JSONStruct{"error", 1001, nil, valid.Errors[0].Message}
		c.ServeJSON()
	} else {
		v := models.GetTeacher(kindergarten_id, status, search, page, prepage)
		if v == nil {
			c.Data["json"] = JSONStruct{"error", 1005, nil, "获取失败"}
		} else {
			c.Data["json"] = JSONStruct{"success", 0, v, "获取成功"}
		}
		c.ServeJSON()
	}
}

// GetClass ...
// @Title 班级列表
// @Description 班级列表
// @Param	kindergarten_id       query	int	     true		"幼儿园ID"
// @Param	class_type            query	int	     true		"班级类型"
// @Param	page                  query	int	     false		"页数"
// @Param	per_page              query	int	     false		"每页显示条数"
// @Success 200 {object} models.Teacher
// @Failure 403
// @router /class [get]
func (c *TeacherController) GetClass() {
	prepage, _ := c.GetInt("per_page", 20)
	page, _ := c.GetInt("page")
	kindergarten_id, _ := c.GetInt("kindergarten_id")
	class_type, _ := c.GetInt("class_type")
	v := models.GetClass(kindergarten_id, class_type, page, prepage)
	if v == nil {
		c.Data["json"] = JSONStruct{"error", 1005, nil, "获取失败"}
	} else {
		c.Data["json"] = JSONStruct{"success", 0, v, "获取成功"}
	}
	c.ServeJSON()
}

// Delete ...
// @Title Delete
// @Description 删除教师
// @Param	teacher_id		path 	int	true		"教师ID"
// @Param	status		    path 	int	true		"状态(status 0:未分班 2:离职)"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (c *TeacherController) Delete() {
	class_type, _ := c.GetInt("class_type")
	status, _ := c.GetInt("status")
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	v := models.DeleteTeacher(id, status, class_type)
	if v == nil {
		c.Data["json"] = JSONStruct{"error", 1004, nil, "删除失败"}
	} else {
		c.Data["json"] = JSONStruct{"success", 0, nil, "删除成功"}
	}
	c.ServeJSON()
}

// GetTeacherInfo ...
// @Title Get Teacher Info
// @Description 教师详情
// @Param	teacher_id       query	int	 true		"教师编号"
// @Success 200 {object} models.Teacher
// @Failure 403 :教师编号为空
// @router /:id [get]
func (c *TeacherController) GetTeacherInfo() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	v := models.GetTeacherInfo(id)
	if v == nil {
		c.Data["json"] = JSONStruct{"error", 1005, nil, "获取失败"}
	} else {
		c.Data["json"] = JSONStruct{"success", 0, v, "获取成功"}
	}
	c.ServeJSON()
}

// Put ...
// @Title 编辑教师
// @Description 编辑教师
// @Param	id		    path 	int	               true		    "教师编号"
// @Param	body		body 	models.Animation	true		"param(json)"
// @Success 200 {object} models.Animation
// @Failure 403 :id is not int
// @router /:id [put]
func (c *TeacherController) Put() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	v := models.Teacher{Id: id}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		v := models.UpdateTeacher(&v)
		if v == nil {
			c.Data["json"] = JSONStruct{"error", 1003, err.Error(), "编辑失败"}
		} else {
			c.Data["json"] = JSONStruct{"success", 0, nil, "编辑成功"}
		}
		c.ServeJSON()
	} else {
		c.Data["json"] = JSONStruct{"error", 1001, err.Error(), "字段必须为json格式"}
		c.ServeJSON()
	}
}

// Post ...
// @Title 教师-录入信息
// @Description 教师-录入信息
// @Param	body		body 	models.Animation	true		"json"
// @Success 201 {int} models.Animation
// @Failure 403 body is empty
// @router / [post]
func (c *TeacherController) Post() {
	var v models.Teacher
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		valid := validation.Validation{}
		valid.Required(v.KindergartenId, "KindergartenId").Message("幼儿园编号不能为空")
		valid.Required(v.UserId, "UserId").Message("用户编号不能为空")
		valid.Required(v.Birthday, "Birthday").Message("出生年月日不能为空")
		valid.Required(v.Name, "Name").Message("用户名不能为空")
		valid.Required(v.Number, "Number").Message("教职工编号不能为空")
		valid.Required(v.NationOrReligion, "NationOrReligion").Message("民族或宗教不能为空")
		valid.Required(v.NativePlace, "NativePlace").Message("籍贯不能为空")
		valid.Required(v.EnterJobTime, "EnterJobTime").Message("参加工作时间不能为空")
		valid.Required(v.Address, "Address").Message("住址不能为空")
		valid.Required(v.EmergencyContact, "EmergencyContact").Message("紧急联系人不能为空")
		valid.Required(v.EmergencyContactPhone, "EmergencyContactPhone").Message("紧急联系人电话不能为空")
		valid.Required(v.Source, "Source").Message("来源不能为空")
		valid.Required(v.TeacherCertificationNumber, "TeacherCertificationNumber").Message("教师认证编号不能为空")
		valid.Required(v.Phone, "Phone").Message("手机号不能为空")
		valid.Required(v.EnterGardenTime, "EnterGardenTime").Message("进入本园时间不能为空")
		if valid.HasErrors() {
			c.Data["json"] = JSONStruct{"error", 1001, nil, valid.Errors[0].Message}
			c.ServeJSON()
		} else {
			v := models.AddTeacher(&v)
			if v == nil {
				c.Data["json"] = JSONStruct{"error", 1003, err.Error(), "保存失败"}
			} else {
				c.Data["json"] = JSONStruct{"success", 0, nil, "保存成功"}
			}
			c.ServeJSON()
		}
	} else {
		c.Data["json"] = JSONStruct{"error", 1001, err.Error(), "字段必须为json格式"}
		c.ServeJSON()
	}
}

// RemoveTeacher ...
// @Title 移除教师
// @Description 移除教师
// @Param	teacher_id		    path 	    int	true		    "教师ID"
// @Param	class_id		    path 	    int	true		    "班级ID"
// @Success 200 {string} delete success!
// @Failure 403 teacher_id is empty
// @router /remove [delete]
func (c *TeacherController) RemoveTeacher() {
	teacher_id, _ := c.GetInt("teacher_id")
	class_id, _ := c.GetInt("class_id")
	valid := validation.Validation{}
	valid.Required(teacher_id, "teacher_id").Message("教师ID不能为空")
	valid.Required(class_id, "class_id").Message("班级ID不能为空")
	if valid.HasErrors() {
		c.Data["json"] = JSONStruct{"error", 1001, nil, valid.Errors[0].Message}
		c.ServeJSON()
	} else {
		v := models.RemoveTeacher(teacher_id, class_id)
		if v == nil {
			c.Data["json"] = JSONStruct{"error", 1004, nil, "移除失败"}
		} else {
			c.Data["json"] = JSONStruct{"success", 0, nil, "移除成功"}
		}
		c.ServeJSON()
	}

}

// Invite ...
// @Title 邀请教师/批量邀请
// @Description 邀请教师/批量邀请
// @Param	phone		        body 	string	 true		"手机号"
// @Param	name		            body 	string   true		"姓名"
// @Param	role  		        body 	int  	 true		"身份"
// @Param	kindergarten_id		body 	int   	 true		"幼儿园ID"
// @Success 201 {int} models.Animation
// @Failure 403 body is empty
// @router /invite [post]
func (c *TeacherController) Invite() {
	var User *UserService
	var Onemore *OnemoreService
	var password string
	var text string
	teacher := c.GetString("teacher")
	var t []inviteTeacher
	json.Unmarshal([]byte(teacher), &t)
	valid := validation.Validation{}
	valid.Required(teacher, "teacher").Message("教师信息不能为空")
	if valid.HasErrors() {
		c.Data["json"] = JSONStruct{"error", 1001, nil, valid.Errors[0].Message}
		c.ServeJSON()
	} else {
		//rpc服务
		client := rpc.NewHTTPClient(beego.AppConfig.String("ONE_MORE_USER_SERVER"))
		client.UseService(&User)
		client = rpc.NewHTTPClient(beego.AppConfig.String("ONE_MORE_SMS_SERVER"))
		client.UseService(&Onemore)
		//获取用户关联表
		for _, value := range t {
			err := User.GetUK(value.Phone)
			if err == nil {
				c.Data["json"] = JSONStruct{"error", 1009, nil, "" + value.Phone + "已被邀请过"}
				c.ServeJSON()
			} else {
				//获取用户信息
				userId, _ := User.GetOne(value.Phone)
				if userId != 0 {
					User.CreateUK(userId, value.KindergartenId, value.Role)
				} else {
					//生成六位验证码
					rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
					vcode := fmt.Sprintf("%06v", rnd.Int31n(1000000))
					//发送短信
					text = "【蓝天白云】您已通过系统成功注册蓝天白云平台账号，您的账号为：" + value.Phone + "（手机号），密码为：" + vcode + "，请您登陆APP进行密码修改。"
					//密码加密
					password = User.Encrypt(vcode)
					res, err := Onemore.Send(value.Phone, text)
					if err == nil {
						if int(res["code"].(float64)) == 0 {
							_, err = User.Create(value.Phone, value.Name, password, value.KindergartenId, value.Role)
							if err == nil {
								c.Data["json"] = JSONStruct{"success", 0, nil, res["msg"].(string)}
								c.ServeJSON()
							}
						} else {
							c.Data["json"] = JSONStruct{"error", 1001, nil, res["msg"].(string)}
							c.ServeJSON()
						}
					} else {
						c.Data["json"] = JSONStruct{"error", 1001, err.Error(), "发送有误"}
						c.ServeJSON()
					}
				}
			}
		}
	}
}

// OrganizationalTeacher ...
// @Title 组织框架教师列表
// @Description 组织框架教师列表
// @Param	kindergarten_id       query	int	     true		"幼儿园ID"
// @Param	type                  query	int	     true		"年级组标识(1 年级组)"
// @Param	person                query	int	     true		"是否为负责人(1 负责 2 不是负责人)"
// @Success 200 {object} models.Teacher
// @Failure 403
// @router /organizational_teacher [get]
func (c *TeacherController) OrganizationalTeacher() {
	ty, _ := c.GetInt("type")
	person, _ := c.GetInt("person")
	class_id, _ := c.GetInt("class_id")
	kindergarten_id, _ := c.GetInt("kindergarten_id")
	valid := validation.Validation{}
	valid.Required(person, "person").Message("身份不能为空")
	valid.Required(kindergarten_id, "kindergarten_id").Message("幼儿园编号不能为空")
	if valid.HasErrors() {
		c.Data["json"] = JSONStruct{"error", 1001, nil, valid.Errors[0].Message}
		c.ServeJSON()
	} else {
		v := models.OrganizationalTeacher(kindergarten_id, ty, person, class_id)
		if v == nil {
			c.Data["json"] = JSONStruct{"error", 1005, nil, "获取失败"}
		} else {
			c.Data["json"] = JSONStruct{"success", 0, v, "获取成功"}
		}
		c.ServeJSON()
	}
}
