package models

import (
	"time"

	"github.com/astaxie/beego/orm"
)

type BabyKindergarten struct {
	Id             int       `orm:"column(baby_kindergarten_id);auto"`
	BabyId         int       `orm:"column(baby_id)"`
	BabyName       string    `orm:"column(baby_name)"`
	KindergartenId int       `orm:"column(kindergarten_id)"`
	Birthday       time.Time `orm:"column(birthday)"`
	Actived        int       `orm:"column(actived)"`
	Status         int       `orm:"column(status)"`
	InviteStatus   int       `orm:"column(invite_status)"`
	CreatedAt      time.Time `orm:"auto_now"`
	UpdateAt       time.Time `orm:"auto_now_add"`
}

func (t *BabyKindergarten) TableName() string {
	return "baby_kindergarten"
}

func init() {
	orm.RegisterModel(new(BabyKindergarten))
}
