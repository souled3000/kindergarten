package healthy

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/hprose/hprose-golang/rpc"
	"kindergarten-service-go/models"
	"math"
	"time"
)

type Drug struct {
	Id             int       `json:"id" orm:"column(id);auto" description:"编号"`
	StudentId      int       `json:"student_id" orm:"column(student_id)" description:"学生ID"`
	Drug           string    `json:"drug" orm:"column(drug)" description:"药片"`
	Symptom        string    `json:"symptom" orm:"column(symptom)" description:"症状"`
	Explain        string    `json:"explain" orm:"column(explain)" description:"用量说明"`
	Url            string    `json:"url" orm:"column(url)" description:"喂药申请图片"`
	KindergartenId int       `json:"kindergarten_id" orm:"column(kindergarten_id)" description:"幼儿园ID"`
	ClassId        int       `json:"class_id" orm:"column(class_id)" description:"班级ID"`
	ClassName      string    `json:"class_name" orm:"column(class_name)" description:"班级名称"`
	UserId         int       `json:"user_id" orm:"column(user_id)"`
	Avatar         string    `json:"avatar" orm:"column(avatar)" description:"头像"`
	NoteTaker      string    `json:"note_taker" orm:"column(note_taker)" description:"喂药申请人"`
	CreatedAt      time.Time `json:"created_at" orm:"column(created_at);auto_now_add"`
}

var User *UserService

func (t *Drug) TableName() string {
	return "healthy_drug"
}

func init() {
	orm.RegisterModel(new(Drug))
}

type UserService struct {
	GetBabyCall func(userId, babyId int) (map[string]map[string]interface{}, error)
}

//申请喂药
func (m Drug) Save(baby_id int) error {
	tmp := m
	o := orm.NewOrm()
	//获取学生信息
	client := rpc.NewHTTPClient(beego.AppConfig.String("ONE_MORE_USER_SERVER"))
	client.UseService(&User)
	if uk, err := User.GetBabyCall(tmp.UserId, baby_id); err != nil {
		return err
	} else {
		avatar := uk["baby"]["avatar"].(string)
		note_taker := uk["baby"]["nickname"].(string) + "的" + uk["user_family"]["identity"].(string)
		tmp.Avatar = avatar
		tmp.NoteTaker = note_taker
	}

	var student models.Student
	err := o.QueryTable("student").Filter("baby_id", baby_id).One(&student)
	if err == nil {
		tmp.KindergartenId = student.KindergartenId
		tmp.StudentId = student.Id
		//通过学生ID获取班级ID
		var member models.OrganizationalMember
		err := o.QueryTable("organizational_member").Filter("member_id", tmp.StudentId).One(&member)
		if err == nil {
			//通过组织架构ID获取班级ID
			var organizational models.Organizational
			err := o.QueryTable("organizational").Filter("id", member.OrganizationalId).Filter("type", 2).Filter("level", 3).One(&organizational)
			if err == nil {
				tmp.ClassId = organizational.Id
				tmp.ClassName = organizational.Name
				o.Insert(&tmp)
			} else {
				return err
			}
		} else {
			return err
		}
		return nil
	} else {

		return err
	}

	return nil
}

//获取喂药记录
func (f *Drug) GetAll(page, perPage, kindergarten_id, class_id, role, types int) (Page, error) {
	var con []interface{}
	where := "1 "
	day_time := time.Now().Format("2006-01-02")
	where += " AND left(healthy_drug.created_at,10) = '" + day_time + "'"
	if kindergarten_id != 0 {
		where += "AND healthy_drug.kindergarten_id = ? "
		con = append(con, kindergarten_id)
	}
	if role != 1 {
		if class_id != 0 {
			where += "AND healthy_drug.class_id = ? "
			con = append(con, class_id)
		}
	}
	if types != 0 {
		where += "AND healthy_drug.types = ? "
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
		sql := qb.Select("healthy_drug.*,student.name").From(f.TableName()).
			LeftJoin("student").On("healthy_drug.student_id = student.student_id").
			Where(where).
			Limit(limit).Offset(start).String()

		if _, err := o.Raw(sql, con).Values(&sxWords); err == nil {

			pageNum := int(math.Ceil(float64(total) / float64(limit)))
			return Page{pageNum, limit, total, sxWords}, nil
		}
	}

	return Page{}, nil
}

//详情
func DrugInfo(id int) map[string]interface{} {
	o := orm.NewOrm()
	v := &Drug{Id: id}
	if err := o.Read(v); err == nil {
		paginatorMap := make(map[string]interface{})
		paginatorMap["data"] = v
		return paginatorMap
	}

	return nil
}
