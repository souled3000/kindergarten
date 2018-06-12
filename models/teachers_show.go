package models

import (
	"time"

	"github.com/astaxie/beego/orm"
)

type TeachersShow struct {
	Id             int       `orm:"column(id);auto"`
	TeacherId      int       `orm:"column(teacher_id)" description:"教师ID"`
	Introduction   string    `orm:"column(introduction);size(100)" description:"介绍"`
	KindergartenId int       `orm:"column(kindergarten_id)" description:"幼儿园ID"`
	CreatedAt      time.Time `orm:"column(created_at);type(datetime)"`
	UpdatedAt      time.Time `orm:"column(updated_at);type(datetime)"`
}

func (t *TeachersShow) TableName() string {
	return "teachers_show"
}

func init() {
	orm.RegisterModel(new(TeachersShow))
}
