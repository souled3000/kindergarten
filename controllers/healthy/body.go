package healthy

import (
	"github.com/astaxie/beego"
	"kindergarten-service-go/models/healthy"
	"strconv"
	"fmt"
)

// BodyController operations for Body
type BodyController struct {
	beego.Controller
}

// URLMapping ...
func (c *BodyController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
}

// Post ...
// @Title
// @Description 添加主题
// @Param	theme		query	string	false	"主题名"
// @Param	total		query	int	false	"总人数"
// @Param	actual		query	int	false	"实际人数"
// @Param	rate		query	string	false	"合格率"
// @Param	test_time	query	string	false	"测评日期"
// @Param	mechanism	query	int	false	"体检机构id"
// @Param	kindergarten_id	query	int	false	"幼儿园id"
// @Param	types	query	int	false	"类型 1，体质测评2，体检"
// @Param	project	query	string	false	"体检项目"
// @Success 201 {int} models.Category
// @Failure 403 body is empty
// @router / [post]
func (c *BodyController) Post() {
	theme := c.GetString("theme")
	total,_ := c.GetInt("total")
	actual,_ := c.GetInt("actual")
	rate,_ := c.GetInt("rate")
	test_time := c.GetString("test_time")
	mechanism,_ := c.GetInt("mechanism")
	kindergarten_id,_ := c.GetInt("kindergarten_id")
	types,_ := c.GetInt("types")
	project := c.GetString("project")
	fmt.Println(project)
	if project == ""{
		project = "column1:左眼,column2:右眼,column3:血小板,column4:龋齿,column5:体脂率"
		fmt.Println(project)
	}
	fmt.Println(project)
	var b healthy.Body
	b.Theme = theme
	b.Total = total
	b.Actual = actual
	b.Rate = rate
	b.TestTime = test_time
	b.Mechanism = mechanism
	b.KindergartenId = kindergarten_id
	b.Types = types
	b.Project = project
	if _,err := healthy.AddBody(&b); err == nil {
		c.Data["json"] = JSONStruct{"success", 0, nil, "添加成功"}
	} else {
		c.Data["json"] = JSONStruct{"error", 1001, err.Error(), "添加失败"}
	}
	c.ServeJSON()
}

// Put ...
// @Title 修改体检主题
// @Description 修改体检主题
// @Param	theme		query	string	false	"主题名"
// @Param	total		query	int	false	"总人数"
// @Param	actual		query	int	false	"实际人数"
// @Param	rate		query	string	false	"合格率"
// @Param	test_time	query	string	false	"测评日期"
// @Param	mechanism	query	int	false	"体检机构id"
// @Param	kindergarten_id	query	int	false	"幼儿园id"
// @Param	types	query	int	false	"类型 1，体质测评2，体检"
// @Param	project	query	string	false	"体检项目"
// @Success 200 {object} models.Cover
// @Failure 403 :id is not int
// @router /:id [put]
func (c *BodyController) Put() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	theme := c.GetString("theme")
	total,_ := c.GetInt("total")
	actual,_ := c.GetInt("actual")
	rate,_ := c.GetInt("rate")
	test_time := c.GetString("test_time")
	mechanism,_ := c.GetInt("mechanism")
	kindergarten_id,_ := c.GetInt("kindergarten_id")
	types,_ := c.GetInt("types")
	project := c.GetString("project")
	var b healthy.Body
	b.Id = id
	b.Theme = theme
	b.Total = total
	b.Actual = actual
	b.Rate = rate
	b.TestTime = test_time
	b.Mechanism = mechanism
	b.KindergartenId = kindergarten_id
	b.Types = types
	b.Project = project
	if err := healthy.UpdataByIdBody(&b); err == nil {
		c.Data["json"] = JSONStruct{"success", 0, nil, "修改成功"}
	} else {
		c.Data["json"] = JSONStruct{"error", 1001, err.Error(), "修改失败"}
	}
	c.ServeJSON()
}

// GetAll ...
// @Title 体检主题列表
// @Description 体检主题列表
// @Param	page	query	int	false	"页"
// @Param	per_page	query	int	false	"每页条数"
// @Param	types	query	int	false	"类型"
// @Param	theme	query	string	false	"名字搜索"
// @Success 200 {object} healthy.Body
// @Failure 403
// @router / [get]
func (c *BodyController) GetAll() {
	page := 1
	per_page := 20
	if v, err := c.GetInt("page"); err == nil {
		page = v
	}
	if v, err := c.GetInt("per_page"); err == nil {
		per_page = v
	}
	types,_ := c.GetInt("type")
	theme := c.GetString("theme")
	if l,err := healthy.GetAllBody(page,per_page,types,theme); err == nil {
		c.Data["json"] = JSONStruct{"success", 0, l, "获取成功"}
	} else {
		c.Data["json"] = JSONStruct{"error", 1001, err.Error(), "获取失败"}
	}
	c.ServeJSON()
}

// GetOne ...
// @Title 主题详情
// @Description 主题详情
// @Param	id		path 	string	true		"自增ID"
// @Success 200 {object} healthy.Body
// @Failure 403
// @router /:id  [get]
func (c *BodyController) GetOne() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	class_id,_ := c.GetInt("class_id")
	if class_id > 0 {
		if l,err := healthy.GetOneBodyClass(id,class_id); err == nil {
			c.Data["json"] = JSONStruct{"success", 0, l, "获取成功"}
		} else {
			c.Data["json"] = JSONStruct{"error", 1001, err.Error(), "获取失败"}
		}
	}else{
		if l,err := healthy.GetOneBody(id); err == nil {
			c.Data["json"] = JSONStruct{"success", 0, l, "获取成功"}
		} else {
			c.Data["json"] = JSONStruct{"error", 1001, err.Error(), "获取失败"}
		}
	}

	c.ServeJSON()
}