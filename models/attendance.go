package models

import (
	"errors"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"time"
	"strings"
)

type Attendance struct {
	Id        int       `json:"id" orm:"column(id);auto"`
	Sid       int       `json:"sid" orm:"column(sid)"`
	Name      string    `json:"name" orm:"column(name)" description:"学生姓名"`
	Today     string    `json:"name" orm:"column(today)" description:"考勤日期"`
	Cls       string    `json:"cls" orm:"column(cls)" description:"班级名称"`
	Morning   time.Time `json:"morning" orm:"column(morning)"`
	Afternoon time.Time `json:"afternoon" orm:"column(afternoon)"`
	Status    int       `json:"status" orm:"column(status)" description:"状态"`
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

type AttendanceRule struct {
	Id   int    `json:"id" orm:"column(id);auto"`
	Kid  int    `json:"kid" orm:"column(kid)" description:"学校ID"`
	Mbeg string `json:"mbeg" orm:"column(m_beg)" description:"上午起始时间"`
	Mend string `json:"mend" orm:"column(m_end)" description:"上午结束时间"`
	Abeg string `json:"abeg" orm:"column(a_beg)" description:"下午起始时间"`
	Aend string `json:"aend" orm:"column(a_end)" description:"下午结束时间"`
	Days string `json:"days" orm:"column(days)" description:"工作日"`
}

func (this *AttendanceRule) TableName() string {
	return "attendance_rule"
}
func init() {
	orm.RegisterModel(new(AttendanceRule))
}

/*
考勤上下限
*/
func AttRule(o AttendanceRule) (err error) {
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
	created, id, err := db.ReadOrCreate(&o, "Kid")
	if !created {
		db.Update(&o)
	}
	beego.Debug(id, err)
	return
}

/*
考勤
*/
func Att(sid, status int) (err error) {
	now := time.Now()
	today := now.Format("2006-01-02")
	noon, _ := time.ParseInLocation("2006-01-02 15:04:05", today+" 12:00:00", time.Local)
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
		err = db.QueryTable("student").Filter("student_id", sid).One(&st, "name", "class_info")
		if err != nil {
			beego.Debug("read stu:", err)
			return
		}
		beego.Debug("stu:", st.Id, st.Name, st.ClassInfo)
		a.Name = st.Name
		a.Cls = st.ClassInfo
	}

	if now.Unix() <= noon.Unix() {
		a.Morning = now
	} else {
		a.Afternoon = now
	}

	a.Status = status
	db.Update(&a)

	beego.Debug(id, a.Name, a.Status, err)
	return
}

/*
一键入园
*/
func ToAll(cid int) (err error) {
	now := time.Now()
	today := now.Format("2006-01-02")
	noon, _ := time.ParseInLocation("2006-01-02 15:04:05", today+" 12:00:00", time.Local)
	db := orm.NewOrm()
	db.Begin()
	defer func() {
		if err != nil {
			db.Rollback()
		} else {
			db.Commit()
		}
	}()
	sql := "select t.member_id from organizational_member t where t.organizational_id = ? and not exists( select t2.sid from aleave t2 where t2.beg <? and t2.end> ? and t.member_id = t2.sid) and not exists ( select sid from attendance t3 where t.member_id = t3.sid and t3.today=? and (t3.morning is not null and t3.afternoon is not null))"
	var ids []int
	db.Raw(sql, cid, now, now, today).QueryRows(&ids)
	var atts []Attendance
	for _, sid := range ids {
		var a Attendance
		a.Sid = sid
		a.Today = today

		var st Student
		st.Id = sid
		e := db.QueryTable("student").Filter("student_id", sid).One(&st, "name", "class_info")
		if e != nil {
			beego.Debug("read stu:", sid, e)
			continue
		}
		beego.Debug("stu:", st.Id, st.Name, st.ClassInfo)
		a.Name = st.Name
		a.Cls = st.ClassInfo

		if now.Unix() <= noon.Unix() {
			a.Morning = now
		} else {
			a.Afternoon = now
		}
		atts = append(atts, a)
	}
	if len(atts) > 0 {
		_, err = db.InsertMulti(len(atts), atts)
	} else {
		err = fmt.Errorf("已全都考勤")
	}
	beego.Debug("TOALL:", len(atts), err)
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
	id, err := db.Insert(&o)
	beego.Debug(id, err)
	return
}
func getWeekNum(weekname string) (n string) {
	switch weekname {
	case "Monday":
		n = "1"
	case "Tuesday":
		n = "2"
	case "Wednesday":
		n = "3"
	case "Thursday":
		n = "4"
	case "Friday":
		n = "5"
	case "Saturday":
		n = "6"
	case "Sunday":
		n = "7"
	}
	return n
}

