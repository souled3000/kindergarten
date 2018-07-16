package models

import (
	"fmt"
	"github.com/astaxie/beego/orm"
	"strconv"
	"strings"
	"time"
)

type KindergartenTime struct {
	Id             int       `json:"id" orm:"column(id);auto"`
	Content        string    `json:"content" orm:"column(content);size(1000)" description:"内容"`
	KindergartenId int       `json:"kindergarten_id" orm:"column(kindergarten_id)"`
	CreatedAt      time.Time `json:"created_at" orm:"auto_now_add"`
	ClassType      int       `json:"class_type" orm:"class_type"`
	ClassId        int       `json:"class_id" orm:"class_id"`
	BeginTime      string    `json:"begin_time" orm:"column(begin_time);size(30)"`
	EndTime        string    `json:"end_time" orm:"column(end_time);size(30)"`
	Type           int       `json:"type" orm:"column(type)"`
	Name           string    `json:"name" orm:"column(name)"`
}

func (t *KindergartenTime) TableName() string {
	return "kindergarten_time"
}

func init() {
	orm.RegisterModel(new(KindergartenTime))
}

//添加时间段
func AddKindergartenTime(m KindergartenTime) (info interface{}, err error) {
	o := orm.NewOrm()
	fmt.Println(m)
	if _, err := o.Insert(&m); err != nil {
		return m, nil
	}
	return nil, err
}

//修改时间段
func UpdataKindergartenTime(m []KindergartenTime) (err error) {
	fmt.Println(m)
	o := orm.NewOrm()
	for _, val := range m {
		if _, err = o.Update(&val); err != nil {
			return err
		}
	}
	return nil
}

//时间安排
func GetKindergartenTimeInfo(class_type int, class_id int) map[string]interface{} {
	var v []KindergartenTime
	o := orm.NewOrm()
	where := "1=1"
	if class_id > 0 {
		where += " and type = 1 and class_id =" + strconv.Itoa(class_id)
	} else {
		where += " and type = 0 and class_type =" + strconv.Itoa(class_type)
	}
	sql := "select * from kindergarten_time where " + where + " order by begin_time asc"
	_, err := o.Raw(sql).QueryRows(&v)

	if err == nil {
		paginatorMap := make(map[string]interface{})
		paginatorMap["data"] = v
		return paginatorMap
	}
	return nil
}

//班级课程表
func GetClassTime(class_id int, kindergarten_time int, times string) map[string]interface{} {
	numbers := make(map[string]int64)
	numbers["Monday"] = 0
	numbers["Tuesday"] = 1
	numbers["Wednesday"] = 2
	numbers["Thursday"] = 3
	numbers["Friday"] = 4
	numbers["Saturday"] = 5
	numbers["Sunday"] = 6
	timeLayout := "2006-01-02"
	loc, _ := time.LoadLocation("")
	theTime, _ := time.ParseInLocation(timeLayout, times, loc)
	begin_num := numbers[theTime.Weekday().String()]
	end_num := 6 - begin_num
	begin_time := time.Unix(theTime.Unix()-begin_num*3600*24, 0).Format(timeLayout)
	end_time := time.Unix(theTime.Unix()+end_num*3600*24, 0).Format(timeLayout)
	sql := "select a.date,b.name,c.begin_time,c.end_time from course_time a left join course_info b on b.id = a.course_id left join kindergarten_time c on c.id= a.kindergarten_time_id where c.class_id = " + strconv.Itoa(class_id) + " and a.date >= '" + begin_time + "' and a.date <='" + end_time + "'"
	o := orm.NewOrm()
	var maps []orm.Params
	if _, err := o.Raw(sql).Values(&maps); err == nil {
		list := make(map[string][]interface{})
		for _, val := range maps {
			one := make(map[string]interface{})
			one[val["date"].(string)] = val
			list[val["begin_time"].(string)+"-"+val["end_time"].(string)] = append(list[val["begin_time"].(string)+"-"+val["end_time"].(string)], val)
		}
		paginatorMap := make(map[string]interface{})
		paginatorMap["data"] = list
		return paginatorMap
	}
	return nil
}

//班级某一天课程
func GetClassDay(class_id int, kindergarten_id int, times string) map[string]interface{} {
	o := orm.NewOrm()
	sql := "select * from kindergarten_time where (type=0 and class_id = 0 and kindergarten_id=" + strconv.Itoa(kindergarten_id) + ") or (type=1 and class_id = " + strconv.Itoa(class_id) + " and kindergarten_id=" + strconv.Itoa(kindergarten_id) + ") order by begin_time asc"
	var v []KindergartenTime
	_, err := o.Raw(sql).QueryRows(&v)
	numbers := make(map[string]int64)
	numbers["Monday"] = 0
	numbers["Tuesday"] = 1
	numbers["Wednesday"] = 2
	numbers["Thursday"] = 3
	numbers["Friday"] = 4
	numbers["Saturday"] = 5
	numbers["Sunday"] = 6
	timeLayout := "2006-01-02"
	loc, _ := time.LoadLocation("")
	theTime, _ := time.ParseInLocation(timeLayout, times, loc)
	begin_num := numbers[theTime.Weekday().String()]

	for key, val := range v {
		namelist := strings.Split(val.Name, ",")
		if len(namelist) > int(begin_num) {
			v[key].Name = namelist[int(begin_num)]
		}
	}

	if err == nil {
		paginatorMap := make(map[string]interface{})
		paginatorMap["data"] = v
		return paginatorMap
	}
	return nil
}
