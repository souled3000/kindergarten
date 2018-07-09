package models

import (
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
)

type ExceptionalChild struct {
	Id             int    `json:"id";orm:"column(id);auto"`
	ChildName      string `json:"child_name";orm:"column(child_name);size(20)"`
	Class          int    `json:"class";orm:"column(class)"`
	Somatotype     int8   `json:"somatotype";orm:"column(somatotype)"`
	Allergen       string `json:"allergen";orm:"column(allergen);size(20)"`
	Source         int8   `json:"source";orm:"column(source)"`
	KindergartenId int    `json:"kindergarten_id";orm:"column(kindergarten_id)"`
	Creator        int    `json:"creator";orm:"column(creator)"`
	StudentId      int    `json:"student_id";orm:"column(student_id)"`
	CreatedAt      string `json:"created_at";orm:"column(created_at);type(datetime)"`
	UpdatedAt      string `json:"updated_at";orm:"column(updated_at);type(datetime)"`
}

func (t *ExceptionalChild) TableName() string {
	return "exceptional_child"
}

func init() {
	orm.RegisterModel(new(ExceptionalChild))
}

// AddExceptionalChild insert a new ExceptionalChild into database and returns
// last inserted Id on success.
func AddExceptionalChild(child_name string, class int, somatotype int8, allergen string, source int8, kindergarten_id int, creator int, student_id int) (id int64, err error) {
	var exceptionalChild ExceptionalChild
	o := orm.NewOrm()
	var infos []orm.Params
	where := " allergen like \"%" + string(allergen) + "%\" AND student_id = ? AND kindergarten_id = ? "
	if n, er := o.Raw("SELECT allergen FROM `exceptional_child` WHERE "+where, student_id, kindergarten_id).Values(&infos); er == nil && n > 0 {
		return 0, err
	} else {
		o.Begin()
		exceptionalChild.ChildName = child_name
		exceptionalChild.Class = class
		exceptionalChild.Somatotype = somatotype
		exceptionalChild.Allergen = allergen
		exceptionalChild.Source = source
		exceptionalChild.KindergartenId = kindergarten_id
		exceptionalChild.Creator = creator
		exceptionalChild.StudentId = student_id
		exceptionalChild.CreatedAt = time.Now().Format("2006-01-02 15:04:05")
		exceptionalChild.UpdatedAt = time.Now().Format("2006-01-02 15:04:05")
		if id, err := o.Insert(&exceptionalChild); err != nil && id <= 0 {
			o.Rollback()
			return id, err
		} else {
			o.Commit()
			return id, nil
		}
	}
	return
}

// GetExceptionalChildById retrieves ExceptionalChild by Id. Returns error if
// Id doesn't exist
func GetExceptionalChildById(id string) (exceptionalChild interface{}, err error) {
	o := orm.NewOrm()
	qb, _ := orm.NewQueryBuilder("mysql")
	sql := qb.Select("ex.id, ex.child_name, ex.somatotype, ex.allergen, ex.source, ex.kindergarten_id, ex.creator, ex.student_id, ex.created_at, ex.updated_at, o.id as class_id, o.name, o.class_type").From("exceptional_child as ex").LeftJoin("organizational as o").On("o.id = ex.class").Where("ex.id=" + id).String()
	var maps []orm.Params
	if num, err := o.Raw(sql).Values(&maps); err == nil && num > 0 {
		var newMaps []orm.Params
		for _, v := range maps {
			if v["class_type"] == nil {
				v["class_name"] = nil
			} else if v["class_type"].(string) == "3" {
				v["class_name"] = "大班" + v["name"].(string) + ""

			} else if v["class_type"].(string) == "2" {
				v["class_name"] = "中班" + v["name"].(string) + ""
			} else {
				v["class_name"] = "小班" + v["name"].(string) + ""
			}

			delete(v, "class_type")
			delete(v, "name")
			newMaps = append(newMaps, v)
		}
		return newMaps, err
	} else {
		return nil, err
	}

}

