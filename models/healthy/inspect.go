package healthy

import (
	"encoding/json"
	"fmt"

	"github.com/astaxie/beego/orm"

	"kindergarten-service-go/models"
	"math"
	"strconv"
	"strings"
	"time"
)

type Inspect struct {
	Id             int       `json:"id" orm:"column(id);auto" description:"编号"`
	StudentId      int       `json:"student_id" orm:"column(student_id)" description:"学生ID"`
	Weight         float64   `json:"weight" orm:"column(weight)" description:"体重"`
	AbnormalWeight string    `json:"abnormal_weight" orm:"column(abnormal_weight)"`
	ClassName      string    `json:"class_name" orm:"column" description:"班级名称"`
	Height         float64   `json:"height" orm:"column(height)" description:"身高"`
	AbnormalHeight string    `json:"abnormal_height" orm:"column(abnormal_height)"`
	Content        string    `json:"content" orm:"column(content);size(100)" description:"备注"`
	Evaluate       int       `json:"evaluate" orm:"column(evaluate)" description:"评价"`
	ClassId        int       `json:"class_id" orm:"column(class_id)" description:"班级ID"`
	KindergartenId int       `json:"kindergarten_id" orm:"column(kindergarten_id)" description:"幼儿园ID"`
	TeacherId      int       `json:"teacher_id" orm:"column(teacher_id)" description:"记录人ID"`
	Types          int       `json:"types" orm:"column(types)" description:"类型"`
	Handel         string    `json:"handel" orm:"column(handel);size(100)" description:"处理方式"`
	Url            string    `json:"url" orm:"column(url);size(1800)" description:"照片"`
	Infect         int       `json:"infect" orm:"column(infect)" description:"是否传染"`
	DrugId         int       `json:"drug_id" orm:"column(drug_id)" description:"喂药申请"`
	Abnormal       string    `json:"abnormal" orm:"column(abnormal)" description:"是否传染"`
	BodyId         int       `json:"body_id" orm:"column(body_id)"`
	Date           string    `json:"date" orm:"column(date)" description:"参检时间"`
	CreatedAt      time.Time `json:"created_at" orm:"auto_now_add;auto_now" description:"创建时间"`
}

type Inspect_add struct {
	Id             int       `json:"id" orm:"column(id);auto" description:"编号"`
	StudentId      int       `json:"student_id" orm:"column(student_id)" description:"学生ID"`
	Weight         float64   `json:"weight" orm:"column(weight)" description:"体重"`
	ClassName      string    `json:"class_name" orm:"column" description:"班级名称"`
	Height         float64   `json:"height" orm:"column(height)" description:"身高"`
	Content        string    `json:"content" orm:"column(content);size(100)" description:"备注"`
	Evaluate       int       `json:"evaluate" orm:"column(evaluate)" description:"评价"`
	ClassId        int       `json:"class_id" orm:"column(class_id)" description:"班级ID"`
	KindergartenId int       `json:"kindergarten_id" orm:"column(kindergarten_id)" description:"幼儿园ID"`
	TeacherId      int       `json:"teacher_id" orm:"column(teacher_id)" description:"记录人ID"`
	Types          int       `json:"types" orm:"column(types)" description:"类型"`
	Handel         string    `json:"handel" orm:"column(handel);size(100)" description:"处理方式"`
	Url            string    `json:"url" orm:"column(url);size(1800)" description:"照片"`
	Infect         int       `json:"infect" orm:"column(infect)" description:"是否传染"`
	DrugId         int       `json:"drug_id" orm:"column(drug_id)" description:"喂药申请"`
	Abnormal       string    `json:"abnormal" orm:"column(abnormal)" description:"是否传染"`
	Date           string    `json:"date" orm:"column(date)" description:"参检时间"`
	CreatedAt      time.Time `json:"created_at" orm:"auto_now_add;auto_now" description:"创建时间"`
	BodyId         int       `json:"body_id" orm:"column(body_id)" description:"班级ID"`
	Info           Column    `json:"info" orm:""`
}

func (t *Inspect) TableName() string {
	return "healthy_inspect"
}

func init() {
	orm.RegisterModel(new(Inspect))
}

type Page struct {
	PageNum int          `json:"page_num"`
	PerPage int          `json:"per_page"`
	Total   int          `json:"total"`
	Data    []orm.Params `json:"data"`
}

//创建记录
func (m Inspect) Save() error {
	tmp := m
	o := orm.NewOrm()

	where := " where 1=1 "
	if tmp.StudentId > 0 {
		where += " and student_id = " + strconv.Itoa(tmp.StudentId)
	}

	if tmp.Date != "" {
		where += " and left(created_at,10) = '" + tmp.Date + "'"
	} else {
		day_time := time.Now().Format("2006-01-02")
		where += " and left(created_at,10) = '" + day_time + "'"
	}

	var drug []Drug
	_, err := o.Raw("SELECT id FROM healthy_drug " + where).QueryRows(&drug)

	if err == nil && len(drug) != 0 {
		m.DrugId = drug[0].Id
	}
	o.Insert(&m)

	return nil
}

