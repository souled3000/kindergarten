package models

import (
	"encoding/json"
	"errors"
	"math"
	"time"

	"github.com/astaxie/beego"

	"github.com/astaxie/beego/orm"
	"github.com/hprose/hprose-golang/rpc"
)

type Kindergarten struct {
	Id               int       `json:"kindergarten_id" orm:"column(kindergarten_id);auto" description:"编号"`
	Name             string    `json:"name" orm:"column(name);size(50)" description:"幼儿园名称"`
	LicenseNo        int       `json:"license_no" orm:"column(license_no)" description:"执照号"`
	KinderGrade      string    `json:"kinder_grade" orm:"column(kinder_grade);size(45)" description:"幼儿园级别"`
	KinderChildNo    int       `json:"kinder_child_no" orm:"column(kinder_child_no)" description:"分校数"`
	Address          string    `json:"address" orm:"column(address);size(50)" description:"地址"`
	TenantId         int       `json:"tenant_id" orm:"column(tenant_id)" description:"租户，企业编号"`
	Status           int8      `json:"status" orm:"column(status)" description:"状态：0:正常，1:删除"`
	Introduce        string    `json:"introduce" orm:"column(introduce);size(255)" description:"幼儿园介绍"`
	IntroducePicture string    `json:"introduce_picture" orm:"column(introduce_picture);size(255)" description:"幼儿园介绍图"`
	CreatedAt        time.Time `json:"created_at" orm:"auto_now_add"`
	UpdatedAt        time.Time `json:"updated_at" orm:"auto_now"`
	DeletedAt        time.Time `json:"deleted_at" orm:"column(deleted_at);type(datetime);null"`
}

func (t *Kindergarten) TableName() string {
	return "kindergarten"
}

func init() {
	orm.RegisterModel(new(Kindergarten))
}

