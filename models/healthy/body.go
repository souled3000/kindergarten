package healthy

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/orm"
	"math"
	"strconv"
	"strings"
	"time"
)

type Body struct {
	Id             int       `json:"id" orm:"column(id);auto" description:"编号"`
	Theme          string    `json:"theme" orm:"column(theme)" description:"测评主题"`
	Total          int       `json:"total" orm:"column(total)" description:"总人数"`
	Actual         int       `json:"actual" orm:"column(actual)" description:"实际参数人数"`
	Rate           int       `json:"rate" orm:"column(rate)" description:"合格率"`
	TestTime       string    `json:"test_time" orm:"column(test_time)" description:"测评时间"`
	Mechanism      int       `json:"mechanism" orm:"column(mechanism)" description:"体检机构"`
	KindergartenId int       `json:"kindergarten_id" orm:"column(kindergarten_id)" description:"幼儿园ID"`
	Types          int       `json:"types" orm:"column(types)"`
	Project        string    `json:"project" orm:"column(project)" description:"体检项目"`
	CreatedAt      time.Time `json:"created_at" orm:"auto_now_add" description:"创建时间"`
}

type Num struct {
	Num int `json:"num"`
}

func (t *Body) TableName() string {
	return "healthy_body"
}

func init() {
	orm.RegisterModel(new(Body))
}
func AddBody(b *Body) (id int64, err error) {
	o := orm.NewOrm()
	var num Num

	sql := "select count(student_id) as num from student where kindergarten_id = " + strconv.Itoa(b.KindergartenId)
	o.Raw(sql).QueryRow(&num)
	b.Total = num.Num
	fmt.Println(num.Num)
	id, err = o.Insert(b)
	return
}

func GetOneBody(id int) (ml map[string]interface{}, err error) {
	o := orm.NewOrm()
	var b Body
	b.Id = id
	var num Num
	var list []map[string]interface{}
	sql := "select count(a.id) as num from healthy_inspect a where a.body_id = " + strconv.Itoa(id)
	o.Raw(sql).QueryRow(&num)
	c_num := num.Num
	list2 := make(map[string]interface{})
	sql = "select count(a.id) as num from healthy_inspect a where a.body_id = " + strconv.Itoa(id) + " and weight is != '' "
	o.Raw(sql).QueryRow(&num)
	fmt.Println(num.Num)
	bili := int(math.Ceil(float64(num.Num) / float64(c_num) * 100.0))
	if bili < 0 {
		list2["bili"] = 0
	} else {
		list2["bili"] = bili
	}
	list2["column"] = "weight"
	list2["columnh"] = "height"
	list2["name"] = "体重体重"
	list = append(list, list2)
	sql = "select count(a.id) as num from healthy_inspect a where a.body_id = " + strconv.Itoa(id) + " and height != '' "
	o.Raw(sql).QueryRow(&num)
	if err := o.Read(&b); err == nil {
		var c string
		project := strings.Split(b.Project, ",")
		for _, val := range project {

			list1 := make(map[string]interface{})
			cloumn := strings.Split(val, ":")
			sql := "select count(b.id) as num from healthy_inspect a left join healthy_column b on a.id= b.inspect_id where a.body_id = " + strconv.Itoa(id) + " and b." + string(cloumn[0]) + " != '' "
			o.Raw(sql).QueryRow(&num)
			fmt.Println(num.Num)
			bili := int(math.Ceil(float64(num.Num) / float64(c_num) * 100.0))
			if bili < 0 {
				list1["bili"] = 0
			} else {
				list1["bili"] = bili
			}
			list1["column"] = cloumn[0]

			if strings.Contains(val, "眼") {
				if c == "" {
					c = cloumn[0]
				} else {
					list1["columnx"] = c
					list1["name"] = "视力"
					list = append(list, list1)
				}

			} else {
				list1["name"] = string(cloumn[1])
				list = append(list, list1)
			}

		}
		bjson, _ := json.Marshal(b)
		json.Unmarshal(bjson, &ml)
		ml["info"] = list
		return ml, nil
	}
	return nil, err
}