//晨、午、晚列表
func (f *Inspect) GetAll(page, perPage, kindergarten_id, class_id, types, role, baby_id int, date, search string) (Page, error) {
	o := orm.NewOrm()
	var con []interface{}
	where := "1 "
	where += "AND ( healthy_inspect.types = 1 Or healthy_inspect.types = 2 Or healthy_inspect.types = 3 ) "
	if kindergarten_id != 0 {
		where += "AND healthy_inspect.kindergarten_id = ? "
		con = append(con, kindergarten_id)
	}
	if role != 1 {
		if class_id != 0 {
			where += "AND healthy_inspect.class_id = ? "
			con = append(con, class_id)
		}
	}
	if types != 0 {
		where += "AND healthy_inspect.types = ? "
		con = append(con, types)
	}

	if class_id != 0 {
		where += "AND healthy_inspect.class_id = ? "
		con = append(con, class_id)
	}

	if date == "" {
		day_time := time.Now().Format("2006-01-02")
		where += " AND left(healthy_inspect.date,10) = '" + day_time + "'"
	} else {
		where += " AND left(healthy_inspect.date,10) = '" + date + "'"
	}
	if search != "" {
		where += "AND ( student.name like ? Or teacher.name like ? Or healthy_inspect.abnormal like ? ) "
		con = append(con, "%"+search+"%")
		con = append(con, "%"+search+"%")
		con = append(con, "%"+search+"%")
	}

	if baby_id != 0 {
		var student_id int
		var student models.Student
		err := o.QueryTable("student").Filter("baby_id", baby_id).One(&student)
		if err != nil {
			student_id = 0
		} else {
			student_id = student.Id
		}
		where += "AND healthy_inspect.student_id = ? "
		con = append(con, student_id)
	}
	qb, _ := orm.NewQueryBuilder("mysql")
	sql := qb.Select("count(*)").From(f.TableName()).
		LeftJoin("student").On("healthy_inspect.student_id = student.student_id").
		LeftJoin("teacher").On("healthy_inspect.teacher_id = teacher.teacher_id").
		Where(where).String()

	var total int
	err := o.Raw(sql, con).QueryRow(&total)
	if err == nil {
		var sxWords []orm.Params

		limit := 10
		if perPage != 0 {
			limit = perPage
		}
		if page <= 0 {
			page = 1
		}
		start := (page - 1) * limit
		qb, _ := orm.NewQueryBuilder("mysql")
		sql := qb.Select("healthy_inspect.*,student.name as student_name,student.avatar,teacher.name as teacher_name").From(f.TableName()).
			LeftJoin("student").On("healthy_inspect.student_id = student.student_id").
			LeftJoin("teacher").On("healthy_inspect.teacher_id = teacher.teacher_id").
			Where(where).
			OrderBy("healthy_inspect.id").Desc().
			Limit(limit).Offset(start).String()

		if _, err := o.Raw(sql, con).Values(&sxWords); err == nil {

			pageNum := int(math.Ceil(float64(total) / float64(limit)))
			return Page{pageNum, limit, total, sxWords}, nil
		}
	}

	return Page{}, nil
}

//删除
func DeleteInspect(id int) map[string]interface{} {
	o := orm.NewOrm()
	v := Inspect{Id: id}
	if err := o.Read(&v); err == nil {
		if num, err := o.Delete(&Inspect{Id: id}); err == nil {
			paginatorMap := make(map[string]interface{})
			paginatorMap["data"] = num
			return paginatorMap
		}
	}
	return nil
}

//统计
func Counts(kindergarten_id int) map[string]interface{} {
	o := orm.NewOrm()
	counts := make(map[string]interface{})
	numbers := make(map[string]int64)
	numbers["Monday"] = 0
	numbers["Tuesday"] = 1
	numbers["Wednesday"] = 2
	numbers["Thursday"] = 3
	numbers["Friday"] = 4
	numbers["Saturday"] = 5
	numbers["Sunday"] = 6
	//每天统计
	where := " kindergarten_id = " + strconv.Itoa(kindergarten_id)
	day_time := time.Now().Format("2006-01-02")
	where += " AND left(date,10) = '" + day_time + "'"
	where += " AND (types = 1 Or types = 2 Or types = 3)"
	where += " AND abnormal != '' "
	type Counts struct {
		Num int
	}
	var count []Counts
	_, err := o.Raw("SELECT count(id) as num FROM healthy_inspect where" + where).QueryRows(&count)
	if err == nil {
		counts["day"] = count[0].Num
	}
	//实际检查人数
	where1 := " kindergarten_id = " + strconv.Itoa(kindergarten_id)
	where1 += " AND left(date,10) = '" + day_time + "'"
	_, err = o.Raw("SELECT count(id) as num FROM healthy_inspect where" + where1).QueryRows(&count)
	if err == nil {
		counts["day_actual"] = count[0].Num
	}
	//每月统计
	month_time := time.Now().Format("2006-01")
	where2 := " kindergarten_id = " + strconv.Itoa(kindergarten_id)
	where2 += " AND left(healthy_inspect.date,7) = '" + month_time + "'"
	where2 += " AND (abnormal != '' Or abnormal_weight = '瘦小' Or abnormal_weight = '肥胖' Or abnormal_weight = '矮小' Or abnormal_weight = '超高' Or abnormal2 = '异常' Or abnormal3 = '重度贫血') "
	_, err = o.Raw("SELECT count(healthy_inspect.id) as num FROM healthy_inspect left join healthy_column  on healthy_inspect.id = healthy_column.inspect_id where" + where2).QueryRows(&count)
	if err == nil {
		counts["month"] = count[0].Num
	}
	//每周统计
	t := time.Now()
	date := time.Now().Unix() - 24*60*60*(numbers[t.Weekday().String()])
	tm := time.Unix(date, 0)
	startTime := tm.Format("2006-01-02 00:00:00")

	where3 := " kindergarten_id = " + strconv.Itoa(kindergarten_id)
	where3 += " AND healthy_inspect.created_at >= '" + startTime + "'"
	where3 += " AND (abnormal != '' Or abnormal_weight = '瘦小' Or abnormal_weight = '肥胖' Or abnormal_weight = '矮小' Or abnormal_weight = '超高' Or abnormal2 = '异常' Or abnormal3 = '重度贫血') "
	_, err = o.Raw("SELECT count(healthy_inspect.id) as num FROM healthy_inspect left join healthy_column  on healthy_inspect.id = healthy_column.inspect_id where" + where3).QueryRows(&count)
	if err == nil {
		counts["week"] = count[0].Num
	}
	//全部统计
	where4 := " kindergarten_id = " + strconv.Itoa(kindergarten_id)
	where4 += " AND types = 1 Or types = 2 Or types = 3"
	all, err := o.Raw("SELECT count(id) as num FROM healthy_inspect where" + where4).QueryRows(&count)
	if err == nil {
		fmt.Println(all)
		counts["all"] = count[0].Num
	}

	return counts
}

