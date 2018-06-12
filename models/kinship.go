package models

import (
	"time"

	"github.com/astaxie/beego/orm"
)

type Kinship struct {
	Id                 int       `json:"kinship_id" orm:"column(kinship_id);auto" description:"编号"`
	StudentId          int       `json:"student_id" orm:"column(student_id)" description:"学生序号"`
	UserId             int       `json:"user_id" orm:"column(user_id)" description:"用户ID"`
	Type               int8      `json:"type" orm:"column(type)" description:"类型，1紧急联系人，2监护人"`
	Relation           string    `json:"relation" orm:"column(relation);size(10)" description:"关系"`
	Name               string    `json:"name" orm:"column(name);size(20)" description:"名字"`
	UnitName           string    `json:"unit_name" orm:"column(unit_name);size(50)" description:"单位名称"`
	ContactInformation string    `json:"contact_information" orm:"column(contact_information);size(11)" description:"联系方式"`
	CreatedAt          time.Time `json:"created_at" orm:"auto_now_add"`
	UpdatedAt          time.Time `json:"updated_at" orm:"auto_now"`
}

func (t *Kinship) TableName() string {
	return "kinship"
}

func init() {
	orm.RegisterModel(new(Kinship))
}
