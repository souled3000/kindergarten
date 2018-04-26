package models

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
)

type StudentExamination struct {
	Id             int       `orm:"column(examination_id);pk" description:"体检编号"`
	KindergartenId int       `orm:"column(kindergarten_id);null" description:"幼儿园ID"`
	StudentId      int       `orm:"column(student_id)" description:"学生ID"`
	Userid         int       `orm:"column(userid)" description:"用户ID"`
	Height         float64   `orm:"column(height);null" description:"身高"`
	Weight         float64   `orm:"column(weight);null" description:"体重"`
	CreatedAt      time.Time `orm:"column(created_at);type(timestamp);null"`
	UpdatedAt      time.Time `orm:"column(updated_at);type(timestamp);null"`
	Status         int8      `orm:"column(status);null" description:"状态：0:正常，1:删除"`
}

func (t *StudentExamination) TableName() string {
	return "student_ examination"
}

func init() {
	orm.RegisterModel(new(StudentExamination))
}

// AddStudentExamination insert a new StudentExamination into database and returns
// last inserted Id on success.
func AddStudentExamination(m *StudentExamination) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetStudentExaminationById retrieves StudentExamination by Id. Returns error if
// Id doesn't exist
func GetStudentExaminationById(id int) (v *StudentExamination, err error) {
	o := orm.NewOrm()
	v = &StudentExamination{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllStudentExamination retrieves all StudentExamination matches certain condition. Returns empty list if
// no records exist
func GetAllStudentExamination(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(StudentExamination))
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

	var l []StudentExamination
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

// UpdateStudentExamination updates StudentExamination by Id and returns error if
// the record to be updated doesn't exist
func UpdateStudentExaminationById(m *StudentExamination) (err error) {
	o := orm.NewOrm()
	v := StudentExamination{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteStudentExamination deletes StudentExamination by Id and returns error if
// the record to be deleted doesn't exist
func DeleteStudentExamination(id int) (err error) {
	o := orm.NewOrm()
	v := StudentExamination{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&StudentExamination{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