/*
*	获取教师下的学生 为教师页面
 */
func GotStds(cid int) (rt []orm.Params) {
	db := orm.NewOrm()
	now := time.Now()

	today := now.Format("2006-01-02")
	var t Organizational
	t.Id = cid
	db.Read(&t, "Id")
	var r AttendanceRule
	r.Kid = t.KindergartenId
	db.Read(&r, "Kid")

	if !strings.Contains(r.Days, getWeekNum(now.Weekday().String())){
		return
	}

	//	mbeg, _ := time.ParseInLocation("2006-01-02 15:04:05", today+" "+r.Mbeg,time.Local)
	mend, _ := time.ParseInLocation("2006-01-02 15:04", today+" "+r.Mend, time.Local)
	abeg, _ := time.ParseInLocation("2006-01-02 15:04", today+" "+r.Abeg, time.Local)
	aend, _ := time.ParseInLocation("2006-01-02 15:04", today+" "+r.Aend, time.Local)

	beego.Debug(today+" "+r.Mend, today+" "+r.Abeg, today+" "+r.Aend)
	beego.Debug(mend, abeg, aend)

	var sql string
	switch {
	case (now.Unix() <= mend.Unix()) || (abeg.Unix() <= now.Unix() && now.Unix() <= aend.Unix()):
		//		sql = "select name,sid,cls,morning,afternoon from (select name,w.sid,cls,morning,afternoon from (select t.name,t.sid,t.cls,a.morning,a.afternoon from (select distinct s.name,s.student_id sid ,s.class_info cls from organizational_member om,organizational_member o2 ,student s where om.type=0 and om.member_id=? and s.student_id = o2.member_id and o2.type=1 and s.status=1 and om.organizational_id=o2.organizational_id) t left join attendance a on t.sid=a.sid and a.today=?) w left join aleave l on w.sid = l.sid and l.beg < ? and l.end > ? where l.sid is null) z"
		sql = "select name,sid,cls,morning,afternoon,z.avatar from (select name,w.sid,cls,date_format(morning,'%H:%i') morning ,date_format(afternoon,'%H:%i') afternoon,w.avatar from (select t.name,t.sid,t.cls,a.morning,a.afternoon,t.avatar from (select distinct s.name,s.student_id sid ,s.class_info cls, s.avatar from organizational_member o2 ,student s where s.student_id = o2.member_id and o2.organizational_id = ? and s.status=1 ) t left join attendance a on t.sid=a.sid and a.today=?) w left join aleave l on w.sid = l.sid and l.beg < ? and l.end > ? where l.sid is null) z"
		db.Raw(sql, cid, today, now, now).Values(&rt)
	case (mend.Unix() < now.Unix() && now.Unix() < abeg.Unix()) || now.Unix() > aend.Unix():
		//		sql = "select t.name,t.sid,t.cls,a.morning,a.afternoon from (select distinct s.name,s.student_id sid ,s.class_info cls from organizational_member om,organizational_member o2 ,student s where om.type=0 and om.member_id=? and s.student_id = o2.member_id and o2.type=1 and s.status=1 and om.organizational_id=o2.organizational_id) t inner join attendance a on t.sid=a.sid and a.today=?"
		sql = "select t.name,t.sid,t.cls,date_format(a.morning,'%H:%i') morning,date_format(a.afternoon,'%H:%i') afternoon,t.avatar from (select distinct s.name,s.student_id sid ,s.class_info cls,s.avatar from organizational_member o2 ,student s where s.student_id = o2.member_id and o2.organizational_id=? and s.status=1 ) t inner join attendance a on t.sid=a.sid and a.today=?"
		db.Raw(sql, cid, today).Values(&rt)
	}
	return
}

/*
*	非正常考勤
 */
