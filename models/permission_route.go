package models

import (
	"github.com/astaxie/beego/orm"
)

type PermissionRoute struct {
	Id           int `json:"id" orm:"column(id);auto"`
	PermissionId int `json:"permission_id" orm:"column(permission_id)" description:"权限ID"`
	RouteId      int `json:"route_id" orm:"column(route_id)" description:"路由ID"`
}

func (t *PermissionRoute) TableName() string {
	return "permission_route"
}

func init() {
	orm.RegisterModel(new(PermissionRoute))
}