/*
web-幼儿园介绍详情
*/
func GetKindergartenById(id int, page, prepage int) map[string]interface{} {
	var v []Kindergarten
	o := orm.NewOrm()
	nums, err := o.QueryTable("kindergarten").Filter("Id", id).Count()
	if err == nil {
		totalpages := int(math.Ceil(float64(nums) / float64(prepage))) //总页数
		if page > totalpages {
			page = totalpages
		}
		if page <= 0 {
			page = 1
		}
		limit := (page - 1) * prepage
		err := o.QueryTable("kindergarten").Filter("Id", id).Limit(prepage, limit).One(&v)
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

/*
oms-设置园长
*/
func AddPrincipal(user_id int, kindergarten_id int, role int) error {
	o := orm.NewOrm()
	o.Begin()
	var or []orm.Params
	var t []orm.Params
	var User *UserService
	client := rpc.NewHTTPClient(beego.AppConfig.String("ONE_MORE_USER_SERVER"))
	client.UseService(&User)
	qb, _ := orm.NewQueryBuilder("mysql")
	sql := qb.Select("*").From("organizational").Where("kindergarten_id = ?").And("type = 1").And("is_fixed = 1").String()
	_, err := o.Raw(sql, kindergarten_id).Values(&or)
	if err == nil {
		if or == nil {
			err = errors.New("幼儿园没有组织")
			return err
		}
	}
	uk, _ := User.GetUKByUserId(user_id)
	if uk == nil {
		User.CreateUK(user_id, kindergarten_id, role)
	} else {
		uks, _ := json.Marshal(uk)
		var ukss map[string]interface{}
		json.Unmarshal(uks, &ukss)
		uk_id := ukss["id"].(float64)
		ukId := int(uk_id)
		err := User.UpdateByUkId(ukId, user_id, kindergarten_id, role)
		if err != nil {
			err := errors.New("用户不存在")
			return err
		}
	}
	if role == 1 || role == 5 || role == 6 {
		userInfo, err := User.GetOneByUserId(user_id)
		if err != nil {
			return err
		}
		userinfo, _ := json.Marshal(userInfo)
		var user map[string]interface{}
		json.Unmarshal(userinfo, &user)
		_, err = o.QueryTable("student").Filter("user_id", user_id).Filter("status", 0).Update(orm.Params{
			"status": 1,
		})
		if err != nil {
			o.Rollback()
			return err
		}
		if err == nil {
			o.QueryTable("teacher").Filter("user_id", user_id).Filter("status", 0).Update(orm.Params{
				"status": 1,
			})
			if err != nil {
				o.Rollback()
				return err
			}
		}
		if role == 1 {
			qb, _ := orm.NewQueryBuilder("mysql")
			sql := qb.Select("t.teacher_id").From("teacher as t").Where("user_id = ?").And("status = 0").And("kindergarten_id = ?").String()
			_, err := o.Raw(sql, user_id, kindergarten_id).Values(&t)
			if err == nil {
				if t == nil {
					sql := "insert into teacher set user_id = ?,name = ?,phone = ?,sex = ?,age = ?,address = ?,kindergarten_id = ?,avatar = ?"
					id, err := o.Raw(sql, user_id, user["name"].(string), user["phone"].(string), user["sex"], user["age"], user["address"].(string), kindergarten_id, user["avatar"].(string)).Exec()
					if err != nil {
						o.Rollback()
						return err
					}
					if err == nil {
						teacherId, _ := id.LastInsertId()
						for _, value := range or {
							sql := "insert into organizational_member set organizational_id = ?,member_id = ?,is_principal = ?,type = ?"
							o.Raw(sql, value["id"], teacherId, 1, 0).Exec()
						}
					} else {
						o.Rollback()
						return err
					}
				}
			}
		} else if role == 5 {
			sql := "insert into teacher set user_id = ?,name = ?,phone = ?,sex = ?,age = ?,address = ?,kindergarten_id = ?,avatar = ?"
			_, err := o.Raw(sql, user_id, user["name"].(string), user["phone"].(string), user["sex"], user["age"], user["address"].(string), kindergarten_id, user["avatar"].(string)).Exec()
			if err != nil {
				o.Rollback()
				return err
			}
		} else {
			sql := "insert into student set user_id = ?,name = ?,phone = ?,sex = ?,age = ?,address = ?,kindergarten_id = ?,avatar = ?"
			_, err := o.Raw(sql, user_id, user["name"].(string), user["phone"].(string), user["sex"], user["age"], user["address"].(string), kindergarten_id, user["avatar"].(string)).Exec()
			if err != nil {
				o.Rollback()
				return err
			}
		}
	}
	if err == nil {
		o.Commit()
		return nil
	}
	return err
}

/*
oms幼儿园列表
*/
func GetAll(page int, prepage int, search string) map[string]interface{} {
	o := orm.NewOrm()
	qb, _ := orm.NewQueryBuilder("mysql")
	var condition []interface{}
	where := "1=1 "
	if search != "" {
		where += " AND name like ?"
		condition = append(condition, "%"+search+"%")
	}
	// 构建查询对象
	sql := qb.Select("count(*)").From("kindergarten").Where(where).String()
	var total int64
	err := o.Raw(sql, condition).QueryRow(&total)
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
		sql := qb.Select("*").From("kindergarten").Where(where).Limit(prepage).Offset(limit).String()
		num, err := o.Raw(sql, condition).Values(&v)
		if err == nil && num > 0 {
			paginatorMap := make(map[string]interface{})
			paginatorMap["total"] = total         //总条数
			paginatorMap["data"] = v              //分页数据
			paginatorMap["page_num"] = totalpages //总页数
			return paginatorMap
		}
	}
	return nil
}

/*
学生姓名搜索班级
*/
func StudentClass(page int, prepage int, name string) map[string]interface{} {
	o := orm.NewOrm()
	qb, _ := orm.NewQueryBuilder("mysql")
	// 构建查询对象
	sql := qb.Select("count(*)").From("organizational_member as om").LeftJoin("student as t").
		On("om.member_id = t.student_id").LeftJoin("organizational as o").
		On("o.id = om.organizational_id").Where("t.name = ?").And("om.type = 1").String()
	var total int64
	err := o.Raw(sql, name).QueryRow(&total)
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
		sql := qb.Select("o.id as class_id", "o.name as class_name").From("organizational_member as om").LeftJoin("student as t").
			On("om.member_id = t.student_id").LeftJoin("organizational as o").
			On("o.id = om.organizational_id").Where("t.name = ?").And("om.type = 1").Limit(prepage).Offset(limit).String()
		num, err := o.Raw(sql, name).Values(&v)
		if err == nil && num > 0 {
			paginatorMap := make(map[string]interface{})
			paginatorMap["total"] = total         //总条数
			paginatorMap["data"] = v              //分页数据
			paginatorMap["page_num"] = totalpages //总页数
			return paginatorMap
		}
	}
	return nil
}
