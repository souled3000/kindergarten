package models

import (
	"time"

	"github.com/astaxie/beego/orm"
)

type GroupView struct {
	Id        int       `orm:"column(id);auto"`
	UserId    int       `orm:"column(user_id)"`
	ClassType int8      `orm:"column(class_type)"`
	CreatedAt time.Time `orm:"auto_now"`
}

func (t *GroupView) TableName() string {
	return "group_view"
}

func init() {
	orm.RegisterModel(new(GroupView))
}
