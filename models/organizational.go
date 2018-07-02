package models

import (
	"errors"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
)

type Organizational struct {
	Id             int       `json:"id" orm:"column(id);auto"`
	KindergartenId int       `json:"kindergarten_id" orm:"column(kindergarten_id)" description:"幼儿园id"`
	ParentId       int       `json:"parent_id" orm:"column(parent_id)" description:"父级id"`
	Name           string    `json:"name" orm:"column(name);size(20)" description:"组织架构名字"`
	IsFixed        int8      `json:"is_fixed" orm:"column(is_fixed)" description:"是否固定的：0不是，1是"`
	Level          int8      `json:"level" orm:"column(level)" description:"等级"`
	ParentIds      string    `json:"parent_ids" orm:"column(parent_ids);size(50)" description:"父级所有id"`
	Type           int8      `json:"type" orm:"column(type)" description:"类型：0普通，1管理层，2年级组"`
	ClassType      int8      `json:"class_type" orm:"column(class_type)" description:"班级类型：1小班，2中班，3大班"`
	CreatedAt      time.Time `json:"created_at" orm:"column(created_at);type(datetime)" description:"添加时间"`
	UpdatedAt      time.Time `json:"updated_at" orm:"column(updated_at);type(datetime)" description:"修改时间"`
}

type OrganizationalTree struct {
	Id             int                  `json:"id" orm:"column(id);auto"`
	KindergartenId int                  `json:"kindergarten_id" orm:"column(kindergarten_id)" description:"幼儿园id"`
	ParentId       int                  `json:"parent_id" orm:"column(parent_id)" description:"父级id"`
	Name           string               `json:"name" orm:"column(name);size(20)" description:"组织架构名字"`
	IsFixed        int8                 `json:"is_fixed" orm:"column(is_fixed)" description:"是否固定的：0不是，1是"`
	Level          int8                 `json:"level" orm:"column(level)" description:"等级"`
	ParentIds      string               `json:"parent_ids" orm:"column(parent_ids);size(50)" description:"父级所有id"`
	Type           int8                 `json:"type" orm:"column(type)" description:"类型：0普通，1管理层，2年级组"`
	ClassType      int8                 `json:"class_type" orm:"column(class_type)" description:"班级类型：1小班，2中班，3大班"`
	CreatedAt      time.Time            `json:"created_at" orm:"column(created_at);type(datetime)" description:"添加时间"`
	UpdatedAt      time.Time            `json:"updated_at" orm:"column(updated_at);type(datetime)" description:"修改时间"`
	Children       []OrganizationalTree `json:"children" orm:"null"`
}

func (t *Organizational) TableName() string {
	return "organizational"
}

func init() {
	orm.RegisterModel(new(Organizational))
}

/*
班级搜索
*/
func GetClassAll(kindergarten_id int, class_type int, page int, prepage int) (paginatorMap map[string]interface{}, err error) {
	o := orm.NewOrm()
	qb, _ := orm.NewQueryBuilder("mysql")
	var condition []interface{}
	where := "1=1 "
	if kindergarten_id == 0 {
		where += " AND kindergarten_id = ?"
	} else {
		where += " AND kindergarten_id = ?"
		condition = append(condition, kindergarten_id)
	}
	if class_type > 0 {
		where += " AND class_type = ?"
		condition = append(condition, class_type)
	}
	// 构建查询对象
	sql := qb.Select("count(*)").From("organizational").Where(where).And("type = 2").And("level = 3").String()
	var total int64
	err = o.Raw(sql, condition).QueryRow(&total)
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
		sql := qb.Select("*").From("organizational").Where(where).And("type = 2").And("level = 3").Limit(prepage).Offset(limit).String()
		num, err := o.Raw(sql, condition).Values(&v)
		if err == nil && num > 0 {
			paginatorMap := make(map[string]interface{})
			paginatorMap["total"] = total         //总条数
			paginatorMap["data"] = v              //分页数据
			paginatorMap["page_num"] = totalpages //总页数
			return paginatorMap, nil
		}
	}
	err = errors.New("获取失败")
	return nil, err
}

