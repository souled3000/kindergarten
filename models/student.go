package models

import (
	"math"
	"time"

	"github.com/astaxie/beego/orm"
)

type Student struct {
	Id               int       `orm:"column(student_id);auto"`
	Name             string    `orm:"column(name);size(20)" description:"姓名"`
	Age              int8      `orm:"column(age)" description:"年龄"`
	Sex              int8      `orm:"column(sex)" description:"性别 0男 1女"`
	NativePlace      string    `orm:"column(native_place);size(20)" description:"籍贯"`
	NationOrReligion string    `orm:"column(nation_or_religion);size(20)" description:"民族或宗教"`
	Number           string    `orm:"column(number);size(11)" description:"学号"`
	ClassInfo        string    `orm:"column(class_info);size(20)" description:"所在班级"`
	Address          string    `orm:"column(address);size(50)" description:"住址"`
	Avatar           string    `orm:"column(avatar);size(150)" description:"头像"`
	Status           int8      `orm:"column(status)" description:"状态 0未分班 1已分班 2离园"`
	UserId           int       `orm:"column(user_id)" description:"用户ID"`
	KindergartenId   int       `orm:"column(kindergarten_id)" description:"幼儿园ID"`
	Phone            string    `orm:"column(phone);size(11)" description:"手机号"`
	HealthStatus     string    `orm:"column(health_status);size(150)" description:"健康状况，多个以逗号隔开"`
	Hobby            string    `orm:"column(hobby);size(150)" description:"兴趣爱好，多个以逗号隔开"`
	CreatedAt        time.Time `orm:"column(created_at);type(datetime)"`
	UpdatedAt        time.Time `orm:"column(updated_at);type(datetime)"`
	DeletedAt        time.Time `orm:"column(deleted_at);type(datetime);null"`
}

func (t *Student) TableName() string {
	return "student"
}

func init() {
	orm.RegisterModel(new(Student))
}

//学生列表
func GetStudent(id int, status int, search string, page int, prepage int) map[string]interface{} {
	var condition []interface{}
	where := "1=1 "
	if id == 0 {
		where += " AND s.kindergarten_id = ?"
	} else {
		where += " AND s.kindergarten_id = ?"
		condition = append(condition, id)
	}
	if status != -1 {
		where += " AND s.status = ?"
		condition = append(condition, status)
	}
	if search != "" {
		where += " AND s.name like ?"
		condition = append(condition, "%"+search+"%")
	}
	o := orm.NewOrm()
	qb, _ := orm.NewQueryBuilder("mysql")

	// 构建查询对象
	sql := qb.Select("count(*)").From("student as s").LeftJoin("organizational_member as om").
		On("s.student_id = om.member_id").LeftJoin("organizational as o").
		On("om.organizational_id = o.id").Where(where).And("isnull(deleted_at)").String()
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
		sql := qb.Select("s.name", "s.avatar", "s.student_id", "s.number", "s.phone", "o.name as class").From("student as s").LeftJoin("organizational_member as om").
			On("s.student_id = om.member_id").LeftJoin("organizational as o").
			On("om.organizational_id = o.id").Where(where).And("isnull(deleted_at)").Limit(prepage).Offset(limit).String()
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

//学生班级列表
func GetStudentClass(id int, class_type int, page int, prepage int) map[string]interface{} {
	var condition []interface{}
	where := "1=1 "
	if id == 0 {
		where += " AND s.kindergarten_id = ?"
	} else {
		where += " AND s.kindergarten_id = ?"
		condition = append(condition, id)
	}
	if class_type != 0 {
		where += " AND o.class_type = ?"
		condition = append(condition, class_type)
	}
	o := orm.NewOrm()
	qb, _ := orm.NewQueryBuilder("mysql")

	// 构建查询对象
	sql := qb.Select("count(*)").From("student as s").LeftJoin("organizational_member as om").
		On("s.student_id = om.member_id").LeftJoin("organizational as o").
		On("om.organizational_id = o.id").Where(where).And("isnull(deleted_at)").String()
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
		sql := qb.Select("s.name", "s.avatar", "s.teacher_id", "s.number", "s.phone", "o.name as class").From("student as s").LeftJoin("organizational_member as om").
			On("s.student_id = om.member_id").LeftJoin("organizational as o").
			On("om.organizational_id = o.id").Where(where).And("isnull(deleted_at)").Limit(prepage).Offset(limit).String()
		num, err := o.Raw(sql, condition).Values(&v)
		if err == nil && num > 0 {
			paginatorMap := make(map[string]interface{})
			data := make(map[string][]interface{})
			paginatorMap["total"] = total //总条数
			for _, val := range v {
				if strc, ok := val["class"].(string); ok {
					data[strc] = append(data[strc], val)
				}
			}
			//分页数据
			paginatorMap["data"] = data
			paginatorMap["page_num"] = totalpages //总页数
			return paginatorMap
		}
	}
	return nil
}
