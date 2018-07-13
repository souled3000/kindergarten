package healthy

import (
	"github.com/astaxie/beego/orm"
	"time"
	"strconv"
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
func AddClass(b *Class,body_id int, class_id int,types int) (err error){
	o := orm.NewOrm()
	var num Num
	sql := "select count(id) as num from organizational_member where type = 1 and organizational_id = "+strconv.Itoa(b.ClassId)
	o.Raw(sql).QueryRow(&num)
	b.ClassTotal = num.Num
	o.Begin()

	var some_err []interface{}
	if _, _, err = o.ReadOrCreate(b, "BodyId","ClassId"); err != nil {
		some_err = append(some_err,err)
	}
	cnt, _ := o.QueryTable("healthy_inspect").Filter("body_id",b.BodyId).Filter("class_id",b.ClassId).Count()
	//判断班级是否添加
	if cnt == 0{
		var inspect []Inspect
		sql = "select b.student_id,a.organizational_id as class_id,b.kindergarten_id,c.name as class_name from organizational_member a left join student b on b.student_id=a.member_id left join organizational c on c.id=a.organizational_id where a.type=1 and a.organizational_id="+strconv.Itoa(class_id)
		if _,err = o.Raw(sql).QueryRows(&inspect); err != nil{
			some_err = append(some_err,err)
		}

		for key,_ := range inspect {
			inspect[key].BodyId = body_id
			inspect[key].Types = types
		}
		if _,err = o.InsertMulti(len(inspect),inspect); err != nil {
			some_err = append(some_err,err)
		}
	}

	if len(some_err) > 0 {
		o.Rollback()
	} else {
		o.Commit()
		return  nil
	}
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
	sql := "select a.*,c.name as class_name,c.class_type,b.theme from healthy_class a left join healthy_body b on a.body_id = b.id left join organizational c on c.id = a.class_id "+where+" order by a.id desc "
	sqlNum := "select count(a.id) as num from healthy_class a left join healthy_body b on a.body_id = b.id left join organizational c on c.id = a.class_id "+where+" order by a.id desc "
	limit := " limit "+strconv.Itoa((page-1)*per_page)+","+strconv.Itoa(per_page)
	ml = make(map[string]interface{})
	if _,err = o.Raw(sql+limit).Values(&d); err == nil {
		type Num struct {
			Num		int 	`json:"num"`
		}
		var total Num
		err = o.Raw(sqlNum).QueryRow(&total);
		pageNum := int(math.Ceil(float64(total.Num) / float64(per_page)))
		for _, val := range d {
			if val["class_type"].(string) == "3"{
				val["class_name"] = "大班"+val["class_name"].(string)
			} else if val["class_type"].(string) == "2"{
				val["class_name"] = "中班"+val["class_name"].(string)
			}else if val["class_type"].(string) == "1"{
				val["class_name"] = "小班"+val["class_name"].(string)
			}
		}
		ml["data"] = d
		ml["total"] = total.Num
		ml["pageNum"] = pageNum
		ml["limit"] = per_page
		return ml,nil
	}
	return nil,err
}

//删除
func DeleteClass(id int) map[string]interface{} {
	o := orm.NewOrm()
	v := Class{Id: id}
	if err := o.Read(&v); err == nil {
		if num, err := o.Delete(&Class{Id: id}); err == nil {
			paginatorMap := make(map[string]interface{})
			paginatorMap["data"] = num
			return paginatorMap
		}
	}
	return nil
}