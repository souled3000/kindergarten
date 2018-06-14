package models

import (
	"fmt"
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
	var User *UserService
	client := rpc.NewHTTPClient(beego.AppConfig.String("ONE_MORE_USER_SERVER"))
	client.UseService(&User)
	uk, _ := User.GetUKByUserId(user_id)
	fmt.Println(uk)
	return nil
}
