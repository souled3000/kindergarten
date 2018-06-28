package app

import (
	"github.com/astaxie/beego"
)

// SituationController operations for Situation
type SituationController struct {
	beego.Controller
}

// URLMapping ...
func (c *SituationController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
}
