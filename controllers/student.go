package controllers

import (
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
	beego.Controller
}

// URLMapping ...
func (c *StudentController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
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
	var prepage int = 20
	var page int
	var kindergarten_id int
	var status int
	var search string
	search = c.GetString("search")
	if v, err := c.GetInt("per_page"); err == nil {
		prepage = v
	}
	if v, err := c.GetInt("page"); err == nil {
		page = v
	}
	if v, err := c.GetInt("kindergarten_id"); err == nil {
		kindergarten_id = v
	}
	if v, err := c.GetInt("status", -1); err == nil {
		status = v
	}
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
	var prepage int = 20
	var page int
	var kindergarten_id int
	var class_type int
	if v, err := c.GetInt("per_page"); err == nil {
		prepage = v
	}
	if v, err := c.GetInt("page"); err == nil {
		page = v
	}
	if v, err := c.GetInt("kindergarten_id"); err == nil {
		kindergarten_id = v
	}
	if v, err := c.GetInt("class_type"); err == nil {
		class_type = v
	}
	v := models.GetStudentClass(kindergarten_id, class_type, page, prepage)
	if v == nil {
		c.Data["json"] = JSONStruct{"error", 1005, nil, "获取失败"}
	} else {
		c.Data["json"] = JSONStruct{"success", 0, v, "获取成功"}
	}
	c.ServeJSON()
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

// GetStudentInfo ...
// @Title Get Student Info
// @Description 学生详情
// @Param	student_id       query	int	 true		"学生编号"
// @Success 200 {object} models.Student
// @Failure 403 :学生编号为空
// @router /:id [get]
func (c *StudentController) GetStudentInfo() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	v := models.GetStudentInfo(id)
	if v == nil {
		c.Data["json"] = JSONStruct{"error", 1005, nil, "获取失败"}
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
		v := models.UpdateStudent(id, student, kinship)
		fmt.Println(v)
		if v == nil {
			fmt.Println(111)
			c.Data["json"] = JSONStruct{"error", 1003, nil, "编辑失败"}
		} else {
			c.Data["json"] = JSONStruct{"success", 0, nil, "编辑成功"}
			fmt.Println(222)
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
		l := models.AddStudent(student, kinship)
		if l == nil {
			c.Data["json"] = JSONStruct{"error", 1003, nil, "保存失败"}
		} else {
			c.Data["json"] = JSONStruct{"success", 0, nil, "保存成功"}
		}
		c.ServeJSON()
	}
}

// Invite ...
// @Title 邀请学生
// @Description 邀请学生
// @Param	phone		        body 	string	true		"手机号"
// @Param	name		        body 	int   	true		"学生姓名"
// @Param	role  		        body 	int  	true		"身份"
// @Param	kindergarten_id		 body 	int   	true		"幼儿园ID"
// @Success 201 {int} models.Student
// @Failure 403 body is empty
// @router /invite [post]
func (c *StudentController) Invite() {
	var User *UserService
	var Onemore *OnemoreService
	var password string
	var text string
	phone := c.GetString("phone")
	name := c.GetString("name")
	role, _ := c.GetInt("role")
	kindergartenId, _ := c.GetInt("kindergarten_id")
	valid := validation.Validation{}
	valid.Required(phone, "phone").Message("手机号不能为空")
	valid.Required(name, "name").Message("姓名不能为空")
	valid.Required(role, "role").Message("身份不能为空")
	valid.Required(kindergartenId, "kindergartenId").Message("幼儿园ID不能为空")
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
		err := User.GetUK(phone)
		if err == nil {
			c.Data["json"] = JSONStruct{"error", 1009, nil, "" + phone + "已被邀请过"}
			c.ServeJSON()
		} else {
			//获取用户信息
			userId, _ := User.GetOne(phone)
			if userId != 0 {
				User.CreateUK(userId, kindergartenId, role)
			} else {
				//生成六位验证码
				rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
				vcode := fmt.Sprintf("%06v", rnd.Int31n(1000000))
				//发送短信
				text = "【蓝天白云】您已通过系统成功注册蓝天白云平台账号，您的账号为：" + phone + "（手机号），密码为：" + vcode + "，请您登陆APP进行密码修改。"
				//密码加密
				password = User.Encrypt(vcode)
				_, err := User.Create(phone, name, password, kindergartenId, role)
				if err == nil {
					_, err := Onemore.Send(phone, text)
					if err == nil {
						c.Data["json"] = JSONStruct{"success", 0, nil, "发送成功"}
						c.ServeJSON()
					} else {
						c.Data["json"] = JSONStruct{"error", 1001, nil, "发送失败"}
						c.ServeJSON()
					}
				}
			}
		}
	}
}