func GotAbnDtl(tid int) (rt []orm.Params) {
	db := orm.NewOrm()
	now := time.Now()
	today := now.Format("2006-01-02")
	var t Organizational
	t.Id = tid
	db.Read(&t, "Id")
	var r AttendanceRule
	r.Kid = t.KindergartenId
	db.Read(&r, "Kid")
	
	if !strings.Contains(r.Days, getWeekNum(now.Weekday().String())){
		return
	}

	//	mbeg, _ := time.ParseInLocation("2006-01-02 15:04:05", today+" "+r.Mbeg)
	mend, _ := time.ParseInLocation("2006-01-02 15:04", today+" "+r.Mend, time.Local)
	abeg, _ := time.ParseInLocation("2006-01-02 15:04", today+" "+r.Abeg, time.Local)
	aend, _ := time.ParseInLocation("2006-01-02 15:04", today+" "+r.Aend, time.Local)

	var sql string
	switch {
	case (now.Unix() <= mend.Unix()) || (abeg.Unix() <= now.Unix() && now.Unix() <= aend.Unix()):
		//		sql := "select t.name,t.sid,t.cls,a.reason from (select distinct s.name,s.student_id sid ,s.class_info cls from organizational_member om,organizational_member o2 ,student s where om.type=0 and om.member_id= ? and s.student_id = o2.member_id and o2.type=1 and s.status=1 and om.organizational_id=o2.organizational_id) t left join aleave a on t.sid=a.sid and a.beg <? and a.end>?"
		sql := "select t.name,t.sid,t.cls,a.reason,a.avatar from (select distinct s.name,s.student_id sid ,s.avatar,s.class_info cls from organizational_member o2 ,student s where s.student_id = o2.member_id and o2.organizational_id=? and s.status=1 ) t left join aleave a on t.sid=a.sid and a.beg <? and a.end>?"
		orm.NewOrm().Raw(sql, tid, now, now).Values(&rt)
	case (mend.Unix() < now.Unix() && now.Unix() < abeg.Unix()):
		//		sql = "select t.name,t.sid,t.cls,a.morning,a.afternoon ,'' reason,-1 type from (select distinct s.name,s.student_id sid ,s.class_info cls from organizational_member om,organizational_member o2 ,student s where om.type=0 and om.member_id= ? and s.student_id = o2.member_id and o2.type=1 and s.status=1 and om.organizational_id=o2.organizational_id) t  left join attendance a on t.sid=a.sid and (a.morning is null || a.morning >?)" +
		//			"union all" +
		//			" select t.name,t.sid,t.cls,null morning,null afternoon,a.reason,type from (select distinct s.name,s.student_id sid ,s.class_info cls from organizational_member om,organizational_member o2 ,student s where om.type=0 and om.member_id= ? and s.student_id = o2.member_id and o2.type=1 and s.status=1 and om.organizational_id=o2.organizational_id) t join aleave a on t.sid=a.sid and a.beg< ? and a.end >?"
		sql = "select t.name,t.sid,t.cls,a.morning,a.afternoon ,'' reason,-1 type,t.avatar from (select distinct s.avatar,s.name,s.student_id sid ,s.class_info cls from organizational_member o2 ,student s where s.student_id = o2.member_id and o2.organizational_id=? and s.status=1 and not exists (select l.sid from aleave l where l.sid=s.student_id and l.beg<? and l.end>?)) t  left join attendance a on t.sid=a.sid and (a.morning is null || a.morning >?)" +
			"union all" +
			" select t.name,t.sid,t.cls,null morning,null afternoon,a.reason,type,t.avatar from (select distinct s.avatar,s.name,s.student_id sid ,s.class_info cls from organizational_member o2 ,student s where s.student_id = o2.member_id and o2.organizational_id=? and s.status=1 ) t join aleave a on t.sid=a.sid and a.beg< ? and a.end >?"
		orm.NewOrm().Raw(sql, tid, now, now, now, tid, now, mend).Values(&rt)
	case now.Unix() > aend.Unix():
		//afternoon<abeg 早退 ；afternoon >aend 补勤
		//		sql = "select t.name,t.sid,t.cls,a.morning,a.afternoon ,'' reason,-1 type from (select distinct s.name,s.student_id sid ,s.class_info cls from organizational_member om,organizational_member o2 ,student s where om.type=0 and om.member_id=? and s.student_id = o2.member_id and o2.type=1 and s.status=1 and om.organizational_id=o2.organizational_id) t  left join attendance a on t.sid=a.sid and (a.afternoon is null || a.afternoon <? || a.afternoon >?)" +
		//			"union all" +
		//			" select t.name,t.sid,t.cls,null morning,null afternoon,a.reason,type from (select distinct s.name,s.student_id sid ,s.class_info cls from organizational_member om,organizational_member o2 ,student s where om.type=0 and om.member_id= ? and s.student_id = o2.member_id and o2.type=1 and s.status=1 and om.organizational_id=o2.organizational_id) t join aleave a on t.sid=a.sid and a.beg< ? and a.end >?"
		sql = "select t.name,t.sid,t.cls,a.morning,a.afternoon ,'' reason,-1 type,t.avatar from (select distinct s.avatar,s.name,s.student_id sid ,s.class_info cls from organizational_member o2 ,student s where s.student_id = o2.member_id and o2.organizational_id=? and s.status=1 and not exists (select l.sid from aleave l where l.sid=s.student_id and l.beg<? and l.end>?)) t  left join attendance a on t.sid=a.sid and (a.afternoon is null || a.afternoon <? || a.afternoon >?)" +
			"union all" +
			" select t.name,t.sid,t.cls,null morning,null afternoon,a.reason,type,t.avatar from (select distinct s.avatar,s.name,s.student_id sid ,s.class_info cls from organizational_member o2 ,student s where s.student_id = o2.member_id and o2.organizational_id=? and s.status=1 ) t join aleave a on t.sid=a.sid and a.beg< ? and a.end >?"
		orm.NewOrm().Raw(sql, tid, now, now, abeg, aend, tid, now, now).Values(&rt)
	}
	return
}