//详情
func InspectInfo(id int) map[string]interface{} {
	o := orm.NewOrm()
	v := &Inspect{Id: id}
	if err := o.Read(v); err == nil {
		paginatorMap := make(map[string]interface{})
		paginatorMap["data"] = v
		return paginatorMap
	}

	return nil
}

//创建记录
func (m Inspect) SaveInspect() error {
	o := orm.NewOrm()
	o.Update(&m)

	return nil
}

func AddlistInspect(data string) (some_err []interface{}) {
	o := orm.NewOrm()
	o.Begin()
	var i []Inspect_add
	json.Unmarshal([]byte(data), &i)
	for key, val := range i {
		var v Inspect
		s := models.Student{Id: val.StudentId}
		o.Read(&s)
		tm2, _ := time.Parse("2006-01-02", s.Birthday)
		timestamp := time.Now().Unix()
		ptime := tm2.Unix()
		yue := float64((timestamp - ptime) / int64(30*24*3600))
		fmt.Println(yue)
		types, _ := CompareHeight(int(s.Sex), yue, val.Height)
		v.AbnormalHeight = types
		weight, _ := CompareWeight(int(s.Sex), yue, val.Weight)
		v.AbnormalWeight = weight
		v.StudentId = val.StudentId //学生ID
		v.Weight = val.Weight       //体重
		v.Height = val.Height       //身高
		v.Types = val.Types
		v.Date = time.Now().Format("2006-01-02")
		v.Content = val.Content
		v.Evaluate = val.Evaluate
		v.ClassId = val.ClassId               //班级ID
		v.KindergartenId = val.KindergartenId //幼儿园ID
		v.BodyId = val.BodyId                 //主题ID
		tmp_i := v
		if create, id, err := o.ReadOrCreate(&v, "StudentId", "ClassId", "KindergartenId", "BodyId"); err != nil {
			some_err = append(some_err, err)
		} else {
			if !create {
				tmp_i.Id = int(id)
				o.Update(&tmp_i)
			}
			i[key].Info.StudentId = val.StudentId
			i[key].Info.InspectId = int(id)
			temp_c1 := Abnormal(i[key].Info, v.BodyId, yue)
			temp_c := temp_c1
			if created, ide, err := o.ReadOrCreate(&temp_c, "StudentId", "InspectId"); err != nil {
				some_err = append(some_err, err)
			} else {
				fmt.Println(temp_c.Column4)
				temp_c1.Id = temp_c.Id
				if temp_c1.Column1 == "" {
					temp_c1.Column1 = temp_c.Column1
					temp_c1.Abnormal1 = temp_c.Abnormal1
				}
				if temp_c1.Column2 == "" {
					temp_c1.Column2 = temp_c.Column2
					temp_c1.Abnormal2 = temp_c.Abnormal2
				}
				if temp_c1.Column3 == "" {
					temp_c1.Column3 = temp_c.Column3
					temp_c1.Abnormal3 = temp_c.Abnormal3
				}
				if temp_c1.Column4 == "" {
					temp_c1.Column4 = temp_c.Column4
					temp_c1.Abnormal4 = temp_c.Abnormal4
				}
				if temp_c1.Column5 == "" {
					temp_c1.Column5 = temp_c.Column5
					temp_c1.Abnormal5 = temp_c.Abnormal5
				}
				if temp_c1.Column6 == "" {
					temp_c1.Column6 = temp_c.Column5
					temp_c1.Abnormal6 = temp_c.Abnormal5
				}
				if temp_c1.Column7 == "" {
					temp_c1.Column7 = temp_c.Column5
					temp_c1.Abnormal7 = temp_c.Abnormal5
				}
				if temp_c1.Column8 == "" {
					temp_c1.Column8 = temp_c.Column5
					temp_c1.Abnormal8 = temp_c.Abnormal5
				}
				if created == false {
					temp_c.Id = int(ide)
					_, err = o.Update(&temp_c1)

				}
			}
		}
	}

	b := Body{Id: i[0].BodyId}
	if err := o.Read(&b); err == nil {
		var num Num
		sql := "select count(a.id) as num from healthy_inspect a where a.body_id = " + strconv.Itoa(i[0].BodyId)
		o.Raw(sql).QueryRow(&num)
		b.Actual = num.Num
		b.Rate = int(math.Ceil(float64(num.Num) / float64(b.Total) * 100.0))
		o.Update(&b)
	} else {
		some_err = append(some_err, err)
	}
	var c Class
	if err := o.QueryTable("healthy_class").Filter("class_id", i[0].ClassId).Filter("body_id", i[0].BodyId).One(&c); err == nil {
		var num Num
		sql := "select count(a.id) as num from healthy_inspect a where a.body_id = " + strconv.Itoa(i[0].BodyId) + " and class_id=" + strconv.Itoa(i[0].ClassId)
		o.Raw(sql).QueryRow(&num)
		c.ClassActual = num.Num
		c.ClassRate = int(math.Ceil(float64(num.Num) / float64(c.ClassTotal) * 100.0))
		o.Update(&c)
	} else {
		some_err = append(some_err, err)
	}

	if len(some_err) > 0 {
		o.Rollback()
	} else {
		o.Commit()
		return nil
	}

	return some_err
}

