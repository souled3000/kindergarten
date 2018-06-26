package healthy

import (
	"github.com/astaxie/beego/orm"
	"time"
)

type Body struct {
	Id             int			`json:"id" orm:"column(id);auto" description:"编号"`
	Theme          string 		`json:"theme" orm:"column(theme)" description:"测评主题"`
	Total          int			`json:"total" orm:"column(total)" description:"总人数"`
	Actual         int			`json:"actual" orm:"column(actual)" description:"实际参数人数"`
	Rate           int			`json:"rate" orm:"column(rate)" description:"合格率"`
	TestTime       time.Time 	`json:"test_time" orm:"column(test_time)" description:"测评时间"`
	Mechanism      int			`json:"mechanism" orm:"column(mechanism)" description:"体检机构"`
	KindergartenId int			`json:"kindergarten_id" orm:"column(kindergarten_id)" description:"幼儿园ID"`
	Types          int			`json:"types" orm:"column(types)"`
	Project        string 		`json:"project" orm:"column(project)" description:"体检项目"`
	CreatedAt      time.Time 	`json:"created_at" orm:"auto_now_add" description:"创建时间"`
}

func (t *Body) TableName() string {
	return "healthy_body"
}

func init() {
	orm.RegisterModel(new(Body))
}