/*
获得某日某年级考勤统计
*/
func CountByGrade(day string, grade int) (rt []orm.Params) {
	db := orm.NewOrm()
	beego.Info(day)
	aday, _ := time.ParseInLocation("2006-01-02", day, time.Local)
	var t Organizational
	t.Id = grade
	db.Read(&t, "Id")
	var r AttendanceRule
	r.Kid = t.KindergartenId
	db.Read(&r, "Kid")

	beego.Debug("KID:", r.Kid)

	//	mbeg, _ := time.ParseInLocation("2006-01-02 15:04", today+" "+r.Mbeg)
	mend, _ := time.ParseInLocation("2006-01-02 15:04", day+" "+r.Mend, time.Local)
	abeg, _ := time.ParseInLocation("2006-01-02 15:04", day+" "+r.Abeg, time.Local)
	aend, _ := time.ParseInLocation("2006-01-02 15:04", day+" "+r.Aend, time.Local)
	var sql string
	sql += "select name,id,max(case type when -1 then amount else 0 end) 'good', max(case type when 0 then amount else 0 end) 'casual', max(case type when 1 then amount else 0 end) 'sick' from ("
	sql += " select t1.name,t1.id,-1 type,count(*) amount from organizational t1 , organizational_member t2 , attendance t3 where t1.parent_id= ? and t1.id = t2.organizational_id and t2.member_id = t3.sid and t3.today = ? and t3.morning < ? and t3.afternoon between ? and ?  group by t1.id"
	sql += " union all"
	sql += " select t4.name,t4.id,t6.type,count(*) amount from organizational t4,organizational_member t5 , aleave t6 where t4.parent_id = ? and t4.id=t5.organizational_id and t5.member_id =t6.sid and t6.beg < ? and t6.end > ? group by t4.id"
	sql += " ) z group by name"
	db.Raw(sql, grade, day, mend, abeg, aend, grade, aday, aday).Values(&rt)
	//	sql += "select name,max(case type when -1 then amount else 0 end) 'good', max(case type when 0 then amount else 0 end) 'casual', max(case type when 1 then amount else 0 end) 'sick' from ("
	//	sql += " select t1.name,-1 type,count(*) amount from organizational t1 , organizational_member t2 , attendance t3 where t1.parent_id=? and t1.id = t2.organizational_id and t2.member_id = t3.sid and t3.today = ? group by t1.id"
	//	sql += " union all"
	//	sql += " select t4.name,t6.type,count(*) amount from organizational t4,organizational_member t5 , aleave t6 where t4.parent_id = ? and t4.id=t5.organizational_id and t5.member_id =t6.sid and t6.beg < ? and t6.end > ? group by t4.id "
	//	sql += " ) z group by name"
	//	db.Raw(sql, grade, day, grade, aday, aday).Values(&rt)

	beego.Debug(rt)
	return
}

