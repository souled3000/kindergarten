package models

import (
	"math"
	"time"

	"github.com/astaxie/beego/orm"
)

type KindergartenLife struct {
	Id             int       `json:"id" orm:"column(id);auto"`
	Content        string    `json:"content" orm:"column(content);size(255)" description:"内容"`
	KindergartenId int       `json:"kindergarten_id" orm:"column(kindergarten_id)"`
	Template       int8      `json:"template" orm:"column(template)" description:"模板"`
	CreatedAt      time.Time `json:"created_at" orm:"auto_now_add"`
	UpdatedAt      time.Time `json:"update_at" orm:"auto_now"`
}

func (t *KindergartenLife) TableName() string {
	return "kindergarten_life"
}

func init() {
	orm.RegisterModel(new(KindergartenLife))
}

//web-添加园内生活
func AddKindergartenLife(m *KindergartenLife) map[string]interface{} {
	o := orm.NewOrm()
	id, err := o.Insert(m)
	if err == nil {
		paginatorMap := make(map[string]interface{})
		paginatorMap["data"] = id //返回数据
		return paginatorMap
	}
	return nil
}

//web-园内生活列表
func GetKindergartenLifeList(page, prepage int) map[string]interface{} {
	var v []KindergartenLife
	o := orm.NewOrm()
	nums, err := o.QueryTable("kindergarten_life").All(&v)
	if err == nil && nums > 0 {
		//根据nums总数，和prepage每页数量 生成分页总数
		totalpages := int(math.Ceil(float64(nums) / float64(prepage))) //page总数
		if page > totalpages {
			page = totalpages
		}
		if page <= 0 {
			page = 1
		}
		limit := (page - 1) * prepage
		num, err := o.QueryTable("kindergarten_life").Limit(prepage, limit).All(&v)
		if err == nil && num > 0 {
			paginatorMap := make(map[string]interface{})
			paginatorMap["total"] = nums          //总条数
			paginatorMap["data"] = v              //分页数据
			paginatorMap["page_num"] = totalpages //总页数
			return paginatorMap
		}
	}
	return nil
}

//Web -园内生活详情
func GetKindergartenLifeInfo(id int, page, prepage int) map[string]interface{} {
	var v []KindergartenLife
	o := orm.NewOrm()
	nums, err := o.QueryTable("kindergarten_life").Filter("Id", id).Count()
	if err == nil {
		totalpages := int(math.Ceil(float64(nums) / float64(prepage))) //总页数
		if page > totalpages {
			page = totalpages
		}
		if page <= 0 {
			page = 1
		}
		limit := (page - 1) * prepage
		err := o.QueryTable("kindergarten_life").Filter("Id", id).Limit(prepage, limit).One(&v)
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

//web-删除园内生活
func DeleteKindergartenLife(id int) map[string]interface{} {
	o := orm.NewOrm()
	v := KindergartenLife{Id: id}
	// ascertain id exists in the database
	if err := o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&KindergartenLife{Id: id}); err == nil {
			paginatorMap := make(map[string]interface{})
			paginatorMap["data"] = num //返回数据
			return paginatorMap
		}
	}
	return nil
}
