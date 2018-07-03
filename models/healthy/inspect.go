package healthy

import (
	"github.com/astaxie/beego/orm"
	"time"
	"math"
	"kindergarten-service-go/models"
	"fmt"
	"encoding/json"
	"strconv"
)

type Inspect struct {
	Id             int			`json:"id" orm:"column(id);auto" description:"编号"`
	StudentId      int			`json:"student_id" orm:"column(student_id)" description:"学生ID"`
	Weight         float64		`json:"weight" orm:"column(weight)" description:"体重"`
	AbnormalWeight string		`json:"abnormal_weight" orm:"column(abnormal_weight)"`
	ClassName	   string		`json:"class_name" orm:"column" description:"班级名称"`
	Height         float64		`json:"height" orm:"column(height)" description:"身高"`
	AbnormalHeight string		`json:"abnormal_height" orm:"column(abnormal_height)"`
	Content        string 		`json:"content" orm:"column(content);size(100)" description:"备注"`
	Evaluate       int			`json:"evaluate" orm:"column(evaluate)" description:"评价"`
	ClassId        int			`json:"class_id" orm:"column(class_id)" description:"班级ID"`
	KindergartenId int			`json:"kindergarten_id" orm:"column(kindergarten_id)" description:"幼儿园ID"`
	TeacherId      int 			`json:"teacher_id" orm:"column(teacher_id)" description:"记录人ID"`
	Types          int			`json:"types" orm:"column(types)" description:"类型"`
	Handel         string 		`json:"handel" orm:"column(handel);size(100)" description:"处理方式"`
	Url            string 		`json:"url" orm:"column(url);size(1800)" description:"照片"`
	Infect         int			`json:"infect" orm:"column(infect)" description:"是否传染"`
	DrugId		   int			`json:"drug_id" orm:"column(drug_id)" description:"喂药申请"`
	Abnormal       string 		`json:"abnormal" orm:"column(abnormal)" description:"是否传染"`
	BodyId		   int			`json:"body_id" orm:"column(body_id)"`
	Date		   string		`json:"date" orm:"column(date)" description:"参检时间"`
	CreatedAt      time.Time 	`json:"created_at" orm:"auto_now_add;auto_now" description:"创建时间"`
}

type Inspect_add struct {
	Id             int			`json:"id" orm:"column(id);auto" description:"编号"`
	StudentId      int			`json:"student_id" orm:"column(student_id)" description:"学生ID"`
	Weight         float64		`json:"weight" orm:"column(weight)" description:"体重"`
	ClassName	   string		`json:"class_name" orm:"column" description:"班级名称"`
	Height         float64		`json:"height" orm:"column(height)" description:"身高"`
	Content        string 		`json:"content" orm:"column(content);size(100)" description:"备注"`
	Evaluate       int			`json:"evaluate" orm:"column(evaluate)" description:"评价"`
	ClassId        int			`json:"class_id" orm:"column(class_id)" description:"班级ID"`
	KindergartenId int			`json:"kindergarten_id" orm:"column(kindergarten_id)" description:"幼儿园ID"`
	TeacherId      int 			`json:"teacher_id" orm:"column(teacher_id)" description:"记录人ID"`
	Types          int			`json:"types" orm:"column(types)" description:"类型"`
	Handel         string 		`json:"handel" orm:"column(handel);size(100)" description:"处理方式"`
	Url            string 		`json:"url" orm:"column(url);size(1800)" description:"照片"`
	Infect         int			`json:"infect" orm:"column(infect)" description:"是否传染"`
	DrugId		   int			`json:"drug_id" orm:"column(drug_id)" description:"喂药申请"`
	Abnormal       string 		`json:"abnormal" orm:"column(abnormal)" description:"是否传染"`
	Date		   string		`json:"date" orm:"column(date)" description:"参检时间"`
	CreatedAt      time.Time 	`json:"created_at" orm:"auto_now_add;auto_now" description:"创建时间"`
	BodyId         int			`json:"body_id" orm:"column(body_id)" description:"班级ID"`
	Info		   Column		`json:"info" orm:""`
}

func (t *Inspect) TableName() string {
	return "healthy_inspect"
}