/*
获取某日某班考勤记录
*/
func GotAttsByDayAndCls(day string, cid int) (rt map[string]interface{}) {
	db := orm.NewOrm()
	beego.Info(day)
	aday, _ := time.ParseInLocation("2006-01-02", day, time.Local)

	var t Organizational
	t.Id = cid
	db.Read(&t, "Id")

	var r AttendanceRule
	r.Kid = t.KindergartenId
	db.Read(&r, "Kid")

	beego.Debug("KID:", r.Kid)

	//	mbeg, _ := time.ParseInLocation("2006-01-02 15:04", today+" "+r.Mbeg)
	mend, _ := time.ParseInLocation("2006-01-02 15:04", day+" "+r.Mend, time.Local)
	abeg, _ := time.ParseInLocation("2006-01-02 15:04", day+" "+r.Abeg, time.Local)
	aend, _ := time.ParseInLocation("2006-01-02 15:04", day+" "+r.Aend, time.Local)

	rt = make(map[string]interface{})

	var sql string
	var r1 []orm.Params
	sql = "select t1.sid,t1.name,t1.cls,date_format(t1.morning, '%H:%i') morning,date_format(t1.afternoon, '%H:%i') afternoon,t4.avatar from attendance t1, organizational t2,organizational_member t3 , student t4 where t2.id=? and t1.sid=t3.member_id and t2.id=t3.organizational_id and t1.today = ? and t1.sid=t4.student_id and t1.morning < ? and t1.afternoon between ? and ? "
	db.Raw(sql, cid, day, mend, abeg, aend).Values(&r1)

	sql = "select t1.sid,t1.name,-1 type,'' reason,date_format(t1.morning, '%H:%i') morning,date_format(t1.afternoon, '%H:%i') afternoon,t4.avatar from attendance t1, organizational t2,organizational_member t3 , student t4 where t2.id=? and t1.sid=t3.member_id and t2.id=t3.organizational_id and t1.today = ? and t1.sid=t4.student_id and ( t1.morning > ? or t1.morning is null ) and ( t1.afternoon < ? or t1.afternoon > ? or t1.afternoon is null) "
	sql += " union all "
	sql += " select t1.sid,t4.name ,t1.type,t1.reason,null morning,null afternoon,t4.avatar from aleave t1,organizational t2,organizational_member t3, student t4 where t4.student_id=t1.sid and t2.id = ? and t1.sid=t3.member_id and t2.id=t3.organizational_id and t1.beg < ? and t1.end >?"
	var r2 []orm.Params
	db.Raw(sql, cid, day, mend, abeg, aend, cid, aday, aday).Values(&r2)
	rt["att"] = r1
	rt["leave"] = r2
	return
}

/*
*考勤总数统计
 */
func TotalCounting(day string, kid int) (rt []orm.Params) {
	db := orm.NewOrm()
	var sql string
	sql = "select max(case name when 'denominator' then n else 0 end) 'denominator', max(case name when 'numerator' then n else 0 end) 'numerator' " +
		" from (" +
		" select 'denominator' name,count(*) n from student t1 where kindergarten_id = ? and status=1 " +
		" union all" +
		" select 'numerator' name,count(*) n from attendance t2,student t3 where t2.morning is not null and t2.afternoon is not null and t3.kindergarten_id =? and t3.student_id=t2.sid )z "
	db.Raw(sql, kid, day).Values(&rt)
	return
}

/*
* 根据学校id求年级
 */
func GotGradeByKid(kid int) (rt []orm.Params) {
	db := orm.NewOrm()
	var sql string
	sql = "select t2.id,t2.name from organizational t ,organizational t2 where t.kindergarten_id = ? and t.level=1 and t.type =2 and t.id = t2.parent_id"
	db.Raw(sql, kid).Values(&rt)
	return
}

func GotRule(kid int) (r AttendanceRule) {
	db := orm.NewOrm()
	r.Kid = kid
	db.Read(&r, "Kid")
	return
}
