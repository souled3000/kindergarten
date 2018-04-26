package models

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
)

type KindergartenFolder struct {
	Id             int       `orm:"column(kindergarten_folder_id);auto" description:"幼儿园文件夹编号"`
	KindergartenId int       `orm:"column(kindergarten_id)" description:"幼儿园编号"`
	FolderId       int       `orm:"column(folder_id)" description:"文件夹ID"`
	FolderName     string    `orm:"column(folder_name);size(10)" description:"目录名称"`
	Type           int8      `orm:"column(type)" description:"类型，1上传，2小班，3中班，4大班，5分享中心"`
	Status         int8      `orm:"column(status)" description:"状态,0正常,1删除"`
	CreatedAt      time.Time `orm:"column(created_at);type(timestamp);null"`
	UpdatedAt      time.Time `orm:"column(updated_at);type(timestamp);null"`
}

func (t *KindergartenFolder) TableName() string {
	return "kindergarten_folder"
}

func init() {
	orm.RegisterModel(new(KindergartenFolder))
}

// AddKindergartenFolder insert a new KindergartenFolder into database and returns
// last inserted Id on success.
func AddKindergartenFolder(m *KindergartenFolder) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetKindergartenFolderById retrieves KindergartenFolder by Id. Returns error if
// Id doesn't exist
func GetKindergartenFolderById(id int) (v *KindergartenFolder, err error) {
	o := orm.NewOrm()
	v = &KindergartenFolder{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllKindergartenFolder retrieves all KindergartenFolder matches certain condition. Returns empty list if
// no records exist
func GetAllKindergartenFolder(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(KindergartenFolder))
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

	var l []KindergartenFolder
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

// UpdateKindergartenFolder updates KindergartenFolder by Id and returns error if
// the record to be updated doesn't exist
func UpdateKindergartenFolderById(m *KindergartenFolder) (err error) {
	o := orm.NewOrm()
	v := KindergartenFolder{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteKindergartenFolder deletes KindergartenFolder by Id and returns error if
// the record to be deleted doesn't exist
func DeleteKindergartenFolder(id int) (err error) {
	o := orm.NewOrm()
	v := KindergartenFolder{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&KindergartenFolder{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
