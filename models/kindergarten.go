package models

import (
	"errors"
	"fmt"
	"math"
	"reflect"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
)

type Kindergarten struct {
	Id               int       `json:"kindergarten_id" orm:"column(kindergarten_id);auto" description:"编号"`
	Name             string    `json:"name" orm:"column(name);size(50)" description:"幼儿园名称"`
	LicenseNo        int       `json:"license_no" orm:"column(license_no)" description:"执照号"`
	KinderGrade      string    `json:"kinder_grade" orm:"column(kinder_grade);size(45)" description:"幼儿园级别"`
	KinderChildNo    int       `json:"kinder_child_no" orm:"column(kinder_child_no)" description:"分校数"`
	Address          string    `json:"address" orm:"column(address);size(50)" description:"地址"`
	TenantId         int       `json:"tenant_id" orm:"column(tenant_id)" description:"租户，企业编号"`
	Status           int8      `json:"status" orm:"column(status)" description:"状态：0:正常，1:删除"`
	Introduce        string    `json:"introduce" orm:"column(introduce);size(255)" description:"幼儿园介绍"`
	IntroducePicture string    `json:"introduce_picture" orm:"column(introduce_picture);size(255)" description:"幼儿园介绍图"`
	CreatedAt        time.Time `json:"created_at" orm:"auto_now_add"`
	UpdatedAt        time.Time `json:"updated_at" orm:"auto_now"`
	DeletedAt        time.Time `json:"deleted_at" orm:"column(deleted_at);type(datetime);null"`
}

func (t *Kindergarten) TableName() string {
	return "kindergarten"
}

func init() {
	orm.RegisterModel(new(Kindergarten))
}

//web-幼儿园介绍详情
func GetKindergartenById(id int, page, prepage int) map[string]interface{} {
	var v []Kindergarten
	o := orm.NewOrm()
	nums, err := o.QueryTable("kindergarten").Filter("Id", id).Count()
	if err == nil {
		totalpages := int(math.Ceil(float64(nums) / float64(prepage))) //总页数
		if page > totalpages {
			page = totalpages
		}
		if page <= 0 {
			page = 1
		}
		limit := (page - 1) * prepage
		err := o.QueryTable("kindergarten").Filter("Id", id).Limit(prepage, limit).One(&v)
		if err == nil {
			paginatorMap := make(map[string]interface{})
			paginatorMap["total"] = nums          //总条数
			paginatorMap["data"] = v              //返回数据
			paginatorMap["page_num"] = totalpages //总页数
			return paginatorMap
		}
	}
	return nil
}

// GetAllKindergarten retrieves all Kindergarten matches certain condition. Returns empty list if
// no records exist
func GetAllKindergarten(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Kindergarten))
	// query k=v
	for k, v := range query {
		// rewrite dot-notation to Object__Attribute
		k = strings.Replace(k, ".", "__", -1)
		if strings.Contains(k, "isnull") {
			qs = qs.Filter(k, (v == "true" || v == "1"))
		} else {
			qs = qs.Filter(k, v)
		}
	}
	// order by:
	var sortFields []string
	if len(sortby) != 0 {
		if len(sortby) == len(order) {
			// 1) for each sort field, there is an associated order
			for i, v := range sortby {
				orderby := ""
				if order[i] == "desc" {
					orderby = "-" + v
				} else if order[i] == "asc" {
					orderby = v
				} else {
					return nil, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
			qs = qs.OrderBy(sortFields...)
		} else if len(sortby) != len(order) && len(order) == 1 {
			// 2) there is exactly one order, all the sorted fields will be sorted by this order
			for _, v := range sortby {
				orderby := ""
				if order[0] == "desc" {
					orderby = "-" + v
				} else if order[0] == "asc" {
					orderby = v
				} else {
					return nil, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
		} else if len(sortby) != len(order) && len(order) != 1 {
			return nil, errors.New("Error: 'sortby', 'order' sizes mismatch or 'order' size is not 1")
		}
	} else {
		if len(order) != 0 {
			return nil, errors.New("Error: unused 'order' fields")
		}
	}

	var l []Kindergarten
	qs = qs.OrderBy(sortFields...)
	if _, err = qs.Limit(limit, offset).All(&l, fields...); err == nil {
		if len(fields) == 0 {
			for _, v := range l {
				ml = append(ml, v)
			}
		} else {
			// trim unused fields
			for _, v := range l {
				m := make(map[string]interface{})
				val := reflect.ValueOf(v)
				for _, fname := range fields {
					m[fname] = val.FieldByName(fname).Interface()
				}
				ml = append(ml, m)
			}
		}
		return ml, nil
	}
	return nil, err
}

// UpdateKindergarten updates Kindergarten by Id and returns error if
// the record to be updated doesn't exist
func UpdateKindergartenById(m *Kindergarten) (err error) {
	o := orm.NewOrm()
	v := Kindergarten{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteKindergarten deletes Kindergarten by Id and returns error if
// the record to be deleted doesn't exist
func DeleteKindergarten(id int) (err error) {
	o := orm.NewOrm()
	v := Kindergarten{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&Kindergarten{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
