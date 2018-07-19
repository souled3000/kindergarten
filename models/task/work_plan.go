package task

import (
	"time"
	"github.com/astaxie/beego/orm"
)

type WorkPlan struct {
	Id             int       `json:"id"`
	Content        string    `json:"content"`
	PlanTime       time.Time `json:"plan_time"`
	Creator        int       `json:"creator"`
	CreatedAt      time.Time `json:"created_at" orm:"auto_now_add"`
	UpdatedAt      time.Time `json:"updated_at" orm:"auto_now"`
}

func (wp *WorkPlan) TableName() string {
	return "work_plan"
}

func init() {
	orm.RegisterModel(new(WorkPlan))
}

func (wp *WorkPlan) Save() (int64, error) {
	return orm.NewOrm().Insert(wp)
}

func (wp *WorkPlan) Get() ([]WorkPlan, error) {
	var wps []WorkPlan

	_, err := orm.NewOrm().QueryTable(wp).Filter("creator", wp.Creator).All(&wps)

	return wps, err
}