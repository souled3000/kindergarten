package healthy

import (
	"github.com/astaxie/beego/orm"
	"time"
)

type Class struct {
	Id          int			`json:"id" orm:"column(id);auto" description:"编号"`
	BadyId      int			`json:"bady_id" orm:"column(bady_id)" description:"体质测评ID"`
	ClassId     int			`json:"class_id" orm:"column(class_id)" description:"参与班级"`
	ClassTotal  int			`json:"class_total" orm:"column(class_total)" description:"班级总人数"`
	ClassActual int			`json:"class_actual" orm:"column(class_actual)" description:"班级实际参赛人数"`
	ClassRate   int			`json:"class_rate" orm:"column(class_rate)" description:"合格率"`
	CreatedAt   time.Time	`json:"created_at" orm:"column(created_at)" description:"创建时间"`
}

func (t *Class) TableName() string {
	return "healthy_class"
}

func init() {
	orm.RegisterModel(new(Class))
}