// GetAllExceptionalChild retrieves all ExceptionalChild matches certain condition. Returns empty list if
// no records exist
func GetAllExceptionalChild(child_name string, somatotype int8, page int64, limit int64, keyword string) (Page, error) {
	o := orm.NewOrm()
	where := "1=1 "

	if child_name != "" {
		where += " AND ex.child_name like \"%" + string(child_name) + "%\""
	}

	if somatotype != 0 {
		Somatotype := strconv.FormatInt(int64(somatotype), 10)

		where += " AND ex.somatotype = " + string(Somatotype)
	}

	// 特殊儿童姓名或者特殊儿童过敏源
	if keyword != "" {
		where += " AND ex.child_name like \"%" + string(keyword) + "%\" OR ex.allergen like \"%" + string(keyword) + "%\""

	}

	totalqb, _ := orm.NewQueryBuilder("mysql")

	tatolsql := totalqb.Select("count(*)").From("exceptional_child as ex").Where(where).String()

	var total int64
	err := o.Raw(tatolsql).QueryRow(&total)

	if err == nil {
		if page <= 0 {
			page = 1
		}

		if limit <= 0 {
			limit = 10
		}

		offset := (page - 1) * limit
		qb, _ := orm.NewQueryBuilder("mysql")
		sql := qb.Select("ex.id, ex.child_name, ex.somatotype, ex.allergen, ex.source, ex.kindergarten_id, ex.creator, ex.student_id, ex.created_at, ex.updated_at, o.id as class_id, o.name, o.class_type").From("exceptional_child as ex").LeftJoin("organizational as o").On("o.id = ex.class").Where(where).OrderBy("ex.id").Desc().Limit(int(limit)).Offset(int(offset)).String()
		var maps []orm.Params
		if num, err := o.Raw(sql).Values(&maps); err == nil && num > 0 {

			var newMap []orm.Params
			for _, v := range maps {

				if v["class_type"] == nil {
					v["class_name"] = nil
				} else if v["class_type"].(string) == "3" {
					v["class_name"] = "大班" + v["name"].(string) + ""
				} else if v["class_type"].(string) == "2" {
					v["class_name"] = "中班" + v["name"].(string) + ""
				} else {
					v["class_name"] = "小班" + v["name"].(string) + ""
				}

				if v["somatotype"] == nil {
					v["somatot_info"] = nil
					v["allergen"] = "无"
				} else if v["somatotype"].(string) == "3" {
					v["somatot_info"] = v["allergen"].(string) + "过敏"
				} else if v["somatotype"].(string) == "2" {
					v["somatot_info"] = "肥胖"
					v["allergen"] = "无"
				} else {
					v["somatot_info"] = "瘦小"
					v["allergen"] = "无"
				}

				t := time.Now()
				currentTime := t.Unix() - (24 * 3600 * 3) //当前时间戳
				BeCharge := v["updated_at"]
				toBeCharge := BeCharge.(string)
				timeLayout := "2006-01-02 15:04:05"
				loc, _ := time.LoadLocation("Local")
				theTime, _ := time.ParseInLocation(timeLayout, toBeCharge, loc)
				sr := theTime.Unix() //将更新时间转为时间戳
				if sr > currentTime {
					v["new"] = 1 //最新
				} else {
					v["new"] = 0
				}

				delete(v, "updated_at")
				delete(v, "class_type")
				delete(v, "name")
				newMap = append(newMap, v)
			}
			pageNum := int(math.Ceil(float64(total) / float64(limit)))

			return Page{newMap, total, pageNum}, nil
		} else {
			return Page{}, err
		}
	}
	return Page{}, err
}

// UpdateExceptionalChild updates ExceptionalChild by Id and returns error if
// the record to be updated doesn't exist
func UpdateExceptionalChildById(id int, child_name string, class int, somatotype int8, allergen string, student_id int) (err error) {
	o := orm.NewOrm()
	err = o.Begin()
	exceptionalChild := ExceptionalChild{Id: id}
	if err = o.Read(&exceptionalChild); err == nil {
		if child_name != "" {
			exceptionalChild.ChildName = child_name
		}

		if class != 0 {
			exceptionalChild.Class = class
		}

		if somatotype != 0 {
			exceptionalChild.Somatotype = somatotype
		}

		if allergen != "" {
			exceptionalChild.Allergen = allergen
		}

		exceptionalChild.UpdatedAt = time.Now().Format("2006-01-02 15:04:05")

		if student_id != 0 {
			exceptionalChild.StudentId = student_id
		}

		if _, err := o.Update(&exceptionalChild); err == nil {
			o.Commit()
			return err
		} else {
			o.Rollback()
			return err
		}
	}
	o.Rollback()
	return err
}

