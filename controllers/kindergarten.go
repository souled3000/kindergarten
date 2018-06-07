package controllers

import (
	"encoding/json"
	"errors"
	"kindergarten-service-go/models"
	"strconv"
	"strings"

	"github.com/astaxie/beego"
)

// KindergartenController operations for Kindergarten
type KindergartenController struct {
	beego.Controller
}

// URLMapping ...
func (c *KindergartenController) URLMapping() {
	c.Mapping("GetIntroduceInfo", c.GetIntroduceInfo)
	c.Mapping("GetIntroduce", c.GetIntroduce)
	c.Mapping("Put", c.Put)
}

// GetIntroduceInfo ...
// @Title 幼儿园介绍详情
// @Description 幼儿园介绍详情
// @Param	id		path 	string	true		"幼儿园ID"
// @Success 200 {object} models.Kindergarten
// @Failure 403 :id is empty
// @router /:id [get]
func (c *KindergartenController) GetIntroduceInfo() {
	prepage, _ := c.GetInt("per_page", 20)
	page, _ := c.GetInt("page")
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	v := models.GetKindergartenById(id, page, prepage)
	if v == nil {
		c.Data["json"] = JSONStruct{"error", 1005, v, "获取失败"}
	} else {
		c.Data["json"] = JSONStruct{"success", 0, v, "获取成功"}
	}
	c.ServeJSON()
}

// GetIntroduce ...
// @Title 幼儿园介绍列表
// @Description 幼儿园介绍列表
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.Kindergarten
// @Failure 403
// @router / [get]
func (c *KindergartenController) GetIntroduce() {
	var fields []string
	var sortby []string
	var order []string
	var query = make(map[string]string)
	var limit int64 = 10
	var offset int64

	// fields: col1,col2,entity.col3
	if v := c.GetString("fields"); v != "" {
		fields = strings.Split(v, ",")
	}
	// limit: 10 (default is 10)
	if v, err := c.GetInt64("limit"); err == nil {
		limit = v
	}
	// offset: 0 (default is 0)
	if v, err := c.GetInt64("offset"); err == nil {
		offset = v
	}
	// sortby: col1,col2
	if v := c.GetString("sortby"); v != "" {
		sortby = strings.Split(v, ",")
	}
	// order: desc,asc
	if v := c.GetString("order"); v != "" {
		order = strings.Split(v, ",")
	}
	// query: k:v,k:v
	if v := c.GetString("query"); v != "" {
		for _, cond := range strings.Split(v, ",") {
			kv := strings.SplitN(cond, ":", 2)
			if len(kv) != 2 {
				c.Data["json"] = errors.New("Error: invalid query key/value pair")
				c.ServeJSON()
				return
			}
			k, v := kv[0], kv[1]
			query[k] = v
		}
	}

	l, err := models.GetAllKindergarten(query, fields, sortby, order, offset, limit)
	if err != nil {
		c.Data["json"] = err.Error()
	} else {
		c.Data["json"] = l
	}
	c.ServeJSON()
}

// Put ...
// @Title Put
// @Description update the Kindergarten
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.Kindergarten	true		"body for Kindergarten content"
// @Success 200 {object} models.Kindergarten
// @Failure 403 :id is not int
// @router /:id [put]
func (c *KindergartenController) Put() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	v := models.Kindergarten{Id: id}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		if err := models.UpdateKindergartenById(&v); err == nil {
			c.Data["json"] = "OK"
		} else {
			c.Data["json"] = err.Error()
		}
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}
