package app

import (
	"github.com/astaxie/beego"
)

// 体检项目
type ColumnController struct {
	beego.Controller
}

// URLMapping ...
func (c *ColumnController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
}
