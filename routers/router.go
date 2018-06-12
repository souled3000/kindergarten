// @APIVersion 1.0.0
// @Title beego Test API
// @Description beego has a very cool tools to autogenerate documents for your API
// @Contact astaxie@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"kindergarten-service-go/controllers"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/plugins/cors"
)

func init() {
	beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type"},
		AllowCredentials: true,
	}))
	ns := beego.NewNamespace("/api/v2/kg",

		beego.NSNamespace("/kindergarten",
			beego.NSInclude(
				&controllers.KindergartenController{},
			),
		),

		beego.NSNamespace("/kindergarten_life",
			beego.NSInclude(
				&controllers.KindergartenLifeController{},
			),
		),

		beego.NSNamespace("/notice",
			beego.NSInclude(
				&controllers.NoticeController{},
			),
		),

		beego.NSNamespace("/organizational",
			beego.NSInclude(
				&controllers.OrganizationalController{},
			),
		),

		beego.NSNamespace("/organizational_member",
			beego.NSInclude(
				&controllers.OrganizationalMemberController{},
			),
		),

		beego.NSNamespace("/permission",
			beego.NSInclude(
				&controllers.PermissionController{},
			),
		),

		beego.NSNamespace("/role",
			beego.NSInclude(
				&controllers.RoleController{},
			),
		),

		beego.NSNamespace("/route",
			beego.NSInclude(
				&controllers.RouteController{},
			),
		),

		beego.NSNamespace("/student",
			beego.NSInclude(
				&controllers.StudentController{},
			),
		),

		beego.NSNamespace("/teacher",
			beego.NSInclude(
				&controllers.TeacherController{},
			),
		),

		beego.NSNamespace("/user_permission",
			beego.NSInclude(
				&controllers.UserPermissionController{},
			),
		),

		beego.NSNamespace("/ping",
			beego.NSInclude(
				&controllers.PingController{},
			),
		),
	)
	beego.AddNamespace(ns)
}