/*
班级成员
*/
func ClassMember(kindergarten_id int, class_type int, class_id int, page int, prepage int) (paginatorMap map[string]interface{}, err error) {
	o := orm.NewOrm()
	var student []orm.Params
	var teacher []orm.Params
	var condition []interface{}
	where := "1=1 "
	if kindergarten_id == 0 {
		where += " AND o.kindergarten_id = ?"
	} else {
		where += " AND o.kindergarten_id = ?"
		condition = append(condition, kindergarten_id)
	}
	if class_type > 0 {
		where += " AND o.class_type = ?"
		condition = append(condition, class_type)
	}
	if class_id > 0 {
		where += " AND o.id = ?"
		condition = append(condition, class_id)
	}
	where += " AND o.type = ?"
	condition = append(condition, 2)
	where += " AND o.level = ?"
	condition = append(condition, 3)

	qb, _ := orm.NewQueryBuilder("mysql")
	sql := qb.Select("s.student_id", "s.name", "s.avatar", "s.number", "s.phone",
		"o.name as class_name", "o.class_type", "om.id").From("student as s").LeftJoin("organizational_member as om").
		On("s.student_id = om.member_id").LeftJoin("organizational as o").
		On("om.organizational_id = o.id").Where(where).And("om.type = 1").And("s.status = 1").String()
	_, err = o.Raw(sql, condition).Values(&student)

	qb, _ = orm.NewQueryBuilder("mysql")
	sql = qb.Select("t.teacher_id", "t.name", "t.avatar", "t.number", "t.phone",
		"o.name as class_name", "om.id").From("teacher as t").LeftJoin("organizational_member as om").
		On("t.teacher_id = om.member_id").LeftJoin("organizational as o").
		On("om.organizational_id = o.id").Where(where).And("t.status = 1").And("om.type = 0").String()
	_, err = o.Raw(sql, condition).Values(&teacher)
	if err == nil {
		paginatorMap := make(map[string]interface{})
		paginatorMap["student"] = student
		paginatorMap["teacher"] = teacher
		return paginatorMap, nil
	}
	err = errors.New("获取失败")
	return nil, err
}

/*
删除班级
*/
func Destroy(class_id int) (paginatorMap map[string]interface{}, err error) {
	o := orm.NewOrm()
	err = o.Begin()
	var t []orm.Params
	var s []orm.Params
	var v Organizational
	paginatorMap = make(map[string]interface{})
	o.QueryTable("organizational").Filter("id", class_id).All(&v)
	if v.IsFixed == 1 {
		err = errors.New("不能删除")
		return nil, err
	}
	qb, _ := orm.NewQueryBuilder("mysql")
	sql := qb.Select("member_id").From("organizational_member").Where("organizational_id = ?").And("type = 0").String()
	_, err = o.Raw(sql, class_id).Values(&t)
	if err == nil {
		//修改teacher
		for key, _ := range t {
			_, err = o.QueryTable("teacher").Filter("teacher_id", t[key]["member_id"]).Update(orm.Params{
				"status": 0,
			})
		}

	}

	//修改学生
	qb, _ = orm.NewQueryBuilder("mysql")
	sql = qb.Select("member_id").From("organizational_member").Where("organizational_id = ?").And("type = 1").String()
	_, err = o.Raw(sql, class_id).Values(&s)
	if err == nil {
		for key, _ := range s {
			_, err = o.QueryTable("student").Filter("student_id", s[key]["member_id"]).Update(orm.Params{
				"status": 0,
			})
		}
	}

	//删除班级
	_, err = o.QueryTable("organizational").Filter("id", class_id).Delete()
	if err != nil {
		err = o.Rollback()
	}
	//删除班级成员
	_, err = o.QueryTable("organizational_member").Filter("organizational_id", class_id).Delete()
	if err != nil {
		err = o.Rollback()
	}
	if err == nil {
		err = o.Commit()
		return paginatorMap, nil
	}
	err = errors.New("删除失败")
	return nil, err
}

