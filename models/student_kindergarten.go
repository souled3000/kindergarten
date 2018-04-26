package models

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
)

type StudentKindergarten struct {
	Id               int       `orm:"column(student_kindergarten_id);auto"`
	StudentId        int       `orm:"column(student_id)" description:"学生ID"`
	KindergartenId   int       `orm:"column(kindergarten_id)" description:"幼儿园ID"`
	EnterGardenTime  time.Time `orm:"column(enter_garden_time);type(date)" description:"入园时间"`
	Weight           float64   `orm:"column(weight)"`
	Height           float64   `orm:"column(height)"`
	LeaveTime        time.Time `orm:"column(leave_time);type(datetime)" description:"离园时间"`
	KindergartenName string    `orm:"column(kindergarten_name);size(50)" description:"幼儿园名称"`
	ClassName        string    `orm:"column(class_name);size(50)" description:"班级名称"`
	IsExist          int8      `orm:"column(is_exist)" description:"是否在系统内"`
	Identity         string    `orm:"column(identity);size(191);null"`
	InKindergarten   string    `orm:"column(in_kindergarten);size(191);null" description:"目前是否在园"`
	Status           int8      `orm:"column(status)" description:"状态：0:正常，1:删除"`
	CreatedAt        time.Time `orm:"column(created_at);type(timestamp);null"`
	UpdatedAt        time.Time `orm:"column(updated_at);type(timestamp);null"`
}

func (t *StudentKindergarten) TableName() string {
	return "student_kindergarten"
}

func init() {
	orm.RegisterModel(new(StudentKindergarten))
}

// AddStudentKindergarten insert a new StudentKindergarten into database and returns
// last inserted Id on success.
func AddStudentKindergarten(m *StudentKindergarten) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetStudentKindergartenById retrieves StudentKindergarten by Id. Returns error if
// Id doesn't exist
func GetStudentKindergartenById(id int) (v *StudentKindergarten, err error) {
	o := orm.NewOrm()
	v = &StudentKindergarten{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllStudentKindergarten retrieves all StudentKindergarten matches certain condition. Returns empty list if
// no records exist
func GetAllStudentKindergarten(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(StudentKindergarten))
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

	var l []StudentKindergarten
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

// UpdateStudentKindergarten updates StudentKindergarten by Id and returns error if
// the record to be updated doesn't exist
func UpdateStudentKindergartenById(m *StudentKindergarten) (err error) {
	o := orm.NewOrm()
	v := StudentKindergarten{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteStudentKindergarten deletes StudentKindergarten by Id and returns error if
// the record to be deleted doesn't exist
func DeleteStudentKindergarten(id int) (err error) {
	o := orm.NewOrm()
	v := StudentKindergarten{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&StudentKindergarten{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
