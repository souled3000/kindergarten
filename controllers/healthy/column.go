package healthy

import (
	"github.com/astaxie/beego"
)

// ColumnController operations for Column
type ColumnController struct {
	beego.Controller
}

// URLMapping ...
func (c *ColumnController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
}