/*
创建班级
*/
func StoreClass(class_type int, kindergarten_id int) (paginatorMap map[string]interface{}, err error) {
	o := orm.NewOrm()
	err = o.Begin()
	var or []orm.Params
	var max_name []orm.Params
	paginatorMap = make(map[string]interface{})
	var condition []interface{}
	where := "1=1 "
	if class_type > 0 {
		where += " AND class_type = ?"
		condition = append(condition, class_type)
	}
	if kindergarten_id > 0 {
		where += " AND kindergarten_id = ?"
		condition = append(condition, kindergarten_id)
	}
	qb, _ := orm.NewQueryBuilder("mysql")
	sql := qb.Select("o.*").From("organizational as o").Where(where).And("type = 2").And("level = 2").String()
	num, err := o.Raw(sql, condition).Values(&or)
	//查出最大班级
	qb, _ = orm.NewQueryBuilder("mysql")
	sql = qb.Select("max(CONVERT(o.name,SIGNED)) as m").From("organizational as o").Where(where).And("type = 2").And("level = 3").String()
	_, err = o.Raw(sql, condition).Values(&max_name)
	if max_name[0]["m"] == nil {
		max_name[0]["m"] = "0"
	}
	m := max_name[0]["m"].(string)
	//班级号
	ml := strings.Replace(m, "班", "", -1)
	new_name, _ := strconv.Atoi(ml)
	name_number := new_name + 1
	name := strconv.Itoa(name_number)
	if num == 0 {
		err = errors.New("班级不存在")
		return nil, err
	}
	//interface 转 int
	le := or[0]["level"].(string)
	level, _ := strconv.Atoi(le)
	lev := level + 1
	//创建 name+1
	parent_ids := or[0]["parent_ids"].(string)
	ids := or[0]["id"].(string)
	qb, _ = orm.NewQueryBuilder("mysql")
	sql = "insert into organizational set kindergarten_id = ?,name = ?,level = ?,parent_ids = ?,class_type = ?,type = ?,parent_id = ?,is_fixed =?"
	res, err := o.Raw(sql, kindergarten_id, ""+name+"班", lev, ""+parent_ids+""+ids+",", class_type, 2, or[0]["id"], or[0]["is_fixed"]).Exec()
	id, _ := res.LastInsertId()
	if err != nil {
		o.Rollback()
	} else {
		o.Commit()
	}
	if err == nil {
		paginatorMap["class_id"] = id
		paginatorMap["name"] = "" + name + "班"
		paginatorMap["class_type"] = class_type
		return paginatorMap, nil
	}
	err = errors.New("创建失败")
	return nil, err
}

/*
组织架构列表
*/
func GetOrganization(kindergarten_id int, page int, prepage int) map[string]interface{} {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Organizational))
	var posts []Organizational
	var Organizational []OrganizationalTree
	if _, err := qs.Filter("kindergarten_id", kindergarten_id).All(&posts); err == nil {
		ml := make(map[string]interface{})
		for _, val := range posts {
			if val.ParentId == 0 {
				next := getNext(posts, val.Id)
				var tree OrganizationalTree
				tree.Id = val.Id
				tree.KindergartenId = val.KindergartenId
				tree.ClassType = val.ClassType
				tree.CreatedAt = val.CreatedAt
				tree.IsFixed = val.IsFixed
				tree.Level = val.Level
				tree.Name = val.Name
				tree.ParentId = val.ParentId
				tree.ParentIds = val.ParentIds
				tree.UpdatedAt = val.UpdatedAt
				tree.Children = next
				Organizational = append(Organizational, tree)
			}
		}
		if err == nil {
			ml["data"] = Organizational
			return ml
		}
	}
	return nil
}

func getNext(posts []Organizational, id int) (Organizational []OrganizationalTree) {
	for _, val := range posts {
		if val.ParentId == id {
			next := getNext(posts, val.Id)
			var tree OrganizationalTree
			tree.Id = val.Id
			tree.KindergartenId = val.KindergartenId
			tree.ClassType = val.ClassType
			tree.CreatedAt = val.CreatedAt
			tree.IsFixed = val.IsFixed
			tree.Level = val.Level
			tree.Name = val.Name
			tree.ParentId = val.ParentId
			tree.ParentIds = val.ParentIds
			tree.UpdatedAt = val.UpdatedAt
			tree.Type = val.Type
			tree.Children = next
			Organizational = append(Organizational, tree)
		}
	}
	return Organizational
}