func Abnormal(info Column, body_id int, yue float64) (ml Column) {
	var minfo map[string]interface{}
	mjson, _ := json.Marshal(info)
	json.Unmarshal(mjson, &minfo)
	o := orm.NewOrm()
	body := Body{Id: body_id}
	o.Read(&body)
	age := int(math.Ceil(float64(yue) / 12.0))
	c_list := strings.Split(body.Project, ",")

	for _, val := range c_list {
		c_one := strings.Split(val, ":")
		abnormal := strings.Replace(c_one[0], "column", "abnormal", 6)

		if c_one[1] == "左眼" || c_one[1] == "右眼" {
			if minfo[c_one[0]] != "" {
				v2, _ := strconv.ParseFloat(minfo[c_one[0]].(string), 64)
				if (age <= 4 && v2 >= 0.6) || (age >= 5 && v2 > 0.8) {
					minfo[abnormal] = "正常"
				} else {
					minfo[abnormal] = "异常"
				}
			}

		} else if c_one[1] == "血小板" {
			if minfo[c_one[0]] != "" {
				v2, _ := strconv.Atoi(minfo[c_one[0]].(string))
				if v2 >= 90 && v2 < 110 {
					minfo[abnormal] = "轻度贫血"
				} else if v2 >= 60 && v2 < 90 {
					minfo[abnormal] = "中度贫血"
				} else if v2 < 60 {
					minfo[abnormal] = "重度贫血"
				} else {
					minfo[abnormal] = "正常"
				}
			}
		}
	}
	mljson, _ := json.Marshal(minfo)
	json.Unmarshal(mljson, &ml)
	return ml
}

func (f *Inspect) Baby(baby_id int) (Page, error) {
	o := orm.NewOrm()
	//var con []interface{}
	where := " (height != 0 Or weight !=0) "
	where += " AND (types = 5 Or types = 4) "
	var sxWords []orm.Params

	if baby_id != 0 {
		var student_id int
		var student models.Student
		err := o.QueryTable("student").Filter("baby_id", baby_id).One(&student)
		if err != nil {
			student_id = 0
		} else {
			student_id = student.Id
		}

		where += " AND student_id = " + strconv.Itoa(student_id)

	}
	qb, _ := orm.NewQueryBuilder("mysql")
	sql := qb.Select("date,height,weight").From(f.TableName()).Where(where).Limit(6).Offset(0).String()

	if _, err := o.Raw(sql).Values(&sxWords); err == nil {

		return Page{0, 0, 0, sxWords}, nil
	}

	return Page{}, nil
}

//获取宝宝健康详情
func (f *Inspect) Situation(baby_id int) (map[string]interface{}, error) {
	o := orm.NewOrm()
	var con []interface{}
	where := "1=1 "
	if baby_id > 0 {
		where += "AND body_id  = " + strconv.Itoa(baby_id)
	}
	var Situation []orm.Params

	qb, _ := orm.NewQueryBuilder("mysql")
	sql := qb.Select("*").From(f.TableName()).Where(where).OrderBy("id").Desc().Limit(1).Offset(0).String()

	if _, err := o.Raw(sql, con).Values(&Situation); err == nil {
		return Situation[0], nil
	}

	return nil, nil
}

