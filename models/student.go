package models

import (
	"encoding/json"
	"math"
	"time"

	"github.com/astaxie/beego/orm"
)

type Student struct {
	Id               int        `json:"student_id" orm:"column(student_id);auto"`
	Name             string     `json:"name" orm:"column(name);size(20)" description:"姓名"`
	Age              int8       `json:"age" orm:"column(age)" description:"年龄"`
	Sex              int8       `json:"sex" orm:"column(sex)" description:"性别 0男 1女"`
	NativePlace      string     `json:"native_place" orm:"column(native_place);size(20)" description:"籍贯"`
	NationOrReligion string     `json:"nation_or_religion" orm:"column(nation_or_religion);size(20)" description:"民族或宗教"`
	Number           string     `json:"number" orm:"column(number);size(11)" description:"学号"`
	ClassInfo        string     `json:"class_info" orm:"column(class_info);size(20)" description:"所在班级"`
	Address          string     `json:"address" orm:"column(address);size(50)" description:"住址"`
	Avatar           string     `json:"avatar" orm:"column(avatar);size(150)" description:"头像"`
	Status           int8       `json:"status" orm:"column(status)" description:"状态 0未分班 1已分班 2离园"`
	UserId           int        `json:"user_id" orm:"column(user_id)" description:"用户ID"`
	KindergartenId   int        `json:"kindergarten_id" orm:"column(kindergarten_id)" description:"幼儿园ID"`
	Phone            string     `json:"phone" orm:"column(phone);size(11)" description:"手机号"`
	HealthStatus     string     `json:"health_status" orm:"column(health_status);size(150)" description:"健康状况，多个以逗号隔开"`
	Hobby            string     `json:"hobby" orm:"column(hobby);size(150)" description:"兴趣爱好，多个以逗号隔开"`
	CreatedAt        time.Time  `json:"Created_at" orm:"auto_now_add"`
	UpdatedAt        time.Time  `json:"updated_at" orm:"auto_now"`
	DeletedAt        time.Time  `json:"deleted_at" orm:"column(deleted_at);type(datetime);null"`
	Kinship          []*Kinship `json:"kinship" orm:"reverse(many)"`
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
		sql := qb.Select("s.name", "s.avatar", "s.student_id", "s.number", "s.phone", "o.name as class").From("student as s").LeftJoin("organizational_member as om").
			On("s.student_id = om.member_id").LeftJoin("organizational as o").
			On("om.organizational_id = o.id").Where(where).And("status = 1").And("isnull(deleted_at)").Limit(prepage).Offset(limit).String()
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

//删除学生
func DeleteStudent(id int, status int, ty int, class_type int) map[string]interface{} {
	o := orm.NewOrm()
	v := Student{Id: id}
	timeLayout := "2006-01-02 15:04:05" //转化所需模板
	loc, _ := time.LoadLocation("")
	timenow := time.Now().Format("2006-01-02 15:04:05")
	if err := o.Read(&v); err == nil {
		if status == 0 && ty == 0 {
			v.Status = 2
		} else if status == 0 && ty == 1 || status == 2 {
			v.DeletedAt, _ = time.ParseInLocation(timeLayout, timenow, loc)
		} else if class_type == 3 && ty == 0 || class_type == 2 && ty == 0 || class_type == 1 && ty == 0 {
			v.Status = 2
		} else if class_type == 3 && ty == 1 || class_type == 2 && ty == 1 || class_type == 1 && ty == 1 {
			v.DeletedAt, _ = time.ParseInLocation(timeLayout, timenow, loc)
		}
		if _, err = o.Update(&v); err == nil {
			paginatorMap := make(map[string]interface{})
			paginatorMap["data"] = nil //返回数据
			return paginatorMap
		}
	}
	return nil
}

//学生详情
func GetStudentInfo(id int) map[string]interface{} {
	o := orm.NewOrm()
	student := Student{Id: id}
	o.Read(&student)
	num, err := o.LoadRelated(&student, "kinship")
	if err == nil && num > 0 {
		paginatorMap := make(map[string]interface{})
		paginatorMap["data"] = student
		return paginatorMap
	}
	return nil
}

//学生--编辑
func UpdateStudent(id int, student string, kinship string) map[string]interface{} {
	o := orm.NewOrm()
	err := o.Begin()
	v := Student{Id: id}
	//编辑学生信息
	var s *Student
	json.Unmarshal([]byte(student), &s)
	if err := o.Read(&v); err == nil {
		s.Id = v.Id
		if _, err = o.Update(s); err == nil {
			if err != nil {
				err = o.Rollback()
			}
		}
	}
	//编辑亲属信息
	var k []Kinship
	json.Unmarshal([]byte(kinship), &k)
	for key, value := range k {
		_, err := o.QueryTable("kinship").Filter("kinship_id", k[key].Id).Update(orm.Params{
			"name": value.Name, "relation": value.Relation, "unit_name": value.UnitName, "contact_information": value.ContactInformation,
		})
		if err == nil {
			err = o.Commit()
		} else {
			err = o.Rollback()
		}
	}
	if err == nil {
		paginatorMap := make(map[string]interface{})
		paginatorMap["data"] = 1 //返回数据
		return paginatorMap
	}
	return nil
}

//学生-录入信息
func AddStudent(student string, kinship string) map[string]interface{} {
	o := orm.NewOrm()
	err := o.Begin()
	//写入学生表
	var s Student
	json.Unmarshal([]byte(student), &s)
	_, err = o.Insert(&s)
	if err != nil {
		err = o.Rollback()
	}
	//写入亲属表
	var k []Kinship
	json.Unmarshal([]byte(kinship), &k)
	id, err := o.InsertMulti(100, &k)
	if err != nil {
		err = o.Rollback()
	} else {
		err = o.Commit()
	}
	if err == nil {
		paginatorMap := make(map[string]interface{})
		paginatorMap["data"] = id //返回数据
		return paginatorMap
	}
	return nil
}
