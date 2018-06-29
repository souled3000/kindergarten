package healthy

import (
	"github.com/astaxie/beego/orm"
	"time"
)

type Body struct {
	Id             int			`json:"id" orm:"column(id);auto" description:"编号"`
	Theme          string 		`json:"theme" orm:"column(theme)" description:"测评主题"`
	Total          int			`json:"total" orm:"column(total)" description:"总人数"`
	Actual         int			`json:"actual" orm:"column(actual)" description:"实际参数人数"`
	Rate           int			`json:"rate" orm:"column(rate)" description:"合格率"`
	TestTime       string 		`json:"test_time" orm:"column(test_time)" description:"测评时间"`
	Mechanism      int			`json:"mechanism" orm:"column(mechanism)" description:"体检机构"`
	KindergartenId int			`json:"kindergarten_id" orm:"column(kindergarten_id)" description:"幼儿园ID"`
	Types          int			`json:"types" orm:"column(types)"`
	Project        string 		`json:"project" orm:"column(project)" description:"体检项目"`
	CreatedAt      time.Time 	`json:"created_at" orm:"auto_now_add" description:"创建时间"`
}

func (t *Body) TableName() string {
	return "healthy_body"
}

func init() {
	orm.RegisterModel(new(Body))
}
func AddBody(b *Body) (id int64, err error){
	o := orm.NewOrm()
	id, err = o.Insert(b)
	return
}

func UpdataByIdBody(b *Body) (err error) {
	o := orm.NewOrm()
	v := Body{Id:b.Id}
	if err := o.Read(&v); err == nil {
		if b.Project != "" {
			v.Project = b.Project
		}

		if b.Types > 0 {
			v.Types = b.Types
		}
		if b.KindergartenId > 0 {
			v.KindergartenId = b.KindergartenId
		}
		if b.Mechanism > 0 {
			v.Mechanism = b.Mechanism
		}
		if b.TestTime != "" {
			v.TestTime = b.TestTime
		}
		if b.Rate > 0 {
			v.Rate = b.Rate
		}
		if b.Actual > 0 {
			v.Actual = b.Actual
		}
		if b.Total > 0 {
			v.Total = b.Total
		}
		_,err = o.Update(&v)
	}

	return err
}

func GetAllBody(page int,per_page int,types int,theme string) (ml map[string]interface{}, err error){
	o := orm.NewOrm()
	qs := o.QueryTable(new(Body))
	if types > 0 {
		qs = qs.Filter("types", types)
	}
	if theme != "" {
		qs = qs.Filter("theme", theme)
	}
	var d []Body

	ml = make(map[string]interface{})
	if _,err = qs.Limit(per_page,(page-1)*per_page).OrderBy("-id").All(&d); err == nil {
		num,_ := qs.Count()
		ml["data"] = d
		ml["total"] = num
		return ml,nil
	}
	return nil,err
}
//添加或查询
func CrBody(theme string, kindergarten_id int,test_time string, types int)(id int64, err error){
	o := orm.NewOrm()
	body := Body{Theme: theme, KindergartenId:kindergarten_id,TestTime:test_time,Types:types}
	// 三个返回参数依次为：是否新创建的，对象 Id 值，错误
	if _, id, err := o.ReadOrCreate(&body, "Theme","KindergartenId","TestTime"); err == nil {
			return  id,nil
	}
	return  id,err
}