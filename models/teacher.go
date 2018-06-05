package models

import (
	"fmt"
	"math"
	"time"

	"github.com/astaxie/beego/orm"
)

type Teacher struct {
	Id                         int       `json:"teacher_id" orm:"column(teacher_id);auto" description:"自增id"`
	Name                       string    `json:"name" orm:"column(name);size(20)" description:"姓名"`
	Age                        int8      `json:"age" orm:"column(age)" description:"年龄"`
	Sex                        int8      `json:"sex" orm:"column(sex)" description:"性别 0男  1女"`
	Avatar                     string    `json:"avatar" orm:"column(avatar);size(150)" description:"头像"`
	Number                     string    `json:"number" orm:"column(number);size(20)" description:"教职工编号"`
	NationOrReligion           string    `json:"nation_or_religion" orm:"column(nation_or_religion);size(10)" description:"民族或宗教"`
	NativePlace                string    `json:"native_place" orm:"column(native_place);size(20)" description:"籍贯"`
	UserId                     int       `json:"user_id" orm:"column(user_id)" description:"用户id"`
	ClassInfo                  string    `json:"class_info" orm:"column(class_info);size(10)" description:"班级信息"`
	Phone                      string    `json:"phone" orm:"column(phone);size(11)" description:"联系电话"`
	EnterGardenTime            time.Time `json:"enter_garden_time" orm:"column(enter_garden_time);type(date)" description:"进入本园时间"`
	EnterJobTime               time.Time `json:"enter_job_time" orm:"column(enter_job_time);type(date)" description:"参加工作时间"`
	KindergartenId             int       `json:"kindergarten_id" orm:"column(kindergarten_id)" description:"幼儿园id"`
	Address                    string    `json:"address" orm:"column(address);size(191)" description:"住址"`
	IdNumber                   string    `json:"id_number" orm:"column(id_number);size(18)" description:"身份证号"`
	EmergencyContact           string    `json:"emergency_contact" orm:"column(emergency_contact);size(20)" description:"紧急联系人"`
	EmergencyContactPhone      string    `json:"emergency_contact_phone" orm:"column(emergency_contact_phone);size(11)" description:"紧急联系人电话"`
	Post                       string    `json:"post" orm:"column(post);size(10)" description:"职务"`
	Source                     string    `json:"source" orm:"column(source);size(191)" description:"来源"`
	TeacherCertificationNumber string    `json:"teacher_certification_number" orm:"column(teacher_certification_number);size(20)" description:"教师资格认证编号"`
	TeacherCertificationStatus int8      `json:"teacher_certification_status" orm:"column(teacher_certification_status)" description:"教师资格证书状态，是否认证"`
	Status                     int8      `json:"status" orm:"column(status)" description:"状态：0未分班，1已分班，2离职"`
	CreatedAt                  time.Time `json:"created_at" orm:"auto_now_add"`
	UpdatedAt                  time.Time `json:"updated_at" orm:"auto_now"`
	DeletedAt                  time.Time `json:"deleted_at" orm:""`
}

func (t *Teacher) TableName() string {
	return "teacher"
}

func init() {
	orm.RegisterModel(new(Teacher))
}

//教师下拉菜单
func GetTeacherById(id int, page, prepage int) map[string]interface{} {
	var v []Teacher
	o := orm.NewOrm()
	qb, _ := orm.NewQueryBuilder("mysql")
	// 构建查询对象
	qb.Select("teacher.*").From("teacher").LeftJoin("teachers_show").
		On("teacher.teacher_id = teachers_show.teacher_id").Where("teacher.kindergarten_id = ?").
		And("teacher.status != ?").And("isnull(teacher.deleted_at)").And("isnull(teachers_show.id)")
	sql := qb.String()
	_, err := o.Raw(sql, id, 2).QueryRows(&v)
	if err == nil {
		paginatorMap := make(map[string]interface{})
		paginatorMap["data"] = v //返回数据
		return paginatorMap
	}
	return nil
}

