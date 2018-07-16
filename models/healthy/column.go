package healthy

import (
	"github.com/astaxie/beego/orm"
	"time"
)

type Column struct {
	Id         int       `json:"id" orm:"column(id);auto"`
	InspectId  int       `json:"inspect_id" orm:"column(inspect_id);" description:"编号"`
	StudentId  int       `json:"student_id" orm:"column(student_id)"`
	Column1    string    `json:"column1" orm:"column(column1)"`
	Abnormal1  string    `json:"abnormal1" orm:"column(abnormal1)"`
	Column2    string    `json:"column2" orm:"column(column2)"`
	Abnormal2  string    `json:"abnormal2" orm:"column(abnormal2)"`
	Column3    string    `json:"column3" orm:"column(column3)"`
	Abnormal3  string    `json:"abnormal3" orm:"column(abnormal3)"`
	Column4    string    `json:"column4" orm:"column(column4)"`
	Abnormal4  string    `json:"abnormal4" orm:"column(abnormal4)"`
	Column5    string    `json:"column5" orm:"column(column5)"`
	Abnormal5  string    `json:"abnormal5" orm:"column(abnormal5)"`
	Column6    string    `json:"column6" orm:"column(column6)"`
	Abnormal6  string    `json:"abnormal6" orm:"column(abnormal6)"`
	Column7    string    `json:"column7" orm:"column(column7)"`
	Abnormal7  string    `json:"abnormal7" orm:"column(abnormal7)"`
	Column8    string    `json:"column8" orm:"column(column8)"`
	Abnormal8  string    `json:"abnormal8" orm:"column(abnormal8)"`
	Column9    string    `json:"column9" orm:"column(column9)"`
	Abnormal9  string    `json:"abnormal9" orm:"column(abnormal9)"`
	Column10   string    `json:"column10" orm:"column(column10)"`
	Abnormal10 string    `json:"abnormal10" orm:"column(abnormal10)"`
	CreatedAt  time.Time `json:"created_at" orm:"column(created_at)"`
}

func (t *Column) TableName() string {
	return "healthy_column"
}

func init() {
	orm.RegisterModel(new(Column))
}
