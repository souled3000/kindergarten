package healthy

import (
	"github.com/astaxie/beego/orm"
	"time"
	"math"
)

type Inspect struct {
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
	o := orm.NewOrm()
	o.Insert(&m);

	return nil
}

//晨、午、晚列表
func (f *Inspect) GetAll(page, perPage, kindergarten_id, class_id, types, role int, date string) (Page, error) {
	var con []interface{}
	where := "1 "

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

	o := orm.NewOrm()
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
		sql := qb.Select("healthy_inspect.*,student.name as student_name,teacher.name as teacher_name").From(f.TableName()).
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