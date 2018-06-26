package healthy

import (
	"github.com/astaxie/beego/orm"
	"time"
)

type Situation struct {
	Id       	int	   	  `json:"id" orm:"column(id);auto"`
	Name     	string    `json:"name" orm:"column(name);"`
	CreatedAt   time.Time `json:"created_at" orm:"auto_now_add"`
}

func init() {
	orm.RegisterModel(new(Situation))
}
																