func GetOneBodyClass(id int, class_id int) (ml map[string]interface{}, err error) {
	o := orm.NewOrm()
	var num Num
	var list []map[string]interface{}
	sql := "select count(a.id) as num from healthy_inspect a where a.class_id=" + strconv.Itoa(class_id) + " and a.body_id = " + strconv.Itoa(id)
	o.Raw(sql).QueryRow(&num)
	c_num := num.Num
	list2 := make(map[string]interface{})
	//list3 := make(map[string]interface{})
	sql = "select count(a.id) as num from healthy_inspect a where a.class_id=" + strconv.Itoa(class_id) + " and a.body_id = " + strconv.Itoa(id) + " and weight != '' "
	o.Raw(sql).QueryRow(&num)
	fmt.Println(num.Num)
	bili := int(math.Ceil(float64(num.Num) / float64(c_num) * 100.0))
	if bili < 0 {
		list2["bili"] = 0
	} else {
		list2["bili"] = bili
	}
	list2["column"] = "weight"
	list2["columnh"] = "height"
	list2["name"] = "身高体重"
	list = append(list, list2)
	sql = "select count(a.id) as num from healthy_inspect a where a.class_id=" + strconv.Itoa(class_id) + " and a.body_id = " + strconv.Itoa(id) + " and height != '' "
	o.Raw(sql).QueryRow(&num)
	fmt.Println(num.Num)
	//list3["bili"] = int(math.Ceil(float64(num.Num)/float64(c_num)*100.0))
	//list3["name"] = "身高"
	//list = append(list,list3)
	sql = "select a.id,a.theme,a.test_time,a.mechanism,a.kindergarten_id,a.types,a.project,b.class_total as total,b.class_actual as actual,b.class_rate as rate from healthy_body a left join healthy_class b on a.id = b.body_id where a.id=" + strconv.Itoa(id) + " and b.class_id=" + strconv.Itoa(class_id)
	var b Body

	if err = o.Raw(sql).QueryRow(&b); err == nil {
		var c string
		project := strings.Split(b.Project, ",")
		for _, val := range project {
			list1 := make(map[string]interface{})
			cloumn := strings.Split(val, ":")
			sql := "select count(b.id) as num from healthy_inspect a left join healthy_column b on a.id= b.inspect_id where  a.class_id=" + strconv.Itoa(class_id) + " and a.body_id = " + strconv.Itoa(id) + " and b." + string(cloumn[0]) + " != '' "
			o.Raw(sql).QueryRow(&num)
			fmt.Println(num.Num)
			bili := int(math.Ceil(float64(num.Num) / float64(c_num) * 100.0))
			if bili < 0 {
				list1["bili"] = 0
			} else {
				list1["bili"] = bili
			}
			list1["column"] = cloumn[0]
			if strings.Contains(val, "眼") {
				if c == "" {
					c = cloumn[0]
				} else {
					list1["columnx"] = c
					list1["name"] = "视力"
					list = append(list, list1)
				}

			} else {
				list1["name"] = string(cloumn[1])
				list = append(list, list1)
			}
		}
		bjson, _ := json.Marshal(b)
		json.Unmarshal(bjson, &ml)
		ml["info"] = list
		return ml, nil
	}
	return nil, err
}

func UpdataByIdBody(b *Body) (err error) {
	o := orm.NewOrm()
	v := Body{Id: b.Id}
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
		if b.Theme != "" {
			v.Theme = b.Theme
		}
		_, err = o.Update(&v)
	}

	return err
}

func GetAllBody(kindergarten_id, page int, per_page int, types int, theme string, search string) (ml map[string]interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Body))
	if types > 0 {
		qs = qs.Filter("types", types)
	}
	if theme != "" {
		qs = qs.Filter("theme", theme)
	}
	if kindergarten_id != 0 {
		qs = qs.Filter("kindergarten_id", kindergarten_id)
	}
	if search != "" {
		qs = qs.Filter("theme__contains", search)
	}

	var d []Body

	ml = make(map[string]interface{})
	if _, err = qs.Limit(per_page, (page-1)*per_page).OrderBy("-test_time").All(&d); err == nil {
		num, _ := qs.Count()
		var dd []map[string]interface{}
		djson, _ := json.Marshal(d)
		json.Unmarshal(djson, &dd)
		for i, _ := range dd {
			dd[i]["name"] = i + 1
			fmt.Println(i)

		}

		ml["data"] = dd
		ml["total"] = num
		return ml, nil
	}
	return nil, err
}

