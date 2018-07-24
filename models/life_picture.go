package models

import (
	"github.com/astaxie/beego/orm"
)

type LifePicture struct {
	Id      int    `json:"id" orm:"column(id);auto"`
	LifeId  int    `json:"life_id" orm:"column(life_id);" description:"生活id"`
	Picture string `json:"picture" orm:"column(picture)"`
	Number  int    `json:"number" orm:"column(number)" description:"编号"`
}

func (t *LifePicture) TableName() string {
	return "life_picture"
}

func init() {
	orm.RegisterModel(new(LifePicture))
}
