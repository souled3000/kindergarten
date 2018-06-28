package healthy

import (
	"github.com/astaxie/beego/orm"
	"time"
	"strconv"
	"fmt"
	"reflect"
	"math"
)

type Class struct {
	Id          int			`json:"id" orm:"column(id);auto" description:"编号"`
	BodyId      int			`json:"body_id" orm:"column(body_id)" description:"体质测评ID"`
	ClassId     int			`json:"class_id" orm:"column(class_id)" description:"参与班级"`
	ClassTotal  int			`json:"class_total" orm:"column(class_total)" description:"班级总人数"`
	ClassActual int			`json:"class_actual" orm:"column(class_actual)" description:"班级实际参赛人数"`
	ClassRate   int			`json:"class_rate" orm:"column(class_rate)" description:"合格率"`
	CreatedAt   time.Time	`json:"created_at" orm:"column(created_at)" description:"创建时间"`
}

func (t *Class) TableName() string {
	return "healthy_class"
}

func init() {
	orm.RegisterModel(new(Class))
}
func AddClass(b *Class) (id int64, err error){
	o := orm.NewOrm()
	id, err = o.Insert(b)
	return
}

func UpdataByIdClass(b *Class) (err error) {
	o := orm.NewOrm()
	v := Class{Id:b.Id}
	if err := o.Read(&v); err == nil {

		if b.ClassId > 0 {
			v.ClassId = b.ClassId
		}
		if b.BodyId > 0 {
			v.BodyId = b.BodyId
		}
		if b.ClassTotal > 0 {
			v.ClassTotal = b.ClassTotal
		}
		if b.ClassActual > 0 {
			v.ClassActual = b.ClassActual
		}
		if b.ClassRate > 0 {
			v.ClassRate = b.ClassRate
		}
		_,err = o.Update(&v)
	}

	return err
}

func GetAllClass(page int,per_page int,class_id int,body_id int) (ml map[string]interface{}, err error){
	o := orm.NewOrm()
	where := " where 1=1"
	if class_id > 0 {
		where += " and a.class_id = "+strconv.Itoa(class_id)
	}
	if body_id > 0 {
		where += " and a.body_id = "+strconv.Itoa(body_id)
	}
	var d []orm.Params
	sql := "select a.*,c.name as class_name,b.theme from healthy_class a left join healthy_body b on a.body_id = b.id left join organizational c on c.id = a.class_id "+where+" order by a.id desc "
	sqlNum := "select count(a.id) as num from healthy_class a left join healthy_body b on a.body_id = b.id left join organizational c on c.id = a.class_id "+where+" order by a.id desc "
	limit := " limit "+strconv.Itoa((page-1)*per_page)+","+strconv.Itoa(per_page)
	ml = make(map[string]interface{})
	if _,err = o.Raw(sql+limit).Values(&d); err == nil {
		type Num struct {
			Num		int 	`json:"num"`
		}
		var total Num
		err = o.Raw(sqlNum).QueryRow(&total);
		fmt.Print(total,reflect.TypeOf(total))
		pageNum := int(math.Ceil(float64(total.Num) / float64(per_page)))
		ml["data"] = d
		ml["total"] = total.Num
		ml["pageNum"] = pageNum
		ml["limit"] = per_page
		return ml,nil
	}
	return nil,err
}
