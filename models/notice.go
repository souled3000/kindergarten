package models

import (
	"math"
	"time"

	"github.com/astaxie/beego/orm"
)

type Notice struct {
	Id             int       `json:"id";orm:"column(id);auto;json:"id";"`
	Title          string    `json:"title";orm:"column(title);size(50)"; description:"标题"`
	Content        string    `json:"content";orm:"column(content);json:"content";size(255)" description:"公告内容"`
	KindergartenId int       `json:"kindergarten_id";orm:"column(kindergarten_id)";json:"kindergarten_id"; description:"幼儿园ID"`
	CreatedAt      time.Time `json:"created_at";orm:"column(created_at);auto_now_add;json:"created_at";type(datetime)"`
	UpdatedAt      time.Time `json:"updated_at";orm:"column(updated_at);auto_now;type(datetime)";`
}

func (t *Notice) TableName() string {
	return "notice"
}

func init() {
	orm.RegisterModel(new(Notice))
}

//公告添加
func AddNotice(m *Notice) map[string]interface{} {
	o := orm.NewOrm()
	id, err := o.Insert(m)
	if err == nil {
		paginatorMap := make(map[string]interface{})
		paginatorMap["data"] = id //返回数据
		return paginatorMap
	}
	return nil
}

//公告列表
func GetNoticeList(page, prepage int) map[string]interface{} {
	var v []Notice
	o := orm.NewOrm()
	nums, err := o.QueryTable("notice").All(&v)
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
		num, err := o.QueryTable("notice").Limit(prepage, limit).All(&v)
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

//删除公告
func DeleteNotice(id int) map[string]interface{} {
	o := orm.NewOrm()
	v := Notice{Id: id}
	// ascertain id exists in the database
	if err := o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&Notice{Id: id}); err == nil {
			paginatorMap := make(map[string]interface{})
			paginatorMap["data"] = num //返回数据
			return paginatorMap
		}
	}
	return nil
}