//教师列表
func GetTeacher(id int, status int, search string, page int, prepage int) map[string]interface{} {
	var condition []interface{}
	where := "1=1 "
	if id == 0 {
		where += " AND t.kindergarten_id = ?"
	} else {
		where += " AND t.kindergarten_id = ?"
		condition = append(condition, id)
	}
	if status != -1 {
		where += " AND t.status = ?"
		condition = append(condition, status)
	}
	if search != "" {
		where += " AND t.name like ?"
		condition = append(condition, "%"+search+"%")
	}
	o := orm.NewOrm()
	qb, _ := orm.NewQueryBuilder("mysql")

	// 构建查询对象
	sql := qb.Select("count(*)").From("teacher as t").LeftJoin("organizational_member as om").
		On("t.teacher_id = om.member_id").LeftJoin("organizational as o").
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
		sql := qb.Select("t.name", "t.avatar", "t.teacher_id", "t.number", "t.phone", "o.name as class").From("teacher as t").LeftJoin("organizational_member as om").
			On("t.teacher_id = om.member_id").LeftJoin("organizational as o").
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

//班级列表
func GetClass(id int, class_type int, page int, prepage int) map[string]interface{} {
	var condition []interface{}
	where := "1=1 "
	if id == 0 {
		where += " AND t.kindergarten_id = ?"
	} else {
		where += " AND t.kindergarten_id = ?"
		condition = append(condition, id)
	}
	if class_type != 0 {
		where += " AND o.class_type = ?"
		condition = append(condition, class_type)
	}
	o := orm.NewOrm()
	qb, _ := orm.NewQueryBuilder("mysql")

	// 构建查询对象
	sql := qb.Select("count(*)").From("teacher as t").LeftJoin("organizational_member as om").
		On("t.teacher_id = om.member_id").LeftJoin("organizational as o").
		On("om.organizational_id = o.id").Where(where).And("status = 1").And("isnull(deleted_at)").String()
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
		sql := qb.Select("t.name", "t.avatar", "t.teacher_id", "t.number", "t.phone", "o.name as class").From("teacher as t").LeftJoin("organizational_member as om").
			On("t.teacher_id = om.member_id").LeftJoin("organizational as o").
			On("om.organizational_id = o.id").Where(where).And("isnull(deleted_at)").And("status = 1").Limit(prepage).Offset(limit).String()
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

//删除教师
func DeleteTeacher(id int, status int, class_type int) map[string]interface{} {
	o := orm.NewOrm()
	v := Teacher{Id: id}
	timeLayout := "2006-01-02 15:04:05" //转化所需模板
	loc, _ := time.LoadLocation("")
	timenow := time.Now().Format("2006-01-02 15:04:05")
	fmt.Println(time.ParseInLocation(timeLayout, timenow, loc))
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
			//			_, err = o.QueryTable("teachers_show").Filter("teacher_id", id).Delete()
			if err == nil {
				paginatorMap := make(map[string]interface{})
				paginatorMap["data"] = nil //返回数据
				return paginatorMap
			}
		}
	}
	return nil
}

//教师详情
func GetTeacherInfo(id int) map[string]interface{} {
	o := orm.NewOrm()
	v := &Teacher{Id: id}
	if err := o.Read(v); err == nil {
		paginatorMap := make(map[string]interface{})
		paginatorMap["data"] = v
		return paginatorMap
	}

	return nil
}

//教师--编辑
func UpdateTeacher(m *Teacher) map[string]interface{} {
	o := orm.NewOrm()
	v := Teacher{Id: m.Id}
	if m.Post == "" {
		m.Post = "普通教师"
	}
	if err := o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			paginatorMap := make(map[string]interface{})
			paginatorMap["data"] = num //返回数据
			return paginatorMap
		}
	}
	return nil
}

//教师-录入信息
//Admin - 动画--添加
func AddTeacher(m *Teacher) map[string]interface{} {
	o := orm.NewOrm()
	if m.Post == "" {
		m.Post = "普通教师"
	}
	id, err := o.Insert(m)
	if err == nil {
		paginatorMap := make(map[string]interface{})
		paginatorMap["data"] = id //返回数据
		return paginatorMap
	}
	return nil
}

//组织框架移除教师
func RemoveTeacher(teacher_id int, class_id int) map[string]interface{} {
	o := orm.NewOrm()
	err := o.Begin()
	_, err = o.QueryTable("organizational_member").Filter("organizational_id", class_id).Filter("member_id", teacher_id).Delete()
	if err != nil {
		err = o.Rollback()
	}
	num, err := o.QueryTable("teacher").Filter("teacher_id", teacher_id).Update(orm.Params{
		"status": 0,
	})
	if err != nil {
		err = o.Rollback()
	}
	_, err = o.QueryTable("teachers_show").Filter("teacher_id", teacher_id).Delete()
	if err != nil {
		err = o.Rollback()
	} else {
		err = o.Commit()
	}
	if err == nil && num > 0 {
		paginatorMap := make(map[string]interface{})
		paginatorMap["data"] = num //返回数据
		return paginatorMap
	}
	return nil
}