//添加或查询
func CrBody(theme string, kindergarten_id int, test_time string, types int) (id int64, err error) {
	o := orm.NewOrm()
	body := Body{Theme: theme, KindergartenId: kindergarten_id, TestTime: test_time, Types: types}
	var num Num
	sql := "select count(student_id) as num from student where kindergarten_id = " + strconv.Itoa(kindergarten_id)
	o.Raw(sql).QueryRow(&num)
	body.Total = num.Num
	// 三个返回参数依次为：是否新创建的，对象 Id 值，错误
	body.Project = "column1:左眼,column2:右眼,column3:血小板,column4:龋齿"
	if _, id, err := o.ReadOrCreate(&body, "Theme", "KindergartenId", "TestTime"); err == nil {
		return id, nil
	}
	return id, err
}

//删除
func Delete(id int) map[string]interface{} {
	o := orm.NewOrm()
	v := Body{Id: id}
	if err := o.Read(&v); err == nil {
		if num, err := o.Delete(&Body{Id: id}); err == nil {
			paginatorMap := make(map[string]interface{})
			paginatorMap["data"] = num
			return paginatorMap
		}
	}
	return nil
}

func GetOneBodyClasss(id int, class_id int) (ml map[string]interface{}, err error) {
	o := orm.NewOrm()
	var num Num
	var list []map[string]interface{}
	sql := "select count(a.id) as num from healthy_inspect a where a.class_id=" + strconv.Itoa(class_id) + " and a.body_id = " + strconv.Itoa(id)
	o.Raw(sql).QueryRow(&num)
	c_num := num.Num
	list2 := make(map[string]interface{})
	list3 := make(map[string]interface{})
	list4 := make(map[string]interface{})
	list5 := make(map[string]interface{})
	sql = "select count(a.id) as num from healthy_inspect a where a.class_id=" + strconv.Itoa(class_id) + " and a.body_id = " + strconv.Itoa(id) + " and weight is not null"
	o.Raw(sql).QueryRow(&num)
	fmt.Println(num.Num)
	bili := int(math.Ceil(float64(num.Num) / float64(c_num) * 100.0))
	if bili < 0 {
		list2["bili"] = 0
	} else {
		list2["bili"] = bili
	}
	list2["column"] = "weight"
	list2["name"] = "体重"
	list4["name"] = "体重评价"
	list4["column"] = "abnormal_weight"
	list = append(list, list2)
	list = append(list, list4)
	sql = "select count(a.id) as num from healthy_inspect a where a.class_id=" + strconv.Itoa(class_id) + " and a.body_id = " + strconv.Itoa(id) + " and height is not null"
	o.Raw(sql).QueryRow(&num)
	fmt.Println(num.Num)
	list3["bili"] = int(math.Ceil(float64(num.Num) / float64(c_num) * 100.0))
	list3["column"] = "height"
	list3["name"] = "身高"
	list5["name"] = "身高评价"
	list5["column"] = "abnormal_height"
	list = append(list, list3)
	list = append(list, list5)
	sql = "select a.id,a.theme,a.test_time,a.mechanism,a.kindergarten_id,a.types,a.project,b.class_total as total,b.class_actual as actual,b.class_rate as rate from healthy_body a left join healthy_class b on a.id = b.body_id where a.id=" + strconv.Itoa(id) + " and b.class_id=" + strconv.Itoa(class_id)
	var b Body

	if err = o.Raw(sql).QueryRow(&b); err == nil {
		project := strings.Split(b.Project, ",")
		for _, val := range project {
			list1 := make(map[string]interface{})
			cloumn := strings.Split(val, ":")
			sql := "select count(b.id) as num from healthy_inspect a left join healthy_column b on a.id= b.inspect_id where  a.class_id=" + strconv.Itoa(class_id) + " and a.body_id = " + strconv.Itoa(id) + " and b." + string(cloumn[0]) + " is not null"
			o.Raw(sql).QueryRow(&num)
			fmt.Println(num.Num)
			bili := int(math.Ceil(float64(num.Num) / float64(c_num) * 100.0))
			if bili < 0 {
				list1["bili"] = 0
			} else {
				list1["bili"] = bili
			}
			list1["column"] = cloumn[0]
			list1["name"] = string(cloumn[1])
			list = append(list, list1)
		}
		bjson, _ := json.Marshal(b)
		json.Unmarshal(bjson, &ml)
		ml["info"] = list
		return ml, nil
	}
	return nil, err
}
