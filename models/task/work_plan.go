package task

import (
	"time"
	"github.com/astaxie/beego/orm"
	"encoding/json"
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

func (wp *WorkPlan) Get() ([]map[string]interface{}, error) {
	var wps []WorkPlan
	var res []map[string]interface{}

	if num, err := orm.NewOrm().QueryTable(wp).Filter("creator", wp.Creator).All(&wps); err == nil && num > 0 {
		for _, val := range wps {
			jsons, _ := json.Marshal(val)
			var maps map[string]interface{}
			json.Unmarshal(jsons, &maps)
			maps["plan_time"] = val.PlanTime.Format("15:04")

			res = append(res, maps)
		}

		return res, err
	} else {
		return res, err
	}
}

func (wp *WorkPlan) Delete() error {
	o := orm.NewOrm()

	if err := o.Read(wp); err != nil {
		return err
	}

	if _, err := o.Delete(wp); err != nil {
		return err
	}

	return nil
}