package services

import (
	"encoding/json"
	"fmt"

	"kindergarten-service-go/models"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/hprose/hprose-golang/rpc"
)

type KindergartenServer struct {
}

func (c *KindergartenServer) Init() {
	server := rpc.NewHTTPService()
	server.AddAllMethods(&KindergartenServer{})
	beego.Handler("/rpc/kindergarten", server)
}

//班级信息
func (c *KindergartenServer) GetKg(user_id int, kindergarten_id int) (value map[string]interface{}, err error) {
	o := orm.NewOrm()
	var v []orm.Params
	var kinder []orm.Params
	var permission []orm.Params
	//幼儿园信息
	qb, _ := orm.NewQueryBuilder("mysql")
	sql := qb.Select("k.name").From("kindergarten as k").Where("kindergarten_id = ?").String()
	_, err = o.Raw(sql, kindergarten_id).Values(&kinder)
	fmt.Println(kinder[0]["name"])
	//权限信息
	qb, _ = orm.NewQueryBuilder("mysql")
	sql = qb.Select("p.identification").From("user_permission as up").LeftJoin("permission as p").
		On("up.permission_id = p.id").Where("up.user_id = ?").String()
	_, err = o.Raw(sql, user_id).Values(&permission)
	//班级信息
	qb, _ = orm.NewQueryBuilder("mysql")
	sql = qb.Select("o.id as class_id", "o.name as class_name").From("teacher as t").LeftJoin("organizational_member as om").
		On("t.teacher_id = om.member_id").LeftJoin("organizational as o").
		On("om.organizational_id = o.id").Where("t.user_id = ?").And("o.type = 2").And("o.level = 3").String()
	_, err = o.Raw(sql, user_id).Values(&v)
	if err == nil {
		if v == nil {
			value := make(map[string]interface{})
			value["kindergarten_name"] = kinder[0]["name"]
			jsons, _ := json.Marshal(permission)
			value["permission"] = jsons
			return value, nil
		} else {
			value := v[0]
			value["kindergarten_name"] = kinder[0]["name"]
			jsons, _ := json.Marshal(permission)
			value["permission"] = jsons
			return value, nil
		}
	}
	return nil, err
}

//班级成员
func (c *KindergartenServer) GetMember(organizational_id int) (ml map[string]interface{}, err error) {
	o := orm.NewOrm()
	var student []orm.Params
	var teacher []orm.Params
	qb, _ := orm.NewQueryBuilder("mysql")
	sql := qb.Select("t.*").From("organizational_member as om").LeftJoin("teacher as t").
		On("om.member_id = t.teacher_id").Where("om.organizational_id = ?").And("om.type = 0").String()
	_, err = o.Raw(sql, organizational_id).Values(&teacher)

	qb, _ = orm.NewQueryBuilder("mysql")
	sql = qb.Select("s.*").From("organizational_member as om").LeftJoin("student as s").
		On("om.member_id = s.student_id").Where("om.organizational_id = ?").And("om.type = 1").String()
	_, err = o.Raw(sql, organizational_id).Values(&student)
	for _, v := range teacher {
		student = append(student, v)
	}
	if err == nil {
		paginatorMap := make(map[string]interface{})
		paginatorMap["data"] = student
		return paginatorMap, nil
	}
	return nil, err
}

func (c *KindergartenServer) GetClass(kindergarten_id int) (ml map[string]interface{}) {
	v := models.GetGroupPermission(kindergarten_id)
	ml = make(map[string]interface{})
	ml["data"] = v
	return ml
}

/*
宝宝所在班级
*/
func (c *KindergartenServer) GetBabyClass(Baby_id int) (ml map[string]interface{}, err error) {
	o := orm.NewOrm()
	var v []orm.Params
	qb, _ := orm.NewQueryBuilder("mysql")
	sql := qb.Select("s.name", "s.student_id", "o.name as class_name", "o.id as class_id").From("student as s").LeftJoin("organizational_member as om").
		On("s.student_id = om.member_id").LeftJoin("organizational as o").
		On("om.organizational_id = o.id").Where("s.baby_id = ?").And("status = 1").And("isnull(deleted_at)").String()
	_, err = o.Raw(sql, Baby_id).Values(&v)
	if err == nil {
		ml = make(map[string]interface{})
		ml["data"] = v
		return ml, nil
	}
	return nil, err
}