//健康异常档案
func (f *Inspect) Abnormals(types, page, perPage, kindergarten_id, class_id int, date, search string) (Page, error) {
	o := orm.NewOrm()
	var con []interface{}
	where := "1 "

	if kindergarten_id != 0 {
		where += "AND healthy_inspect.kindergarten_id = ? "
		con = append(con, kindergarten_id)
	}
	if search != "" {
		where += "AND ( student.name like ? Or teacher.name like ? Or healthy_inspect.abnormal like ? ) "
		con = append(con, "%"+search+"%")
		con = append(con, "%"+search+"%")
		con = append(con, "%"+search+"%")
	}
	if date == "" {
		day_time := time.Now().Format("2006-01-02")
		where += " AND left(healthy_inspect.date,10) = '" + day_time + "'"
	} else {
		where += " AND left(healthy_inspect.date,10) = '" + date + "'"
	}

	if types != 0 {
		where += " AND healthy_inspect.types = ? "
		con = append(con, types)
	}

	if class_id != 0 {
		where += " AND healthy_inspect.class_id = ? "
		con = append(con, class_id)
	}

	qb, _ := orm.NewQueryBuilder("mysql")
	sql := qb.Select("count(*)").From(f.TableName()).
		LeftJoin("student").On("healthy_inspect.student_id = student.student_id").
		LeftJoin("teacher").On("healthy_inspect.teacher_id = teacher.teacher_id").
		Where(where).String()

	var total int
	err := o.Raw(sql, con).QueryRow(&total)
	if err == nil {
		var sxWords []orm.Params

		limit := 10
		if perPage != 0 {
			limit = perPage
		}
		if page <= 0 {
			page = 1
		}
		start := (page - 1) * limit
		qb, _ := orm.NewQueryBuilder("mysql")
		sql := qb.Select("healthy_inspect.*," +
			"healthy_column.abnormal1," +
			"healthy_column.abnormal2," +
			"healthy_column.abnormal3," +
			"healthy_column.abnormal4," +
			"healthy_column.abnormal5," +
			"healthy_column.abnormal6," +
			"organizational.name as className,s" +
			"tudent.name as student_name," +
			"student.avatar," +
			"teacher.name as teacher_name").From(f.TableName()).
			LeftJoin("student").On("healthy_inspect.student_id = student.student_id").
			LeftJoin("teacher").On("healthy_inspect.teacher_id = teacher.teacher_id").
			LeftJoin("organizational").On("healthy_inspect.class_id = organizational.id").
			LeftJoin("healthy_column").On("healthy_column.inspect_id = healthy_inspect.id").
			Where(where).
			OrderBy("healthy_inspect.id").Desc().
			Limit(limit).Offset(start).String()

		if _, err := o.Raw(sql, con).Values(&sxWords); err == nil {

			pageNum := int(math.Ceil(float64(total) / float64(limit)))
			return Page{pageNum, limit, total, sxWords}, nil
		}
	}

	return Page{}, nil
}

//体检详情
func (f *Inspect) Project(page, perPage, kindergarten_id, class_id, body_id, baby_id int, column string) (Page, error) {
	o := orm.NewOrm()
	var con []interface{}
	where := "1 "
	where += "AND healthy_inspect.types = 5 "

	if class_id > 0 {
		where += " AND healthy_inspect.class_id = " + strconv.Itoa(class_id)
	}
	if baby_id > 0 {
		var student_id int
		var student models.Student
		err := o.QueryTable("student").Filter("baby_id", baby_id).One(&student)
		if err != nil {
			student_id = 0
		} else {
			student_id = student.Id
		}
		where += " AND healthy_inspect.student_id = " + strconv.Itoa(student_id)
	}
	if body_id > 0 {
		where += " AND healthy_inspect.body_id = " + strconv.Itoa(body_id)
	}
	if kindergarten_id > 0 {
		where += " AND healthy_inspect.kindergarten_id = " + strconv.Itoa(kindergarten_id)
	}
	qb, _ := orm.NewQueryBuilder("mysql")
	sql := qb.Select("count(*)").From(f.TableName()).Where(where).String()
	if column != "" {
		if column == "weight" {
			where += " AND healthy_inspect.weight != 0"
		} else {
			where += " AND healthy_column." + column + " != '' "
		}
	}
	var total int
	err := o.Raw(sql, con).QueryRow(&total)
	if err == nil {
		var sxWords []orm.Params

		limit := 10
		if perPage != 0 {
			limit = perPage
		}
		if page <= 0 {
			page = 1
		}
		start := (page - 1) * limit
		qb, _ := orm.NewQueryBuilder("mysql")
		sql := qb.Select("healthy_inspect.*,student.name as student_name,student.student_id as studentId,student.avatar,healthy_column.*").From(f.TableName()).
			LeftJoin("student").On("healthy_inspect.student_id = student.student_id").
			LeftJoin("healthy_column").On("healthy_inspect.id = healthy_column.inspect_id").
			Where(where).
			Limit(limit).Offset(start).String()

		if _, err := o.Raw(sql, con).Values(&sxWords); err == nil {

			pageNum := int(math.Ceil(float64(total) / float64(limit)))
			return Page{pageNum, limit, total, sxWords}, nil
		}
	}

	return Page{}, nil
}

//体检详情
func (f *Inspect) ProjectNew(page, perPage, kindergarten_id, class_id, body_id, baby_id int, column string) (Page, error) {
	o := orm.NewOrm()
	var con []interface{}
	where := "1 "
	where += " AND healthy_inspect.types = 5 "

	fmt.Println(where)

	if class_id > 0 {
		where += " AND healthy_inspect.class_id = " + strconv.Itoa(class_id)
	}
	if baby_id > 0 {
		var student_id int
		var student models.Student
		err := o.QueryTable("student").Filter("baby_id", baby_id).One(&student)
		if err != nil {
			student_id = 0
		} else {
			student_id = student.Id
		}
		where += " AND healthy_inspect.student_id = " + strconv.Itoa(student_id)
	}
	if body_id > 0 {
		where += " AND healthy_inspect.body_id = " + strconv.Itoa(body_id)
	}
	if kindergarten_id > 0 {
		where += " AND healthy_inspect.kindergarten_id = " + strconv.Itoa(kindergarten_id)
	}
	qb, _ := orm.NewQueryBuilder("mysql")
	sql := qb.Select("count(*)").From(f.TableName()).Where(where).String()

	if column != "" {
		if column == "weight" {
			where += " AND healthy_inspect.weight = 0 "
		} else {
			where += " AND (healthy_column." + column + " is null Or healthy_column." + column + " = '' )"
		}
	}

	var total int
	err := o.Raw(sql, con).QueryRow(&total)
	if err == nil {
		var sxWords []orm.Params

		limit := 10
		if perPage != 0 {
			limit = perPage
		}
		if page <= 0 {
			page = 1
		}
		start := (page - 1) * limit
		qb, _ := orm.NewQueryBuilder("mysql")
		sql := qb.Select("healthy_inspect.*,student.name as student_name,student.student_id as studentId,student.avatar,healthy_column.*").From(f.TableName()).
			LeftJoin("student").On("healthy_inspect.student_id = student.student_id").
			LeftJoin("healthy_column").On("healthy_inspect.id = healthy_column.inspect_id").
			Where(where).
			Limit(limit).Offset(start).String()

		if _, err := o.Raw(sql, con).Values(&sxWords); err == nil {

			pageNum := int(math.Ceil(float64(total) / float64(limit)))
			return Page{pageNum, limit, total, sxWords}, nil
		}
	}

	return Page{}, nil
}

