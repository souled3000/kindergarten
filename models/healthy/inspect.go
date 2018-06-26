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
	Height         float64		`json:"height" orm:"column(height)" description:"身高"`
	Content        string 		`json:"content" orm:"column(content);size(100)" description:"备注"`
	Evaluate       int			`json:"evaluate" orm:"column(evaluate)" description:"评价"`
	ClassId        int			`json:"class_id" orm:"column(class_id)" description:"班级ID"`
	KindergartenId int			`json:"kindergarten_id" orm:"column(kindergarten_id)" description:"幼儿园ID"`
	NoteTaker      int 			`json:"note_taker" orm:"column(note_taker)" description:"记录人"`
	Types          int			`json:"types" orm:"column(types)" description:"类型"`
	Handel         string 		`json:"handel" orm:"column(handel);size(100)" description:"处理方式"`
	Url            string 		`json:"url" orm:"column(url);size(1800)" description:"照片"`
	Infect         int			`json:"infect" orm:"column(infect)" description:"是否传染"`
	DrugId		   int			`json:"drug_id" orm:"column(drug_id)" description:"喂药申请"`
	Abnormal       string 		`json:"abnormal" orm:"column(abnormal)" description:"是否传染"`
	CreatedAt      time.Time 	`json:"created_at" orm:"auto_now_add" description:"创建时间"`
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
func (f *Inspect) GetAll(page, perPage, kindergarten_id, class_id, types int) (Page, error) {
	var con []interface{}
	where := "1 "

	if kindergarten_id != 0 {
		where += "AND healthy_inspect.kindergarten_id = ? "
		con = append(con, kindergarten_id)
	}
	if class_id != 0 {
		where += "AND healthy_inspect.class_id = ? "
		con = append(con, class_id)
	}
	if types != 0 {
		where += "AND healthy_inspect.types = ? "
		con = append(con, types)
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
		sql := qb.Select("healthy_inspect.*,student.name,organizational.name as class").From(f.TableName()).
			LeftJoin("student").On("healthy_inspect.student_id = student.student_id").
			LeftJoin("organizational").On("healthy_inspect.class_id = organizational.id").
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
