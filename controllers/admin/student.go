package admin

import (
	"encoding/json"
	"fmt"
	"kindergarten-service-go/models"
	"math/rand"
	"strconv"
	"time"

	"github.com/astaxie/beego/validation"
	"github.com/hprose/hprose-golang/rpc"

	"github.com/astaxie/beego"
)

//学生
type StudentController struct {
	BaseController
}

// GetStudent ...
// @Title 学生列表
// @Description 学生列表
// @Param	kindergarten_id       query	int	     true		"幼儿园ID"
// @Param	status                query	int	     false		"状态"
// @Param	search                query	int	     false		"搜索条件"
// @Param	page                  query	int	     false		"页数"
// @Param	per_page              query	int	     false		"每页显示条数"
// @Success 200 {object} models.Student
// @Failure 403
// @router / [get]
func (c *StudentController) GetStudent() {
	search := c.GetString("search")
	prepage, _ := c.GetInt("per_page", 20)
	page, _ := c.GetInt("page")
	kindergarten_id, _ := c.GetInt("kindergarten_id")
	status, _ := c.GetInt("status", -1)
	v := models.GetStudent(kindergarten_id, status, search, page, prepage)
	if v == nil {
		c.Data["json"] = JSONStruct{"error", 1005, nil, "获取失败"}
	} else {
		c.Data["json"] = JSONStruct{"success", 0, v, "获取成功"}
	}
	c.ServeJSON()
}

// GetStudentClass ...
// @Title 学生班级搜索
// @Description 学生班级搜索
// @Param	kindergarten_id       query	int	     true		"幼儿园ID"
// @Param	class_type            query	int	     true		"班级类型"
// @Param	page                  query	int	     false		"页数"
// @Param	per_page              query	int	     false		"每页显示条数"
// @Success 200 {object} models.Student
// @Failure 403
// @router /class [get]
func (c *StudentController) GetStudentClass() {
	prepage, _ := c.GetInt("per_page", 20)
	page, _ := c.GetInt("page")
	kindergarten_id, _ := c.GetInt("kindergarten_id")
	class_type, _ := c.GetInt("class_type")
	v := models.GetStudentClass(kindergarten_id, class_type, page, prepage)
	if v == nil {
		c.Data["json"] = JSONStruct{"error", 1005, nil, "获取失败"}
	} else {
		c.Data["json"] = JSONStruct{"success", 0, v, "获取成功"}
	}
	c.ServeJSON()
}

// RemoveStudent ...
// @Title RemoveStudent
// @Description 移除学生
// @Param	student_id		path 	    int	true		"学生ID"
// @Param	class_id		    path 	    int	true		"班级ID"
// @Success 200 {string} delete success!
// @Failure 403 student_id is empty
// @router /remove [delete]
func (c *StudentController) RemoveStudent() {
	student_id, _ := c.GetInt("student_id")
	class_id, _ := c.GetInt("class_id")
	v := models.RemoveStudent(class_id, student_id)
	if v == nil {
		c.Data["json"] = JSONStruct{"error", 1004, nil, "移除失败"}
	} else {
		c.Data["json"] = JSONStruct{"success", 0, nil, "移除成功"}
	}
	c.ServeJSON()
}

// GetStudentInfo ...
// @Title Get Student Info
// @Description 学生详情
// @Param	student_id       query	int	 true		"学生编号"
// @Success 200 {object} models.Student
// @Failure 403 :学生编号为空
// @router /:id [get]
func (c *StudentController) GetOne() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	v, err := models.GetStudentInfo(id)
	if err != nil {
		c.Data["json"] = JSONStruct{"error", 1005, nil, err.Error()}
	} else {
		c.Data["json"] = JSONStruct{"success", 0, v, "获取成功"}
	}
	c.ServeJSON()
}

