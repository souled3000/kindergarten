package models

import (
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
)

type KindergartenLife struct {
	Id             int       `json:"id" orm:"column(id);auto"`
	Content        string    `json:"content" orm:"column(content);size(255)" description:"内容"`
	KindergartenId int       `json:"kindergarten_id" orm:"column(kindergarten_id)"`
	Template       int       `json:"template" orm:"column(template)" description:"模板"`
	CreatedAt      time.Time `json:"created_at" orm:"auto_now_add"`
	UpdatedAt      time.Time `json:"update_at" orm:"auto_now"`
}

func (t *KindergartenLife) TableName() string {
	return "kindergarten_life"
}

func init() {
	orm.RegisterModel(new(KindergartenLife))
}

/*
添加园内生活
*/
func AddKindergartenLife(content string, template int, kindergarten_id int, picture string, number int) (err error) {
	o := orm.NewOrm()
	photo := strings.Split(picture, ",")
	m := KindergartenLife{Content: content, Template: template, KindergartenId: kindergarten_id}
	id, err := o.Insert(&m)
	ids := strconv.FormatInt(id, 10)
	lifeId, _ := strconv.Atoi(ids)
	if err == nil {
		for _, v := range photo {
			l := LifePicture{LifeId: lifeId, Picture: v, Number: number}
			_, err = o.Insert(&l)
		}
		if err != nil {
			return err
		}
	}
	return err
}

/*
园内生活列表
*/
func GetKindergartenLifeList(page, prepage int, kindergarten_id int) (paginatorMap map[string]interface{}, err error) {
	o := orm.NewOrm()
	qb, _ := orm.NewQueryBuilder("mysql")
	// 构建查询对象
	sql := qb.Select("count(*)").From("kindergarten_life as kl").LeftJoin("life_picture as lp").
		On("kl.id = lp.life_id").Where("kl.kindergarten_id = ?").String()
	var total int64
	err = o.Raw(sql, kindergarten_id).QueryRow(&total)
	if err == nil {
		var v []orm.Params
		//根据nums总数，和prepage每页数量 生成分页总数
		totalpages := int(math.Ceil(float64(total) / float64(prepage))) //page总数
		if page > totalpages {
			page = totalpages
		}
		if page <= 0 {
			page = 1
		}
		limit := (page - 1) * prepage
		qb, _ := orm.NewQueryBuilder("mysql")
		sql := qb.Select("*").From("kindergarten_life as kl").LeftJoin("life_picture as lp").
			On("kl.id = lp.life_id").Where("kl.kindergarten_id = ?").Limit(prepage).Offset(limit).String()
		num, err := o.Raw(sql, kindergarten_id).Values(&v)
		if err == nil && num > 0 {
			paginatorMap := make(map[string]interface{})
			paginatorMap["total"] = total //总条数
			paginatorMap["data"] = v
			paginatorMap["page_num"] = totalpages //总页数
			return paginatorMap, nil
		}
	}
	return nil, err
}

/*
Web -园内生活详情
*/
func GetKindergartenLifeInfo(id int) (ml map[string]interface{}, err error) {
	o := orm.NewOrm()
	var v []orm.Params
	qb, _ := orm.NewQueryBuilder("mysql")
	sql := qb.Select("*").From("kindergarten_life as kl").LeftJoin("life_picture as lp").
		On("kl.id = lp.life_id").Where("kl.id = ?").String()
	_, err = o.Raw(sql, id).Values(&v)
	if err == nil {
		paginatorMap := make(map[string]interface{})
		paginatorMap["data"] = v
		return paginatorMap, nil
	}
	return nil, err
}

/*
web-删除园内生活
*/
func DeleteKindergartenLife(id int) (err error) {
	o := orm.NewOrm()
	v := KindergartenLife{Id: id}
	if err = o.Read(&v); err == nil {
		if _, err = o.Delete(&KindergartenLife{Id: id}); err == nil {
			return nil
		}
	}
	return err
}

/*
编辑园内生活
*/
func UpdateKL(id int, content string, template string, kindergarten_id int) (err error) {
	o := orm.NewOrm()
	_, err = o.QueryTable("kindergarten_life").Filter("id", id).Update(orm.Params{
		"content": content, "template": template, "kindergarten_id": kindergarten_id,
	})
	return err
}