/*
添加组织架构
*/
func AddOrganization(name string, ty int, parent_id int, kindergarten_id int, class_type int) (paginatorMap map[string]interface{}, err error) {
	o := orm.NewOrm()
	var v []orm.Params
	var kinder []orm.Params
	paginatorMap = make(map[string]interface{})
	qb, _ := orm.NewQueryBuilder("mysql")
	sql := qb.Select("*").From("kindergarten").Where("kindergarten_id = ?").String()
	num, err := o.Raw(sql, kindergarten_id).Values(&kinder)
	if num == 0 {
		err = errors.New("幼儿园不存在")
		return nil, err
	}
	if parent_id != 0 {
		qb, _ = orm.NewQueryBuilder("mysql")
		sql = qb.Select("*").From("organizational").Where("id = ?").String()
		num, err = o.Raw(sql, parent_id).Values(&v)
		if num == 0 {
			err = errors.New("上一级架构不存在")
			return nil, err
		}
		//interface 转 int
		parent_ids := v[0]["parent_ids"].(string)
		id := v[0]["id"].(string)

		t := v[0]["type"].(string)
		typ, _ := strconv.Atoi(t)

		lev := v[0]["level"].(string)
		leve, _ := strconv.Atoi(lev)
		le := leve + 1
		if typ == 1 {
			if leve >= 2 {
				err = errors.New("管理层不能超过2级")
				return nil, err
			}
		} else {
			if leve >= 3 {
				err = errors.New("管理层不能超过3级")
				return nil, err
			}
		}
		qb, _ = orm.NewQueryBuilder("mysql")
		sql = "insert into organizational set parent_id = ?,name = ?,level = ?,parent_ids = ?,type = ?,kindergarten_id =?,class_type = ?"
		_, err = o.Raw(sql, parent_id, name, le, ""+parent_ids+""+id+",", ty, kindergarten_id, class_type).Exec()
	} else {
		qb, _ = orm.NewQueryBuilder("mysql")
		sql = "insert into organizational set name = ?,type = ?,kindergarten_id =?"
		_, err = o.Raw(sql, name, ty, kindergarten_id).Exec()
	}
	if err == nil {
		paginatorMap["data"] = nil
		return paginatorMap, nil
	}
	return nil, err
}

/*
删除组织架构
*/
func DelOrganization(organization_id int) (paginatorMap map[string]interface{}, err error) {
	o := orm.NewOrm()
	var v []orm.Params
	var val []orm.Params
	var organ []orm.Params
	paginatorMap = make(map[string]interface{})
	qb, _ := orm.NewQueryBuilder("mysql")
	sql := qb.Select("*").From("organizational").Where("id = ?").String()
	_, err = o.Raw(sql, organization_id).Values(&v)
	is_fixe := v[0]["is_fixed"].(string)
	level := v[0]["level"].(string)
	if is_fixe == "1" || (is_fixe == "0" && level == "2") {
		err = errors.New("不能删除")
		return nil, err
	}
	qb, _ = orm.NewQueryBuilder("mysql")
	sql = qb.Select("count(*) as num").From("organizational").Where("parent_ids = ?").String()
	_, err = o.Raw(sql, organization_id).Values(&val)
	num := val[0]["num"].(string)
	nums, _ := strconv.Atoi(num)
	if nums > 0 {
		qb, _ = orm.NewQueryBuilder("mysql")
		sql = qb.Select("organizational.*").From("organizational").Where("parent_ids = ?").String()
		_, err = o.Raw(sql, organization_id).Values(&organ)
		for key, _ := range organ {
			_, err = o.QueryTable("organizational").Filter("id", organ[key]["id"]).Delete()
			_, err = o.QueryTable("organizational_member").Filter("organizational_id", organ[key]["id"]).Delete()
		}
	} else {
		_, err = o.QueryTable("organizational").Filter("id", organization_id).Delete()
		_, err = o.QueryTable("teacher").Filter("teacher_id", organization_id).Delete()
	}
	if err == nil {
		paginatorMap["data"] = nil
		return paginatorMap, nil
	}
	return nil, err
}

/*
编辑组织架构
*/
func UpOrganization(organization_id int, name string) (paginatorMap map[string]interface{}, err error) {
	o := orm.NewOrm()
	var v []orm.Params
	paginatorMap = make(map[string]interface{})
	qb, _ := orm.NewQueryBuilder("mysql")
	sql := qb.Select("*").From("organizational").Where("id = ?").String()
	_, err = o.Raw(sql, organization_id).Values(&v)
	//	level := v[0]["level"].(string)
	is_fixe := v[0]["is_fixed"].(string)
	if is_fixe == "1" {
		err = errors.New("不能编辑")
		return nil, err
	}
	_, err = o.QueryTable("organizational").Filter("id", organization_id).Update(orm.Params{
		"name": name,
	})
	if err == nil {
		paginatorMap["data"] = nil
		return paginatorMap, nil
	}
	err = errors.New("编辑组织架构失败")
	return nil, err
}