func init() {
	orm.RegisterModel(new(Inspect))
}

type Page struct {
	PageNum int `json:"page_num"`
	PerPage int `json:"per_page"`
	Total int `json:"total"`
	Data []orm.Params `json:"data"`
}

//创建记录
func (m Inspect) Save() error {
	tmp := m
	o := orm.NewOrm()

	where := " where 1=1 "
	if tmp.StudentId > 0{
		where += " and student_id = "+strconv.Itoa(tmp.StudentId)
	}

	if tmp.Date != "" {
		where += " and left(created_at,10) = '"+tmp.Date+"'"
	}else {
		day_time := time.Now().Format("2006-01-02")
		where += " and left(created_at,10) = '"+ day_time +"'"
	}

	var drug []Drug
	_, err := o.Raw("SELECT id FROM healthy_drug "+where).QueryRows(&drug)
	if err == nil {
		m.DrugId = drug[0].Id
	}
	o.Insert(&m);

	return nil
}

//晨、午、晚列表
func (f *Inspect) GetAll(page, perPage, kindergarten_id, class_id, types, role, baby_id int, date string) (Page, error) {
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

	if date != "" {
		where += "AND sx_user.date like ? "
		con = append(con, "%"+date+"%")
	}

	if baby_id != 0 {
		var student_id int
		var student models.Student
		err := o.QueryTable("student").Filter("baby_id", baby_id).One(&student)
		if err != nil{
			student_id = 0
		}else {
			student_id = student.Id
		}

		fmt.Print(student)
		where += "AND healthy_inspect.student_id = ? "
		con = append(con, student_id)
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
		qb, _ := orm.NewQueryBuilder("mysql")
		sql := qb.Select("healthy_inspect.*,student.name as student_name,student.avatar,teacher.name as teacher_name").From(f.TableName()).
			LeftJoin("student").On("healthy_inspect.student_id = student.student_id").
			LeftJoin("teacher").On("healthy_inspect.teacher_id = teacher.teacher_id").
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
func Counts(kindergarten_id int) map[string]interface{}  {
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
	day_time := time.Now().Format("2006-01-02")
	day, err := o.QueryTable("healthy_inspect").Filter("created_at__contains", day_time).Filter("kindergarten_id",kindergarten_id).Filter("abnormal__isnull", false).Count()
	if err == nil{
		counts["day"] = day
	}
	//实际检查人数
	day_actual, err := o.QueryTable("healthy_inspect").Filter("created_at__contains", day_time).Filter("kindergarten_id",kindergarten_id).Count()
	if err == nil{
		counts["day_actual"] = day_actual
	}
	//每月统计
	month_time := time.Now().Format("2006-01")
	month, err := o.QueryTable("healthy_inspect").Filter("created_at__contains", month_time).Filter("kindergarten_id",kindergarten_id).Filter("abnormal__isnull", false).Count()
	if err == nil{
		counts["month"] = month
	}
	//每周统计
	t := time.Now()
	date := time.Now().Unix() - 24 * 60 * 60 *(numbers[t.Weekday().String()])
	tm := time.Unix(date, 0)
	startTime := tm.Format("2006-01-02 00:00:00")
	week, err := o.QueryTable("healthy_inspect").Filter("created_at__gte",startTime).Filter("kindergarten_id",kindergarten_id).Filter("abnormal__isnull", false).Count()
	if err == nil{
		counts["week"] = week
	}
	//全部统计
	all, err := o.QueryTable("healthy_inspect").Filter("kindergarten_id",kindergarten_id).Filter("abnormal__isnull", false).Count()
	if err == nil{
		counts["all"] = all
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
	o.Update(&m);

	return nil
}

func AddlistInspect(data string) (some_err []interface{}) {
	o := orm.NewOrm()
	o.Begin()
	var i []Inspect_add
	json.Unmarshal([]byte(data),&i)
	for key,val := range i{
		var v Inspect
		s := models.Student{Id:val.StudentId}
		o.Read(&s)
		tm2, _ := time.Parse("2006-01-02", s.Birthday)
		timestamp := time.Now().Unix()
		ptime := tm2.Unix()
		yue := float64((timestamp-ptime)/int64(30*24*3600))

		fmt.Println(yue)


		types,_ := CompareHeight(int(s.Sex),yue,val.Height)
		v.AbnormalHeight = types
		weight, _:=CompareWeight(int(s.Sex),yue,val.Weight)
		v.AbnormalWeight = weight
		v.StudentId = val.StudentId
		v.Weight = val.Weight
		v.Height = val.Height
		v.Content = val.Content
		v.Evaluate = val.Evaluate
		v.ClassId = val.ClassId
		v.ClassName = val.ClassName
		v.KindergartenId = val.KindergartenId
		v.TeacherId = val.TeacherId
		v.Types = val.Types
		v.BodyId = val.BodyId
		v.Infect = val.Infect
		v.Abnormal = val.Abnormal
		v.Date = val.Date
		v.Url = val.Url
		tmp_i := v
		if create, id, err := o.ReadOrCreate(&v, "StudentId","ClassId","KindergartenId","Types","BodyId"); err != nil {
			some_err = append(some_err,err)
		} else {
			if !create {
				tmp_i.Id = int(id)
				o.Update(&tmp_i)
			}
			i[key].Info.StudentId = val.StudentId
			i[key].Info.InspectId = int(id)
			temp_c := i[key].Info
			if created, ide, err := o.ReadOrCreate(&i[key].Info, "StudentId","InspectId"); err != nil {
				some_err = append(some_err,err)
			} else {
				if !created {
					temp_c.Id = int(ide)
					o.Update(&temp_c)
				}
			}
		}
	}

	if len(some_err) > 0 {
		o.Rollback()
	} else {
		o.Commit()
		return  nil
	}

	return some_err
}

func (f *Inspect) Baby(baby_id int) (Page, error) {
	o := orm.NewOrm()
	var con []interface{}
	where := " (height != 0 Or weight !=0) "
	where += " AND (types = 5 Or types = 4) "
	var sxWords []orm.Params


	qb, _ := orm.NewQueryBuilder("mysql")
	sql := qb.Select("date,height,weight").From(f.TableName()).Where(where).Limit(5).Offset(0).String()

	if _, err := o.Raw(sql, con).Values(&sxWords); err == nil {

		return Page{0, 0, 0, sxWords}, nil
	}

	return Page{}, nil
}

//获取宝宝健康详情
func (f *Inspect) Situation(baby_id int) (map[string]interface{}, error) {
	o := orm.NewOrm()
	var con []interface{}
	where := "1=1 "
	if baby_id > 0{
		where += "AND body_id  = "+strconv.Itoa(baby_id)
	}
	var Situation []orm.Params


	qb, _ := orm.NewQueryBuilder("mysql")
	sql := qb.Select("*").From(f.TableName()).Where(where).OrderBy("id").Desc().Limit(1).Offset(0).String()

	if _, err := o.Raw(sql, con).Values(&Situation); err == nil {
		return Situation[0],nil
	}

	return nil, nil
}

//健康异常档案
func (f *Inspect) Abnormals(page, perPage, kindergarten_id, class_id int, date, search string) (Page, error) {
	o := orm.NewOrm()
	var con []interface{}
	where := "1 "

	if kindergarten_id != 0 {
		where += "AND healthy_inspect.kindergarten_id = ? "
		con = append(con, kindergarten_id)
	}

	qb, _ := orm.NewQueryBuilder("mysql")
	sql := qb.Select("count(*)").From(f.TableName()).Where(where).String()

	var total int
	err := o.Raw(sql, con).QueryRow(&total)
	if err == nil {
		if search != "" {
			where += "AND ( student.name like ? Or teacher.name like ? Or healthy_inspect.abnormal like ? ) "
			con = append(con, "%"+search+"%")
			con = append(con, "%"+search+"%")
			con = append(con, "%"+search+"%")
		}
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
			OrderBy("id").Desc().
			Limit(limit).Offset(start).String()

		if _, err := o.Raw(sql, con).Values(&sxWords); err == nil {

			pageNum := int(math.Ceil(float64(total) / float64(limit)))
			return Page{pageNum, limit, total, sxWords}, nil
		}
	}

	return Page{}, nil
}