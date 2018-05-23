package models

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
)

type Kinship struct {
	Id int `json:"kinship_id" orm:"column(kinship_id);auto" description:"编号"`
	//StudentId          int       `json:"student_id" orm:"column(student_id)" description:"学生序号"`
	UserId             int       `json:"user_id" orm:"column(user_id)" description:"用户ID"`
	Type               int8      `json:"type" orm:"column(type)" description:"类型，1紧急联系人，2监护人"`
	Relation           string    `json:"relation" orm:"column(relation);size(10)" description:"关系"`
	Name               string    `json:"name" orm:"column(name);size(20)" description:"名字"`
	UnitName           string    `json:"unit_name" orm:"column(unit_name);size(50)" description:"单位名称"`
	ContactInformation string    `json:"contact_information" orm:"column(contact_information);size(11)" description:"联系方式"`
	CreatedAt          time.Time `json:"created_at" orm:"auto_now_add"`
	UpdatedAt          time.Time `json:"updated_at" orm:"auto_now"`
	Student            *Student  `json:"-" orm:"rel(fk)"`
}

func (t *Kinship) TableName() string {
	return "kinship"
}

func init() {
	orm.RegisterModel(new(Kinship))
}

// AddKinship insert a new Kinship into database and returns
// last inserted Id on success.
func AddKinship(m *Kinship) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetKinshipById retrieves Kinship by Id. Returns error if
// Id doesn't exist
func GetKinshipById(id int) (v *Kinship, err error) {
	o := orm.NewOrm()
	v = &Kinship{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllKinship retrieves all Kinship matches certain condition. Returns empty list if
// no records exist
func GetAllKinship(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Kinship))
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

	var l []Kinship
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

// UpdateKinship updates Kinship by Id and returns error if
// the record to be updated doesn't exist
func UpdateKinshipById(m *Kinship) (err error) {
	o := orm.NewOrm()
	v := Kinship{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteKinship deletes Kinship by Id and returns error if
// the record to be deleted doesn't exist
func DeleteKinship(id int) (err error) {
	o := orm.NewOrm()
	v := Kinship{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&Kinship{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
