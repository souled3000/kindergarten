package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/astaxie/beego/orm"
)

type Student struct {
	Id               int       `json:"student_id" orm:"column(student_id);auto"`
	Name             string    `json:"name" orm:"column(name);size(20)" description:"姓名"`
	Age              int8      `json:"age" orm:"column(age)" description:"年龄"`
	Sex              int8      `json:"sex" orm:"column(sex)" description:"性别 0男 1女"`
	NativePlace      string    `json:"native_place" orm:"column(native_place);size(20)" description:"籍贯"`
	NationOrReligion string    `json:"nation_or_religion" orm:"column(nation_or_religion);size(20)" description:"民族或宗教"`
	Number           string    `json:"number" orm:"column(number);size(11)" description:"学号"`
	ClassInfo        string    `json:"class_info" orm:"column(class_info);size(20)" description:"所在班级"`
	Address          string    `json:"address" orm:"column(address);size(50)" description:"住址"`
	Avatar           string    `json:"avatar" orm:"column(avatar);size(150)" description:"头像"`
	Status           int8      `json:"status" orm:"column(status)" description:"状态 0未分班 1已分班 2离园"`
	BabyId           int       `json:"baby_id" orm:"column(baby_id)" description:"宝宝ID"`
	Birthday         string    `json:"birthday" orm:"column(birthday)" description:"生日"`
	KindergartenId   int       `json:"kindergarten_id" orm:"column(kindergarten_id)" description:"幼儿园ID"`
	Phone            string    `json:"phone" orm:"column(phone);size(11)" description:"手机号"`
	HealthStatus     string    `json:"health_status" orm:"column(health_status);size(150)" description:"健康状况，多个以逗号隔开"`
	Hobby            string    `json:"hobby" orm:"column(hobby);size(150)" description:"兴趣爱好，多个以逗号隔开"`
	IsMuslim         int       `json:"is_muslim" orm:"column(is_muslim)"`
	BabyName         string    `json:"baby_name" orm:"column(baby_name)"`
	CreatedAt        time.Time `json:"Created_at" orm:"auto_now_add"`
	UpdatedAt        time.Time `json:"updated_at" orm:"auto_now"`
	DeletedAt        time.Time `json:"deleted_at" orm:"column(deleted_at);type(datetime);null"`
}

type inviteStudent struct {
	Name           string `json:"name"`
	Phone          string `json:"phone"`
	BabyId         int    `json:"baby_id"`
	Birthday       string `json:"birthday"`
	KindergartenId int    `json:"kindergarten_id"`
}

func (t *Student) TableName() string {
	return "student"
}

func init() {
	orm.RegisterModel(new(Student))
}

