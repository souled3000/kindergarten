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
	"kindergarten-service-go/controllers/admin"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/plugins/cors"
	"kindergarten-service-go/controllers/healthy"
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

		beego.NSNamespace("/healthy/drug",
			beego.NSInclude(
				&healthy.DrugController{},
			),
		),

		beego.NSNamespace("/healthy/inspect",
			beego.NSInclude(
				&healthy.InspectController{},
			),
		),

		beego.NSNamespace("/healthy/body",
			beego.NSInclude(
				&healthy.BodyController{},
			),
		),

		beego.NSNamespace("/healthy/class",
			beego.NSInclude(
				&healthy.ClassController{},
			),
		),

		beego.NSNamespace("/healthy/situation",
			beego.NSInclude(
				&healthy.SituationController{},
			),
		),

		beego.NSNamespace("/healthy/column",
			beego.NSInclude(
				&healthy.ColumnController{},
			),
		),

		beego.NSNamespace("/admin/kindergarten",
			beego.NSInclude(
				&admin.KindergartenController{},
			),
		),

		beego.NSNamespace("/admin/kindergarten_life",
			beego.NSInclude(
				&admin.KindergartenLifeController{},
			),
		),

		beego.NSNamespace("/admin/notice",
			beego.NSInclude(
				&admin.NoticeController{},
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

		beego.NSNamespace("/admin/organizational",
			beego.NSInclude(
				&admin.OrganizationalController{},
			),
		),

		beego.NSNamespace("/admin/organizational_member",
			beego.NSInclude(
				&admin.OrganizationalMemberController{},
			),
		),

		beego.NSNamespace("/admin/permission",
			beego.NSInclude(
				&admin.PermissionController{},
			),
		),

		beego.NSNamespace("/admin/role",
			beego.NSInclude(
				&admin.RoleController{},
			),
		),

		beego.NSNamespace("/admin/route",
			beego.NSInclude(
				&admin.RouteController{},
			),
		),

		beego.NSNamespace("/admin/student",
			beego.NSInclude(
				&admin.StudentController{},
			),
		),

		beego.NSNamespace("/admin/teacher",
			beego.NSInclude(
				&admin.TeacherController{},
			),
		),

		beego.NSNamespace("/user_permission",
			beego.NSInclude(
				&controllers.UserPermissionController{},
			),
		),

		beego.NSNamespace("/admin/user_permission",
			beego.NSInclude(
				&admin.UserPermissionController{},
			),
		),

		beego.NSNamespace("/ping",
			beego.NSInclude(
				&admin.PingController{},
			),
		),

		beego.NSNamespace("/admin/ping",
			beego.NSInclude(
				&admin.PingController{},
			),
		),
		beego.NSNamespace("/app/visitors",
			beego.NSInclude(
				&controllers.KindergartenVisitorsController{},
			),
		),
		beego.NSNamespace("/app/special_child",
			beego.NSInclude(
				&controllers.ExceptionalChildController{},
			),
		),

		beego.NSNamespace("/admin/special_child",
			beego.NSInclude(
				&admin.ExceptionalChildController{},
			),
		),
	)
	beego.AddNamespace(ns)
}