// DeleteExceptionalChild deletes ExceptionalChild by Id and returns error if
// the record to be deleted doesn't exist
func DeleteExceptionalChild(id int) (err error) {
	o := orm.NewOrm()
	o.Begin()
	e := &ExceptionalChild{Id: id}
	// ascertain id exists in the database
	if err = o.Read(e); err == nil {
		if _, err = o.Delete(&ExceptionalChild{Id: id}); err == nil {
			o.Commit()
			return nil
		}
	}
	o.Rollback()
	return err
}

// 根据过敏源获取过敏儿童
func GetAllergenChild(allergen string) (allergenChild []map[string]interface{}, err error) {
	param := strings.Split(allergen, ",")
	o := orm.NewOrm()
	var allergens []orm.Params
	for _, v := range param {
		if v != "" {
			maps := make(map[string]interface{})

			where := " allergen like \"%" + string(v) + "%\""
			if _, err := o.Raw("SELECT allergen FROM `exceptional_child` WHERE " + where).Values(&allergens); err == nil {
				if childName, err := GetChildName(where); err == nil {
					maps["allergen"] = v
					maps["child_name"] = childName
					allergenChild = append(allergenChild, maps)
				} else {
					return allergenChild, err
				}
			} else {
				return allergenChild, err
			}
		}
	}

	return allergenChild, nil
}

// 获取过敏儿童名称
func GetChildName(where string) (childName string, err error) {
	o := orm.NewOrm()
	var lists []orm.ParamsList
	if _, err := o.Raw("SELECT child_name FROM `exceptional_child` WHERE " + where).ValuesList(&lists); err == nil {
		var str string
		for _, row := range lists {
			str += row[0].(string) + ","
			s := []rune(str)
			childName = string(s[:len(s)-1])
		}
		return childName, err
	}
	return childName, err
}

// 根据baby_id获取过敏源
func GetAllergen(id int) (allergen []interface{}, err error) {
	o := orm.NewOrm()
	qb, _ := orm.NewQueryBuilder("mysql")
	idStr := strconv.Itoa(id)
	where := " stu.baby_id=" + idStr
	var maps []orm.Params
	sql := qb.Select("ex.allergen, ex.id, COUNT(DISTINCT ex.allergen) as field").From("student as stu").LeftJoin("exceptional_child as ex").On("ex.student_id = stu.student_id").Where(where).GroupBy("ex.`allergen`").String()
	if num, err := o.Raw(sql).Values(&maps); err == nil && num > 0 {
		for _, row := range maps {
			if row["allergen"] != nil {
				delete(row, "field")
				allergen = append(allergen, row)
			}
		}
		return allergen, err
	}
	return nil, err
}

// 过敏食物报备
func AllergenPreparation(child_name string, somatotype int8, allergens string, source int8, kindergarten_id int, creator int, baby_id int) (id int64, err error) {
	var exceptionalChild ExceptionalChild
	o := orm.NewOrm()
	allergen := strings.Split(allergens, ",")
	o.Begin()
	var maps []orm.Params
	if num, err := o.Raw("SELECT stu.student_id,org.member_id FROM student as stu LEFT JOIN organizational_member as org ON stu.student_id = org.member_id WHERE baby_id = ? LIMIT 1", baby_id).Values(&maps); err == nil && num > 0 {
		student_id, _ := strconv.Atoi(maps[0]["student_id"].(string))
		class, _ := strconv.Atoi(maps[0]["member_id"].(string))
		for _, v := range allergen {
			if v != "" {
				var infos []orm.Params
				where := " allergen like \"%" + string(v) + "%\" AND student_id = ? AND kindergarten_id = ? "
				if n, er := o.Raw("SELECT allergen FROM `exceptional_child` WHERE "+where, student_id, kindergarten_id).Values(&infos); er != nil || n == 0 {
					exceptionalChild.Id = 0
					exceptionalChild.ChildName = child_name
					exceptionalChild.Class = class
					exceptionalChild.Somatotype = somatotype
					exceptionalChild.Allergen = v
					exceptionalChild.Source = source
					exceptionalChild.KindergartenId = kindergarten_id
					exceptionalChild.Creator = creator
					exceptionalChild.StudentId = student_id
					exceptionalChild.CreatedAt = time.Now().Format("2006-01-02 15:04:05")
					exceptionalChild.UpdatedAt = time.Now().Format("2006-01-02 15:04:05")
					if id, err = o.Insert(&exceptionalChild); err != nil && id <= 0 {
						o.Rollback()
						return id, err
					}
				}
			} else {
				o.Rollback()
				return id, err
			}
		}
		o.Commit()
		return id, nil
	}
	return id, err
}
