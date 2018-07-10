package models

import (
	"errors"
	"fmt"
	//	"math"
	//	"strconv"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"time"
)

type Attendance struct {
	Id      int       `json:"id" orm:"column(id);auto"`
	Sid     int       `json:"sid" orm:"column(sid)"`
	Name    string    `json:"name" orm:"column(name)" description:"学生姓名"`
	Today   string    `json:"name" orm:"column(today)" description:"考勤日期"`
	Cls     string    `json:"cls" orm:"column(cls)" description:"班级名称"`
	AttTime time.Time `json:"att_time" orm:"auto_now"`
	Status  int       `json:"status" orm:"column(status)" description:"状态"`
}

func (t *Attendance) TableName() string {
	return "attendance"
}
func init() {
	orm.RegisterModel(new(Attendance))
}

type Leave struct {
	Id        int       `json:"id" orm:"column(id);auto"`
	Sid       int       `json:"sid" orm:"column(sid)" description:"学生ID"`
	Applicant string    `json:"applicant" orm:"column(applicant)" description:"申请人"`
	Reason    string    `json:"reason" orm:"column(reason)" description:"理由"`
	Beg       time.Time `json:"beg" orm:"column(beg);type(datetime);null"`
	End       time.Time `json:"end" orm:"column(end);type(datetime);null"`
	Status    int       `json:"status" orm:"column(status)" description:"状态"`
	Type      int       `json:"Type" orm:"column(type)" description:"0:事假;1:病假"`
	AppTime   time.Time `json:"app" orm:"auto_now_add"`
}

func (t *Leave) TableName() string {
	return "aleave"
}
func init() {
	orm.RegisterModel(new(Leave))
}

/*
考勤
*/
func Att(sid, status int) (err error) {
	today := time.Now().Format("2006-01-02")
	fmt.Println(sid, status)
	db := orm.NewOrm()
	db.Begin()
	defer func() {
		if err != nil {
			db.Rollback()
			err = errors.New("保存失败")
		} else {
			db.Commit()
		}
	}()

	var a Attendance
	a.Sid = sid
	a.Today = today

	created, id, err := db.ReadOrCreate(&a, "Sid", "Today")
	if created {
		var st Student
		db.QueryTable("student").Filter("student_id", sid).One(&st, "name", "class_info")
		err = db.Read(&st, "name", "class_info")
		if err != nil {
			beego.Debug("read stu:", err)
			return
		}
		beego.Debug("stu:", st.Id, st.Name, st.ClassInfo)
		a.Name = st.Name
		a.Cls = st.ClassInfo
	}
	a.Status = status
	db.Update(&a)

	beego.Debug(id, a.Name, a.Status, a.AttTime, err)
	return
}

/*
请假
*/
func AskForLeave(o Leave) (err error) {
	fmt.Println(o)
	db := orm.NewOrm()
	db.Begin()
	defer func() {
		if err != nil {
			db.Rollback()
			err = errors.New("保存失败")
		} else {
			db.Commit()
		}
	}()
	id,err:=db.Insert(&o)
	beego.Debug(id,err)
	return
}


/*
获取教师下的学生 为教师页面
*/
func GotStdsByTeaID(tid int) (rt []orm.Params) {
	sql := "select distinct s.name,s.student_id ,s.class_info from organizational_member om,organizational_member o2 ,student s  where om.type=0 and om.member_id= ? and s.student_id = o2.member_id and o2.type=1 and s.status=1"
	orm.NewOrm().Raw(sql, tid).Values(&rt)
	return
}

/*
获得某日考勤记录
*/
func Query(date time.Time, orgid int) (rt []Attendance) {
	db := orm.NewOrm()
	var condition []interface{}
	condition = append(condition, orgid)
	sql := "select a.* from attendance a,organizatinal_member om  om.kindergarten_id = ? and a.organizational_id = om.organizational_id"
	db.Raw(sql, condition).QueryRows(&rt)
	return
}
