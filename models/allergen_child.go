package models

import (
	"strings"
	"github.com/astaxie/beego/orm"
)

func GetAllergenChild(allergen string) (allergenChild []map[string]interface{}, err error) {
	param := strings.Split(allergen, ",")
	where := "1=1 "
	qb, _ := orm.NewQueryBuilder("mysql")
	o := orm.NewOrm()
	var allergens []orm.Params
	for _, v := range param {
		if v != "" {
			maps := make(map[string]interface{})
			where += " AND allergen like \"%" + string(v) + "%\""
			sql := qb.Select("allergen").From("exceptional_child").Where(where).String()
			if num, err := o.Raw(sql).Values(&allergens); err == nil && num > 0 {
				if childName, err := GetChildName(where); err == nil {
					maps["allergen"] = v
			    	maps["child_name"] = childName
			        allergenChild = append(allergenChild, maps)
			        return allergenChild, err
			    }
			} else {
				return allergenChild, err
			}
		}
	}
	return
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
	 		childName = string(s[:len(s) - 1])
		}
		return childName, err
	 }
	 return  childName ,err
}