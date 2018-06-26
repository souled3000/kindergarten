package healthy

import (
	"github.com/astaxie/beego"
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
