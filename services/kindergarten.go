package services

import (
	"encoding/json"
	"fmt"

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