//宝宝健康指数
func (f *Inspect) Personal(baby_id int) ([]orm.Params, error) {
	o := orm.NewOrm()
	var con []interface{}
	where := "1 "
	where += "AND ( healthy_inspect.types = 1 Or healthy_inspect.types = 2 Or healthy_inspect.types = 3 ) "
	day_time := time.Now().Format("2006-01-02")
	where += " AND left(date,10) = '" + day_time + "'"
	wheres := " left(date,10) = '" + day_time + "'"
	wheres += " AND abnormal is not null "
	wheres += " AND ( types = 1 Or types = 2 Or types = 3 ) "
	if baby_id != 0 {
		var student_id int
		var student models.Student
		err := o.QueryTable("student").Filter("baby_id", baby_id).One(&student)
		if err != nil {
			student_id = 0
		} else {
			student_id = student.Id
		}

		where += " AND healthy_inspect.student_id = ? "
		con = append(con, student_id)
		wheres += " AND student_id = " + strconv.Itoa(student_id)
	}
	qb, _ := orm.NewQueryBuilder("mysql")
	sql := qb.Select("count(*)").From(f.TableName()).Where(wheres).String()

	var total int
	err := o.Raw(sql).QueryRow(&total)

	if err == nil {
		var sxWords []orm.Params

		qb, _ := orm.NewQueryBuilder("mysql")
		sql := qb.Select("healthy_inspect.*,student.name as student_name,student.avatar").From(f.TableName()).
			LeftJoin("student").On("healthy_inspect.student_id = student.student_id").
			Where(where).
			OrderBy("id").Desc().
			Limit(1).Offset(0).String()
		if _, err := o.Raw(sql, con).Values(&sxWords); err == nil {
			if len(sxWords) != 0 {
				sxWords[0]["index"] = 100 - total*20

				return sxWords, nil
			}
		}
	}

	return nil, nil
}

//体重统计
func (f *Inspect) Weights(kindergarten_id int, date string) ([]orm.Params, error) {
	o := orm.NewOrm()
	var counts []orm.Params
	wheres := "1"
	wheres += " AND (healthy_inspect.types = 4 Or healthy_inspect.types = 5)"
	if date != "" {
		wheres += " AND left( healthy_body.test_time,4) = '" + date + "'"
	}
	wheres += " group by body_id"

	wheres += " order by healthy_body.id desc "
	wheres += " limit 6 "
	_, err := o.Raw("SELECT healthy_body.test_time as date, sum(if(abnormal_weight = '肥胖',1,0)) as fat,sum(if(abnormal_weight = '超重',1,0)) as heavy ,sum(if(abnormal_weight = '瘦小',1,0)) as thin FROM healthy_inspect left join healthy_body on healthy_inspect.body_id = healthy_body.id  WHERE " + wheres).Values(&counts)
	if err == nil {
		return counts, nil
	}
	return nil, nil
}

//身高统计
func (f *Inspect) Heights(kindergarten_id int, date string) ([]orm.Params, error) {
	o := orm.NewOrm()
	var counts []orm.Params
	wheres := "1"
	wheres += " AND (healthy_inspect.types = 4 Or healthy_inspect.types = 5)"
	if date != "" {
		wheres += " AND left( healthy_body.test_time,4) = '" + date + "'"
	}
	wheres += " group by body_id"

	wheres += " order by healthy_body.id desc "
	wheres += " limit 6 "
	_, err := o.Raw("SELECT healthy_body.test_time as date, sum(if(abnormal_height = '超高',1,0)) as fat,sum(if(abnormal_height = '偏矮',1,0)) as heavy ,sum(if(abnormal_height = '矮小',1,0)) as thin FROM healthy_inspect left join healthy_body on healthy_inspect.body_id = healthy_body.id  WHERE " + wheres).Values(&counts)
	if err == nil {
		return counts, nil
	}
	return nil, nil
}

//全园统计
func (f *Inspect) Country(kindergarten_id int) ([]orm.Params, error) {
	o := orm.NewOrm()
	var counts []orm.Params
	wheres := " 1 "
	wheres += " AND kindergarten_id = " + strconv.Itoa(kindergarten_id)

	wheres += " order by id desc "
	wheres += " limit 6 "

	_, err := o.Raw("SELECT test_time,(total - actual) as cha,actual,rate FROM healthy_body WHERE" + wheres).Values(&counts)

	if err == nil {
		return counts, nil
	}
	return nil, nil
}

