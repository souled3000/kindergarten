package models

import (
	"math"
	"time"

	"github.com/astaxie/beego/orm"
)

type SlideShow struct {
	Id             int       `json:"id" orm:"column(id);auto"`
	Title          string    `json:"Title" orm:"column(Title);" description:"标题"`
	KindergartenId int       `json:"kindergarten_id" orm:"column(kindergarten_id)"`
	Content        string    `json:"Content" orm:"column(Content)" description:"内容"`
	Picture        string    `json:"picture" orm:"column(picture)" description:"图片"`
	CreatedAt      time.Time `json:"created_at" orm:"auto_now_add"`
	UpdatedAt      time.Time `json:"update_at" orm:"auto_now"`
}

func (t *SlideShow) TableName() string {
	return "slide_show"
}

func init() {
	orm.RegisterModel(new(SlideShow))
}

/*
添加轮播图
*/
func AddSlideShow(title string, content string, kindergarten_id int, picture string) (err error) {
	o := orm.NewOrm()
	m := SlideShow{Title: title, Content: content, KindergartenId: kindergarten_id, Picture: picture}
	_, err = o.Insert(&m)
	return err
}

/*
轮播图列表
*/
func GetSlideShowList(page int, prepage int, kindergarten_id int) (paginatorMap map[string]interface{}, err error) {
	var v []SlideShow
	o := orm.NewOrm()
	nums, err := o.QueryTable("slide_show").Filter("kindergarten_id", kindergarten_id).All(&v)
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
		num, err := o.QueryTable("slide_show").Filter("kindergarten_id", kindergarten_id).Limit(prepage, limit).All(&v)
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
轮播图详情
*/
func GetSlideShow(id int) (ml map[string]interface{}, err error) {
	var v []SlideShow
	o := orm.NewOrm()
	err = o.QueryTable("slide_show").Filter("Id", id).One(&v)
	if err == nil {
		paginatorMap := make(map[string]interface{})
		paginatorMap["data"] = v
		return paginatorMap, nil
	}
	return nil, err
}

/*
删除轮播图
*/
func DeleteSlideShow(id int) (err error) {
	o := orm.NewOrm()
	_, err = o.QueryTable("slide_show").Filter("id", id).Delete()
	return err
}

/*
编辑轮播图
*/
func UpdateSlideShow(id int, title string, content string, kindergarten_id int, picture string) (err error) {
	o := orm.NewOrm()
	_, err = o.QueryTable("slide_show").Filter("id", id).Update(orm.Params{
		"Content": content, "Title": title, "KindergartenId": kindergarten_id, "Picture": picture,
	})
	return err
}