/*
学生列表
*/
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
	sql := qb.Select("count(*)").From("student as s").Where(where).And("isnull(deleted_at)").String()
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
		sql := qb.Select("s.name", "s.avatar", "s.student_id", "s.number", "s.phone").
			From("student as s").Where(where).And("isnull(deleted_at)").Limit(prepage).Offset(limit).String()
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
学生班级列表
*/
func GetStudentClass(id int, class_type int, page int, prepage int) map[string]interface{} {
	var condition []interface{}
	where := "1=1 "
	if id == 0 {
		where += " AND o.kindergarten_id = ?"
	} else {
		where += " AND o.kindergarten_id = ?"
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

/*
删除学生
*/
func DeleteStudent(id int, status int, ty int, class_type int) error {
	o := orm.NewOrm()
	o.Begin()
	v := Student{Id: id}
	timeLayout := "2006-01-02 15:04:05" //转化所需模板
	loc, _ := time.LoadLocation("")
	timenow := time.Now().Format("2006-01-02 15:04:05")
	if err := o.Read(&v); err == nil {
		if status == 0 {
			v.Status = 2
		} else if status == 2 {
			v.DeletedAt, _ = time.ParseInLocation(timeLayout, timenow, loc)
		}

		if class_type == 3 || class_type == 2 || class_type == 1 {
			v.Status = 0
		}
		if _, err = o.Update(&v); err == nil {
			_, err = o.QueryTable("organizational_member").Filter("member_id", id).Delete()
			if err != nil {
				o.Rollback()
				return err
			}
			_, err = o.QueryTable("exceptional_child").Filter("student_id", id).Delete()
			if err != nil {
				o.Rollback()
				return err
			}
		}
	}
	return nil
}

/*
学生详情
*/
func GetStudentInfo(id int) (paginatorMap map[string]interface{}, err error) {
	o := orm.NewOrm()
	var kinships []Kinship
	//学生信息
	student := Student{Id: id}
	err = o.Read(&student)
	//亲属信息
	qb, _ := orm.NewQueryBuilder("mysql")
	sql := qb.Select("*").From("kinship").Where("student_id = ?").String()
	num, err := o.Raw(sql, id).QueryRows(&kinships)
	if err == nil && num > 0 {
		paginatorMap := make(map[string]interface{})
		paginatorMap["student"] = student
		paginatorMap["kinship"] = kinships
		return paginatorMap, nil
	}
	err = errors.New("获取失败")
	return nil, err
}

/*
学生编辑
*/
func UpdateStudent(id int, student string, kinship string) (paginatorMap map[string]interface{}, err error) {
	o := orm.NewOrm()
	err = o.Begin()
	v := Student{Id: id}
	//编辑学生信息
	var s *Student
	json.Unmarshal([]byte(student), &s)
	if err := o.Read(&v); err == nil {
		s.Id = v.Id
		if _, err = o.Update(s); err == nil {
			//写入亲属表
			_, err = o.QueryTable("kinship").Filter("student_id", id).Delete()
			if err == nil {
				var k []Kinship
				json.Unmarshal([]byte(kinship), &k)
				_, err = o.InsertMulti(100, &k)
				if err != nil {
					o.Rollback()
					return nil, nil
				} else {
					o.Commit()
					return nil, err
				}
			}
		} else {
			o.Rollback()
		}
	}

	if err == nil {
		return nil, nil
	}
	err = errors.New("编辑失败")
	return nil, err
}

/*
学生-录入信息
*/
func AddStudent(student string, kinship string) (paginatorMap map[string]interface{}, err error) {
	o := orm.NewOrm()
	o.Begin()
	paginatorMap = make(map[string]interface{})
	//写入学生表
	var s Student
	json.Unmarshal([]byte(student), &s)
	id, err := o.Insert(&s)
	fmt.Println(id)
	if err != nil {
		o.Rollback()
		return nil, err
	}
	if kinship != "" {
		ids := strconv.FormatInt(id, 10)
		kid, _ := strconv.Atoi(ids)
		//写入亲属表
		var k []Kinship
		json.Unmarshal([]byte(kinship), &k)
		for key, _ := range k {
			k[key].StudentId = kid
		}
		id, err = o.InsertMulti(100, &k)
		if err != nil {
			o.Rollback()
			return nil, err
		}
	}
	_, err = o.QueryTable("baby_kindergarten").Filter("baby_id", s.BabyId).Update(orm.Params{
		"actived": 0,
	})
	if err != nil {
		o.Rollback()
		err = errors.New("保存失败")
		return nil, err
	} else {
		o.Commit()
		return nil, nil
	}
}

/*
移除学生
*/
func RemoveStudent(class_id int, student_id int) error {
	o := orm.NewOrm()
	o.Begin()
	_, err := o.QueryTable("organizational_member").Filter("organizational_id", class_id).Filter("member_id", student_id).Delete()
	if err != nil {
		o.Rollback()
		return err
	}
	_, err = o.QueryTable("exceptional_child").Filter("student_id", student_id).Delete()
	if err != nil {
		o.Rollback()
		return err
	}
	_, err = o.QueryTable("student").Filter("student_id", student_id).Update(orm.Params{
		"status": 0,
	})
	if err != nil {
		o.Rollback()
		return err
	} else {
		o.Commit()
		return nil
	}
	return err
}

/*
邀请学生
*/
func Invite(student string) error {
	o := orm.NewOrm()
	var someError error
	var baby []orm.Params
	var s []inviteStudent
	json.Unmarshal([]byte(student), &s)
	timeLayout := "2006-01-02 15:04:05" //转化所需模板
	loc, _ := time.LoadLocation("")
	timenow := time.Now().Format("2006-01-02 15:04:05")
	createTime, _ := time.ParseInLocation(timeLayout, timenow, loc)
	for _, v := range s {
		t, _ := time.Parse("2006-01-02 15:04:05", v.Birthday)
		qb, _ := orm.NewQueryBuilder("mysql")
		sql := qb.Select("*").From("baby_kindergarten").Where("baby_id = ?").And("status = 0").String()
		_, err := o.Raw(sql, v.BabyId).Values(&baby)
		if err == nil {
			if baby != nil {
				someError = errors.New("" + string(v.Name) + "已被邀请过")
				continue
			} else {
				sql = "insert into baby_kindergarten set kindergarten_id = ?,baby_id = ?,baby_name = ?,created_at = ?,birthday = ?,phone = ?"
				_, err := o.Raw(sql, v.KindergartenId, v.BabyId, v.Name, createTime, t, v.Phone).Exec()
				if err != nil {
					err = errors.New("邀请失败")
					return err
				}
			}
		}
	}
	return someError
}

/*
未激活baby
*/
func GetBabyInfo(kindergarten_id int) (paginatorMap map[string]interface{}, err error) {
	o := orm.NewOrm()
	var baby []orm.Params
	qb, _ := orm.NewQueryBuilder("mysql")
	sql := qb.Select("baby_id", "baby_name", "kindergarten_id", "birthday", "phone").From("baby_kindergarten").Where("kindergarten_id = ?").And("actived = 1").And("status = 0").String()
	_, err = o.Raw(sql, kindergarten_id).Values(&baby)
	if err == nil {
		paginatorMap := make(map[string]interface{})
		paginatorMap["data"] = baby
		return paginatorMap, nil
	}
	err = errors.New("获取失败")
	return nil, err
}

/*
学生名字获取班级
*/
func GetNameClass(name string) (paginatorMap map[string]interface{}, err error) {
	o := orm.NewOrm()
	var class []orm.Params
	qb, _ := orm.NewQueryBuilder("mysql")
	sql := qb.Select("s.student_id", "o.id as class_id", "o.name as class_name", "class_type").From("student as s").LeftJoin("organizational_member as om").
		On("s.student_id = om.member_id").LeftJoin("organizational as o").
		On("om.organizational_id = o.id").Where("s.name = ?").And("om.type = 1").String()
	_, err = o.Raw(sql, name).Values(&class)
	for key, val := range class {
		if val["class_type"].(string) == "3" {
			class[key]["class"] = "大班" + val["class_name"].(string) + ""
		} else if val["class_type"].(string) == "2" {
			class[key]["class"] = "中班" + val["class_name"].(string) + ""
		} else {
			class[key]["class"] = "小班" + val["class_name"].(string) + ""
		}
	}
	if err == nil {
		paginatorMap := make(map[string]interface{})
		paginatorMap["data"] = class
		return paginatorMap, nil
	}
	err = errors.New("获取失败")
	return nil, err
}
