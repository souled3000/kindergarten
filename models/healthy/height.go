package healthy

import (
	"fmt"
	"github.com/astaxie/beego/orm"
)

type Height struct {
	Id      int     `json:"id" orm:"column(id);auto"`
	Age     int     `json:"age" orm:"column(age)"`
	Small   float64 `json:"small" orm:"column(small)"`
	SdOne   float64 `json:"sd_one" orm:"column(sd_one)"`
	SdTwo   float64 `json:"sd_two" orm:"column(sd_two)"`
	SdThree float64 `json:"sd_three" orm:"column(sd_three)"`
	Large   float64 `json:"large" orm:"column(large)"`
	Type    int     `json:"type" orm:"column(type)"`
}

func init() {
	orm.RegisterModel(new(Height))
}

func (t *Height) TableName() string {
	return "healthy_height"
}

//根据年龄判断身高
func CompareHeight(sex int, age, weight float64) (types string, err error) {

	o := orm.NewOrm()
	var status string
	var info Height
	err = o.QueryTable("healthy_height").Filter("type", sex).Filter("age__lte", age).One(&info)
	fmt.Println(info)
	if err == nil {

		if weight < info.Small {
			status = "矮小" //矮小
		} else if weight >= info.Small && weight < info.SdOne {
			status = "偏矮" //偏矮
		} else if weight >= info.SdOne && weight < info.SdThree {
			status = "正常" //正常
		} else if weight >= info.SdThree && weight < info.Large {
			status = "偏高" //偏胖
		} else if weight >= info.Large {
			status = "超高"
		}
	} else {
		return "", err
	}

	return status, nil
}
