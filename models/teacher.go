package models

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
)

type Teacher struct {
	Id                         int       `orm:"column(teacher_id);auto" description:"自增id"`
	Name                       string    `orm:"column(name);size(20)" description:"姓名"`
	Age                        int8      `orm:"column(age)" description:"年龄"`
	Sex                        int8      `orm:"column(sex)" description:"性别 0男  1女"`
	Avatar                     string    `orm:"column(avatar);size(150)" description:"头像"`
	Number                     string    `orm:"column(number);size(20)" description:"教职工编号"`
	NationOrReligion           string    `orm:"column(nation_or_religion);size(10)" description:"民族或宗教"`
	NativePlace                string    `orm:"column(native_place);size(20)" description:"籍贯"`
	UserId                     int       `orm:"column(user_id)" description:"用户id"`
	ClassInfo                  string    `orm:"column(class_info);size(10)" description:"班级信息"`
	Phone                      string    `orm:"column(phone);size(11)" description:"联系电话"`
	EnterGardenTime            time.Time `orm:"column(enter_garden_time);type(date)" description:"进入本园时间"`
	EnterJobTime               time.Time `orm:"column(enter_job_time);type(date)" description:"参加工作时间"`
	KindergartenId             int       `orm:"column(kindergarten_id)" description:"幼儿园id"`
	Address                    string    `orm:"column(address);size(191)" description:"住址"`
	IdNumber                   string    `orm:"column(id_number);size(18)" description:"身份证号"`
	EmergencyContact           string    `orm:"column(emergency_contact);size(20)" description:"紧急联系人"`
	EmergencyContactPhone      string    `orm:"column(emergency_contact_phone);size(11)" description:"紧急联系人电话"`
	Post                       string    `orm:"column(post);size(10)" description:"职务"`
	Source                     string    `orm:"column(source);size(191)" description:"来源"`
	TeacherCertificationNumber string    `orm:"column(teacher_certification_number);size(20)" description:"教师资格认证编号"`
	TeacherCertificationStatus int8      `orm:"column(teacher_certification_status)" description:"教师资格证书状态，是否认证"`
	Status                     int8      `orm:"column(status)" description:"状态：0未分班，1已分班，2离职"`
	CreatedAt                  time.Time `orm:"column(created_at);type(datetime)"`
	UpdatedAt                  time.Time `orm:"column(updated_at);type(datetime)"`
	DeletedAt                  time.Time `orm:"column(deleted_at);type(datetime);null"`
}

func (t *Teacher) TableName() string {
	return "teacher"
}

func init() {
	orm.RegisterModel(new(Teacher))
}

// AddTeacher insert a new Teacher into database and returns
// last inserted Id on success.
func AddTeacher(m *Teacher) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetTeacherById retrieves Teacher by Id. Returns error if
// Id doesn't exist
func GetTeacherById(id int) (v *Teacher, err error) {
	o := orm.NewOrm()
	v = &Teacher{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllTeacher retrieves all Teacher matches certain condition. Returns empty list if
// no records exist
func GetAllTeacher(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Teacher))
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

	var l []Teacher
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

// UpdateTeacher updates Teacher by Id and returns error if
// the record to be updated doesn't exist
func UpdateTeacherById(m *Teacher) (err error) {
	o := orm.NewOrm()
	v := Teacher{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteTeacher deletes Teacher by Id and returns error if
// the record to be deleted doesn't exist
func DeleteTeacher(id int) (err error) {
	o := orm.NewOrm()
	v := Teacher{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&Teacher{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
