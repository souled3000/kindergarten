package healthy

import (
	"github.com/astaxie/beego/orm"
	"strconv"
)

type Weight struct {
	Id       	int	   	 	`json:"id" orm:"column(id);auto"`
	Age			string	  	`json:"age" orm:"column(age)"`
	Small		string	  	`json:"small" orm:"column(small)"`
	SdOne		string	  	`json:"sd_one" orm:"column(sd_one)"`
	SdTwo		string		`json:"sd_two" orm:"column(sd_two)"`
	SdThree		string		`json:"sd_three" orm:"column(sd_three)"`
	SdFour		string		`json:"sd_four" orm:"column(sd_four)"`
	SdFive 		string		`json:"sd_five" orm:"column(sd_five)"`
	Large		string		`json:"large" orm:"column(large)"`
	Type		int			`json:"type" orm:"column(type)"`
}

func init() {
	orm.RegisterModel(new(Weight))
}

func (t *Weight) TableName() string {
	return "healthy_weight"
}

func (w *Weight) Compare(sex int,age,weight float64)  (types string,err error) {

	o := orm.NewOrm()
	var status string
	var info Weight
	err = o.QueryTable("healthy_weight").Filter("type",sex).Filter("age", age).One(&info)
	if err == nil{

		small := info.Small
		sml,_ := strconv.ParseFloat(small, 64)

		large := info.Large
		len,_ := strconv.ParseFloat(large, 64)

		if weight < sml {
			status = "1"
		}
		if weight > len {
			status = "2"
		}else {
			status = "0"
		}
	}else {
		return "",err
	}

	return status,nil
}