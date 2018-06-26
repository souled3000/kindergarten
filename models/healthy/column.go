package healthy

import (
	"time"
	"github.com/astaxie/beego/orm"
)

type Column struct {
	Id		  int			`json:"id" orm:"column(id);auto"`
	InspectId int			`json:"inspect_id" orm:"column(inspect_id);" description:"编号"`
	StudentId int			`json:"student_id" orm:"column(student_id)"`
	Column1   string 		`json:"column_1" orm:"column(column_1)"`
	Column2   string 		`json:"column_2" orm:"column(column_2)"`
	Column3   string 		`json:"column_3" orm:"column(column_3)"`
	Column4   string 		`json:"column_4" orm:"column(column_4)"`
	Column5   string 		`json:"column_5" orm:"column(column_5)"`
	Column6   string 		`json:"column_6" orm:"column(column_6)"`
	Column7   string 		`json:"column_7" orm:"column(column_7)"`
	Column8   string 		`json:"column_8" orm:"column(column_8)"`
	Column9   string 		`json:"column_9" orm:"column(column_9)"`
	Column10  string 		`json:"column_10" orm:"column(column_10)"`
	CreatedAt time.Time		`json:"created_at" orm:"column(created_at)"`
}

func (t *Column) TableName() string {
	return "healthy_column"
}

func init() {
	orm.RegisterModel(new(Column))
}
