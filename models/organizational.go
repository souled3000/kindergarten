package models

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
)

type Organizational struct {
	Id             int       `orm:"column(id);auto"`
	KindergartenId int       `orm:"column(kindergarten_id)" description:"幼儿园id"`
	ParentId       int       `orm:"column(parent_id)" description:"父级id"`
	Name           string    `orm:"column(name);size(20)" description:"组织架构名字"`
	IsFixed        int8      `orm:"column(is_fixed)" description:"是否固定的：0不是，1是"`
	Level          int8      `orm:"column(level)" description:"等级"`
	ParentIds      string    `orm:"column(parent_ids);size(50)" description:"父级所有id"`
	Type           int8      `orm:"column(type)" description:"类型：0普通，1管理层，2年级组"`
	ClassType      int8      `orm:"column(class_type)" description:"班级类型：1小班，2中班，3大班"`
	CreatedAt      time.Time `orm:"column(created_at);type(datetime)" description:"添加时间"`
	UpdatedAt      time.Time `orm:"column(updated_at);type(datetime)" description:"修改时间"`
}

func (t *Organizational) TableName() string {
	return "organizational"
}

func init() {
	orm.RegisterModel(new(Organizational))
}

// AddOrganizational insert a new Organizational into database and returns
// last inserted Id on success.
func AddOrganizational(m *Organizational) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetOrganizationalById retrieves Organizational by Id. Returns error if
// Id doesn't exist
func GetOrganizationalById(id int) (v *Organizational, err error) {
	o := orm.NewOrm()
	v = &Organizational{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllOrganizational retrieves all Organizational matches certain condition. Returns empty list if
// no records exist
func GetAllOrganizational(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Organizational))
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

	var l []Organizational
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

// UpdateOrganizational updates Organizational by Id and returns error if
// the record to be updated doesn't exist
func UpdateOrganizationalById(m *Organizational) (err error) {
	o := orm.NewOrm()
	v := Organizational{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteOrganizational deletes Organizational by Id and returns error if
// the record to be deleted doesn't exist
func DeleteOrganizational(id int) (err error) {
	o := orm.NewOrm()
	v := Organizational{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&Organizational{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
