package healthy

import (
	"github.com/astaxie/beego/orm"
	"time"
)

type Situation struct {
	Id        int       `json:"id" orm:"column(id);auto"`
	Name      string    `json:"name" orm:"column(name);"`
	Type      int       `json:"type" orm:"column(type)"`
	CreatedAt time.Time `json:"created_at" orm:"auto_now_add"`
}

func init() {
	orm.RegisterModel(new(Situation))
}

func (t *Situation) TableName() string {
	return "healthy_situation"
}

//创建异常记录
func (m Situation) Post() error {
	o := orm.NewOrm()
	o.Insert(&m)

	return nil
}

//删除
func DeleteSituation(id int) map[string]interface{} {
	o := orm.NewOrm()
	v := Situation{Id: id}
	if err := o.Read(&v); err == nil {
		if num, err := o.Delete(&Situation{Id: id}); err == nil {
			paginatorMap := make(map[string]interface{})
			paginatorMap["data"] = num
			return paginatorMap
		}
	}
	return nil
}

//病例列表
func (f *Situation) GetAll(types int) ([]orm.Params, error) {
	o := orm.NewOrm()
	var con []interface{}
	where := "1 "
	where += "AND type = ? "
	con = append(con, types)

	var sxWords []orm.Params

	qb, _ := orm.NewQueryBuilder("mysql")
	sql := qb.Select("*").From(f.TableName()).Where(where).String()

	if _, err := o.Raw(sql, con).Values(&sxWords); err == nil {

		return sxWords, nil
	}

	return nil, nil
}
