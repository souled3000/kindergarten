package healthy

import (
	"github.com/astaxie/beego/orm"
	"time"
)

type Drug struct {
	Id        	int			`json:"id" orm:"column(id);auto" description:"编号"`
	StudentId 	int			`json:"student_id" orm:"column(student_id)" description:"学生ID"`
	Drug      	string 		`json:"drug" orm:"column(drug)" description:"药片"`
	Symptom	  	string		`json:"symptom" orm:"column(symptom)" description:"症状"`
	Explain   	string 		`json:"explain" orm:"column(explain)" description:"用量说明"`
	Url       	string		`json:"url" orm:"column(url)" description:"喂药申请图片"`
	UserId    	int			`json:"user_id" orm:"column(user_id)"`
	CreatedAt   time.Time	`json:"created_at" orm:"column(created_at);auto_now_add"`
}

func (t *Drug) TableName() string {
	return "healthy_drug"
}

func init() {
	orm.RegisterModel(new(Drug))
}

//申请喂药
func (m Drug) Save() error {
	o := orm.NewOrm()
	o.Insert(&m);

	return nil
}