// UpdateStudent ...
// @Title 编辑学生
// @Description 编辑学生
// @Param	id		    path 	int	               true		    "学生编号"
// @Param	body		body 	models.Student	       true		"param(json)"
// @Success 200 {object} models.Student
// @Failure 403 :id is not int
// @router /:id [put]
func (c *StudentController) UpdateStudent() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	student := c.GetString("student")
	kinship := c.GetString("kinship")
	valid := validation.Validation{}
	valid.Required(student, "student").Message("学生信息不能为空")
	valid.Required(kinship, "kinship").Message("亲属信息不能为空")
	if valid.HasErrors() {
		c.Data["json"] = JSONStruct{"error", 1001, nil, valid.Errors[0].Message}
		c.ServeJSON()
	} else {
		_, err := models.UpdateStudent(id, student, kinship)
		if err != nil {
			c.Data["json"] = JSONStruct{"error", 1003, nil, err.Error()}
		} else {
			c.Data["json"] = JSONStruct{"success", 0, nil, "编辑成功"}
		}
		c.ServeJSON()
	}
}

// Post ...
// @Title 学生-录入信息
// @Description 学生-录入信息
// @Param	body		body 	models.Animation	true		"json"
// @Success 201 {int} models.Student
// @Failure 403 body is empty
// @router / [post]
func (c *StudentController) Post() {
	student := c.GetString("student")
	kinship := c.GetString("kinship")
	valid := validation.Validation{}
	valid.Required(student, "student").Message("学生信息不能为空")
	valid.Required(kinship, "kinship").Message("亲属信息不能为空")
	if valid.HasErrors() {
		c.Data["json"] = JSONStruct{"error", 1001, nil, valid.Errors[0].Message}
		c.ServeJSON()
	} else {
		_, err := models.AddStudent(student, kinship)
		if err != nil {
			c.Data["json"] = JSONStruct{"error", 1003, nil, err.Error()}
		} else {
			c.Data["json"] = JSONStruct{"success", 0, nil, "保存成功"}
		}
		c.ServeJSON()
	}
}

// Invite ...
// @Title 邀请学生/批量邀请
// @Description 邀请学生/批量邀请
// @Param	phone		        body 	string	true		"手机号(json)"
// @Param	name		            body 	int   	true		"学生姓名(json)"
// @Param	role  		        body 	int  	true		"身份(json)"
// @Param	kindergarten_id		body 	int   	true		"幼儿园ID(json)"
// @Success 201 {int} models.Student
// @Failure 403 body is empty
// @router /invite [post]
func (c *StudentController) Invite() {
	var User *UserService
	var Onemore *OnemoreService
	var password string
	var text string
	student := c.GetString("student")
	var s []inviteTeacher
	json.Unmarshal([]byte(student), &s)
	valid := validation.Validation{}
	valid.Required(student, "student").Message("学生信息不能为空")
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
		for _, value := range s {
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
					_, err := User.Create(value.Phone, value.Name, password, value.KindergartenId, value.Role)
					if err == nil {
						res, err := Onemore.Send(value.Phone, text)
						if err == nil {
							if int(res["code"].(float64)) == 0 {
								c.Data["json"] = JSONStruct{"success", 0, nil, res["msg"].(string)}
								c.ServeJSON()
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
}

// DeleteStudent ...
// @Title DeleteStudent
// @Description 删除学生
// @Param	student_id		path 	int	true		"学生ID"
// @Param	status		    path 	int	true		"状态(status 0:未分班 2:离园)"
// @Param	type		        path 	int	true		"删除类型（type 0:学生离园 1:删除档案）"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (c *StudentController) DeleteStudent() {
	class_type, _ := c.GetInt("class_type")
	status, _ := c.GetInt("status")
	ty, _ := c.GetInt("type")
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	v := models.DeleteStudent(id, status, ty, class_type)
	if v == nil {
		c.Data["json"] = JSONStruct{"error", 1004, nil, "删除失败"}
	} else {
		c.Data["json"] = JSONStruct{"success", 0, nil, "删除成功"}
	}
	c.ServeJSON()
}