//后台体检项目
func (f *Inspect) Projects(page, perPage, kindergarten_id, class_id, body_id, baby_id int, search string) (Page, error) {
	o := orm.NewOrm()
	var con []interface{}
	where := "1 "
	where += "AND healthy_inspect.types = 5 "

	if class_id > 0 {
		where += " AND healthy_inspect.class_id = " + strconv.Itoa(class_id)
	}

	if baby_id > 0 {
		var student_id int
		var student models.Student
		err := o.QueryTable("student").Filter("baby_id", baby_id).One(&student)
		if err != nil {
			student_id = 0
		} else {
			student_id = student.Id
		}
		where += " AND healthy_inspect.student_id = " + strconv.Itoa(student_id)
	}
	if body_id > 0 {
		where += " AND healthy_inspect.body_id = " + strconv.Itoa(body_id)
	}
	if kindergarten_id > 0 {
		where += " AND healthy_inspect.kindergarten_id = " + strconv.Itoa(kindergarten_id)
	}
	if search != "" {
		where += " AND student.name like " + "'%" + search + "%'"
	}
	qb, _ := orm.NewQueryBuilder("mysql")
	//sql := qb.Select("").From(f.TableName()).Where(where).String()
	sql := qb.Select("count(*)").From(f.TableName()).
		LeftJoin("student").On("healthy_inspect.student_id = student.student_id").
		LeftJoin("healthy_column").On("healthy_inspect.id = healthy_column.inspect_id").
		Where(where).String()
	var total int
	err := o.Raw(sql, con).QueryRow(&total)
	if err == nil {
		var sxWords []orm.Params

		limit := 10
		if perPage != 0 {
			limit = perPage
		}
		if page <= 0 {
			page = 1
		}
		start := (page - 1) * limit
		qb, _ := orm.NewQueryBuilder("mysql")
		sql := qb.Select("healthy_inspect.*,healthy_inspect.id as hid,student.sex,student.age,student.name as student_name,student.student_id as studentId,student.avatar,healthy_column.*").From(f.TableName()).
			LeftJoin("student").On("healthy_inspect.student_id = student.student_id").
			LeftJoin("healthy_column").On("healthy_inspect.id = healthy_column.inspect_id").
			Where(where).
			Limit(limit).Offset(start).String()

		if _, err := o.Raw(sql, con).Values(&sxWords); err == nil {

			pageNum := int(math.Ceil(float64(total) / float64(limit)))
			return Page{pageNum, limit, total, sxWords}, nil
		}
	}

	return Page{}, nil
}

//删除
func DeleteStudent(id int) map[string]interface{} {
	o := orm.NewOrm()
	v := Inspect{Id: id}
	err := o.Begin()
	if err = o.Read(&v); err == nil {
		_, err = o.Raw("UPDATE healthy_body SET actual = actual - 1 where id = ? ", v.BodyId).Exec()
		if err != nil {
			err = o.Rollback()
			return nil
		}
		_, err = o.Raw("UPDATE healthy_class SET class_actual = class_actual - 1 where class_id = ? and body_id = ? ", v.ClassId, v.BodyId).Exec()
		if err != nil {
			err = o.Rollback()
			return nil
		}
		if num, err := o.Delete(&Inspect{Id: id}); err == nil {
			paginatorMap := make(map[string]interface{})
			paginatorMap["data"] = num
			o.Commit()
			return paginatorMap
		} else {
			err = o.Rollback()
			return nil
		}
	}
	return nil
}

//添加备注
func (f *Inspect) Contents(content string) error {
	o := orm.NewOrm()
	if err := o.Read(f); err == nil {
		f.Content = content
		if _, err := o.Update(f, "Content"); err != nil {
			return err
		}
	} else {
		return err
	}

	return nil
}

//SaaS统计
func Countss(kindergarten_id int) map[string]interface{} {
	o := orm.NewOrm()
	counts := make(map[string]interface{})
	numbers := make(map[string]int64)
	numbers["Monday"] = 0
	numbers["Tuesday"] = 1
	numbers["Wednesday"] = 2
	numbers["Thursday"] = 3
	numbers["Friday"] = 4
	numbers["Saturday"] = 5
	numbers["Sunday"] = 6
	//每天统计
	where := " kindergarten_id = " + strconv.Itoa(kindergarten_id)
	day_time := time.Now().Format("2006-01-02")
	where += " AND left(date,10) = '" + day_time + "'"
	where += " AND (types = 1 Or types = 2 Or types = 3)"
	where += " AND abnormal != '' "
	type Counts struct {
		Num int
	}
	var count []Counts
	_, err := o.Raw("SELECT count(id) as num FROM healthy_inspect where" + where).QueryRows(&count)
	if err == nil {
		counts["day"] = count[0].Num
	}
	//每月统计
	month_time := time.Now().Format("2006-01")
	where2 := " kindergarten_id = " + strconv.Itoa(kindergarten_id)
	where2 += " AND (types = 1 Or types = 2 Or types = 3)"
	where2 += " AND left(healthy_inspect.date,7) = '" + month_time + "'"
	_, err = o.Raw("SELECT count(healthy_inspect.id) as num FROM healthy_inspect where" + where2).QueryRows(&count)
	if err == nil {
		counts["month"] = count[0].Num
	}
	//每周统计
	t := time.Now()
	date := time.Now().Unix() - 24*60*60*(numbers[t.Weekday().String()])
	tm := time.Unix(date, 0)
	startTime := tm.Format("2006-01-02 00:00:00")

	where3 := " kindergarten_id = " + strconv.Itoa(kindergarten_id)
	where3 += " AND healthy_inspect.created_at >= '" + startTime + "'"
	_, err = o.Raw("SELECT count(healthy_inspect.id) as num FROM healthy_inspect where" + where3).QueryRows(&count)
	if err == nil {
		counts["week"] = count[0].Num
	}
	//全部统计
	where4 := " kindergarten_id = " + strconv.Itoa(kindergarten_id)
	where4 += " AND types = 1 Or types = 2 Or types = 3"
	all, err := o.Raw("SELECT count(id) as num FROM healthy_inspect where" + where4).QueryRows(&count)
	if err == nil {
		fmt.Println(all)
		counts["all"] = count[0].Num
	}

	return counts
}

