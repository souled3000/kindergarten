package models

import (
	"math"
	"time"

	"github.com/astaxie/beego/orm"
)

type FacilitiesDisplay struct {
	Id             int       `json:"id" orm:"column(id);auto;"`
	Picture        string    `json:"picture" orm:"column(picture);"; description:"图片"`
	Order          int       `json:"order" orm:"column(order);" description:"排序"`
	KindergartenId int       `json:"kindergarten_id" orm:"column(kindergarten_id)";description:"幼儿园ID"`
	CreatedAt      time.Time `json:"created_at" orm:"auto_now_add"`
	UpdatedAt      time.Time `json:"updated_at" orm:"auto_now"`
}

func (t *FacilitiesDisplay) TableName() string {
	return "facilities_display"
}

func init() {
	orm.RegisterModel(new(FacilitiesDisplay))
}

/*
添加设施
*/
func Store(order int, picture string, kindergarten_id int) (err error) {
	o := orm.NewOrm()
	facli := FacilitiesDisplay{Order: order, Picture: picture, KindergartenId: kindergarten_id}
	_, err = o.Insert(&facli)
	return err
}

/*
设施列表
*/
func GetList(page int, prepage int, kindergarten_id int) (ml map[string]interface{}, err error) {
	var v []Notice
	o := orm.NewOrm()
	nums, err := o.QueryTable("facilities_display").Filter("kindergarten_id", kindergarten_id).All(&v)
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
		num, err := o.QueryTable("notice").Filter("kindergarten_id", kindergarten_id).Limit(prepage, limit).All(&v)
		if err == nil && num > 0 {
			paginatorMap := make(map[string]interface{})
			paginatorMap["total"] = nums          //总条数
			paginatorMap["data"] = v              //分页数据
			paginatorMap["page_num"] = totalpages //总页数
			return paginatorMap, nil
		}
	}
	return nil, err

}

/*
设施详情
*/
func GetOne(id int) (ml map[string]interface{}, err error) {
	var v []FacilitiesDisplay
	o := orm.NewOrm()
	err = o.QueryTable("facilities_display").Filter("Id", id).One(&v)
	if err == nil {
		paginatorMap := make(map[string]interface{})
		paginatorMap["data"] = v
		return paginatorMap, nil
	}
	return nil, err
}

/*
删除设施
*/
func Delete(id int) (err error) {
	o := orm.NewOrm()
	_, err = o.QueryTable("facilities_display").Filter("id", id).Delete()
	return err
}

/*
编辑设施
*/
func Update(id int, picture string, order int, kindergarten_id int) (err error) {
	o := orm.NewOrm()
	_, err = o.QueryTable("facilities_display").Filter("id", id).Update(orm.Params{
		"order": order, "picture": picture, "kindergarten_id": kindergarten_id,
	})
	return err
}