/*
班级成员
*/
func Principal(class_id int, page int, prepage int) map[string]interface{} {
	o := orm.NewOrm()
	var teacher []orm.Params
	var condition []interface{}
	where := "1=1 "
	if class_id > 0 {
		where += " AND om.organizational_id = ?"
		condition = append(condition, class_id)
	}
	qb, _ := orm.NewQueryBuilder("mysql")
	sql := qb.Select("t.*", "om.id").From("teacher as t").LeftJoin("organizational_member as om").
		On("t.teacher_id = om.member_id").Where(where).And("om.type = 0").String()
	num, err := o.Raw(sql, condition).Values(&teacher)
	if err == nil && num > 0 {
		paginatorMap := make(map[string]interface{})
		paginatorMap["data"] = teacher
		return paginatorMap
	}
	return nil
}

/*
幼儿园所有班级
*/
func GetkinderClass(kindergarten_id int, user_id int) (paginatorMapmap map[string]interface{}, err error) {
	o := orm.NewOrm()
	var v []orm.Params
	var teacher []orm.Params
	qb, _ := orm.NewQueryBuilder("mysql")
	paginatorMap := make(map[string]interface{})
	sql := qb.Select("class_type", "id as class_id", "name as class_name").From("organizational").Where("kindergarten_id = ?").And("type = 2").And("level = 3").String()
	_, err = o.Raw(sql, kindergarten_id).Values(&v)
	if v == nil {
		err = errors.New("未创建班级")
		return nil, err
	}
	for key, val := range v {
		if val["class_type"].(string) == "3" {
			v[key]["class"] = "大班" + val["class_name"].(string) + ""
		} else if val["class_type"].(string) == "2" {
			v[key]["class"] = "中班" + val["class_name"].(string) + ""
		} else {
			v[key]["class"] = "小班" + val["class_name"].(string) + ""
		}
	}
	qb, _ = orm.NewQueryBuilder("mysql")
	sql = qb.Select("teacher_id").From("teacher").Where("user_id = ?").String()
	_, err = o.Raw(sql, user_id).Values(&teacher)
	if teacher == nil && err == nil {
		paginatorMap["data"] = v
		paginatorMap["teacher_id"] = nil
		return paginatorMap, nil
	} else {
		paginatorMap["data"] = v
		paginatorMap["teacher_id"] = teacher[0]["teacher_id"]
		return paginatorMap, nil
	}
	return nil, err
}

/*
幼儿园所有班级学生
*/
func GetClassStudent(class_id int) (paginatorMapmap map[string]interface{}, err error) {
	o := orm.NewOrm()
	var v []orm.Params
	qb, _ := orm.NewQueryBuilder("mysql")
	paginatorMap := make(map[string]interface{})
	qb, _ = orm.NewQueryBuilder("mysql")
	sql := qb.Select("s.*").From("organizational_member as om").LeftJoin("student as s").
		On("om.member_id = s.student_id").Where("om.organizational_id = ?").And("om.type = 1").String()
	num, err := o.Raw(sql, class_id).Values(&v)
	if err == nil && num > 0 {
		paginatorMap["data"] = v
		return paginatorMap, nil
	}
	return nil, err
}

/*
宝宝所在班级
*/
func GetBabyClass(babyIds string) (paginatorMapmap map[string]interface{}, err error) {
	o := orm.NewOrm()
	var v []orm.Params
	var class []interface{}
	baby_id := strings.Split(babyIds, ",")
	paginatorMapmap = make(map[string]interface{})
	for _, s := range baby_id {
		qb, _ := orm.NewQueryBuilder("mysql")
		sql := qb.Select("s.name", "s.student_id", "s.baby_id", "o.name as class_name", "o.id as class_id", "class_type").From("student as s").LeftJoin("organizational_member as om").
			On("s.student_id = om.member_id").LeftJoin("organizational as o").
			On("om.organizational_id = o.id").Where("s.baby_id = ?").And("s.status = 1").And("om.type = 1").And("om.is_principal = 1").And("isnull(s.deleted_at)").String()
		_, err = o.Raw(sql, s).Values(&v)
		for key, val := range v {
			if val["class_type"].(string) == "3" {
				v[key]["class"] = "大班" + val["class_name"].(string) + ""
			} else if val["class_type"].(string) == "2" {
				v[key]["class"] = "中班" + val["class_name"].(string) + ""
			} else {
				v[key]["class"] = "小班" + val["class_name"].(string) + ""
			}
		}
		class = append(class, v)
	}
	if err == nil {
		paginatorMapmap["data"] = class
		return paginatorMapmap, nil
	}
	return nil, err
}