//宝宝健康指数
func (f *Inspect) PersonalInfo(baby_id int) (Page, error) {
	o := orm.NewOrm()
	var con []interface{}
	where := "1 "
	where += "AND ( healthy_inspect.types = 1 Or healthy_inspect.types = 2 Or healthy_inspect.types = 3 ) "
	day_time := time.Now().Format("2006-01-02")
	where += " AND left(date,10) = '" + day_time + "'"
	wheres := " left(date,10) = '" + day_time + "'"
	wheres += " AND abnormal is not null "
	wheres += " AND ( types = 1 Or types = 2 Or types = 3 ) "
	var student_id int
	if baby_id != 0 {
		fmt.Println(baby_id)
		var student models.Student
		err := o.QueryTable("student").Filter("baby_id", baby_id).One(&student)
		if err != nil {
			student_id = 0
		} else {
			student_id = student.Id
		}

		fmt.Println(student_id)
		where += " AND healthy_inspect.student_id = ? "
		con = append(con, student_id)
		wheres += " AND student_id = " + strconv.Itoa(student_id)
	}
	qb, _ := orm.NewQueryBuilder("mysql")
	sql := qb.Select("count(*)").From(f.TableName()).Where(wheres).String()

	var total int
	err := o.Raw(sql).QueryRow(&total)

	if err == nil {
		var sxWords []orm.Params

		qb, _ := orm.NewQueryBuilder("mysql")
		sql := qb.Select("healthy_inspect.*,student.name as student_name,student.age,student.avatar").From(f.TableName()).
			LeftJoin("student").On("healthy_inspect.student_id = student.student_id").
			Where(where).
			OrderBy("id").Desc().
			Limit(2).Offset(0).String()
		if _, err := o.Raw(sql, con).Values(&sxWords); err == nil {
			if len(sxWords) != 0 {
				qb, _ := orm.NewQueryBuilder("mysql")
				wheres2 := " types = 5 "
				wheres2 += " AND student_id = " + strconv.Itoa(student_id)
				sql1 := qb.Select("*").From(f.TableName()).Where(wheres2).OrderBy("id").Desc().Limit(1).Offset(0).String()
				var Heights []orm.Params
				if _, err := o.Raw(sql1).Values(&Heights); err == nil {
					sxWords[0]["height"] = Heights[0]["height"]
					sxWords[0]["weight"] = Heights[0]["weight"]
				}
				sxWords[0]["index"] = 100 - total*20
				return Page{0, 0, total, sxWords}, nil
			} else {
				qb, _ := orm.NewQueryBuilder("mysql")
				wheres2 := " healthy_inspect.types = 5 "
				wheres2 += " AND healthy_inspect.student_id = " + strconv.Itoa(student_id)
				sql1 := qb.Select("healthy_inspect.*,student.name as student_name,student.age,student.avatar").From(f.TableName()).
					LeftJoin("student").On("healthy_inspect.student_id = student.student_id").
					Where(wheres2).
					OrderBy("id").Desc().Limit(1).Offset(0).String()
				var Heights []orm.Params
				if _, err := o.Raw(sql1).Values(&Heights); err == nil {
					Heights[0]["index"] = 100
				}
				return Page{0, 0, total, Heights}, nil
			}
		}
	}

	return Page{}, nil
}

//主题列表
func (f *Inspect) Body(page, perPage, kindergarten_id, baby_id int) (Page, error) {
	o := orm.NewOrm()
	var con []interface{}
	where := "1 "
	where += "AND healthy_inspect.types = 5 "
	if baby_id > 0 {
		var student_id int
		var student models.Student
		err := o.QueryTable("student").Filter("baby_id", baby_id).One(&student)
		if err != nil {
			student_id = 0
		} else {
			student_id = student.Id
		}
		where += " AND healthy_inspect.student_id = " + strconv.Itoa(student_id)
	}
	if kindergarten_id > 0 {
		where += " AND healthy_inspect.kindergarten_id = " + strconv.Itoa(kindergarten_id)
	}
	qb, _ := orm.NewQueryBuilder("mysql")
	sql := qb.Select("count(*)").From(f.TableName()).Where(where).String()
	var total int
	err := o.Raw(sql, con).QueryRow(&total)
	if err == nil {
		var sxWords []orm.Params

		limit := 10
		if perPage != 0 {
			limit = perPage
		}
		if page <= 0 {
			page = 1
		}
		start := (page - 1) * limit
		where += " AND healthy_body.status = 1 "
		qb, _ := orm.NewQueryBuilder("mysql")
		sql := qb.Select("healthy_body.*").From(f.TableName()).
			LeftJoin("healthy_body").On("healthy_inspect.body_id = healthy_body.id").
			Where(where).
			Limit(limit).Offset(start).String()

		if _, err := o.Raw(sql, con).Values(&sxWords); err == nil {

			pageNum := int(math.Ceil(float64(total) / float64(limit)))
			return Page{pageNum, limit, total, sxWords}, nil
		}
	}

	return Page{}, nil
}
