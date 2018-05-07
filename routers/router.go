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
)

func init() {
	ns := beego.NewNamespace("/api/v2/kg",

		beego.NSNamespace("/facilities_display",
			beego.NSInclude(
				&controllers.FacilitiesDisplayController{},
			),
		),

		beego.NSNamespace("/grade",
			beego.NSInclude(
				&controllers.GradeController{},
			),
		),

		beego.NSNamespace("/group_view",
			beego.NSInclude(
				&controllers.GroupViewController{},
			),
		),

		beego.NSNamespace("/kindergarten",
			beego.NSInclude(
				&controllers.KindergartenController{},
			),
		),

		beego.NSNamespace("/kindergarten_courseware",
			beego.NSInclude(
				&controllers.KindergartenCoursewareController{},
			),
		),

		beego.NSNamespace("/kindergarten_folder",
			beego.NSInclude(
				&controllers.KindergartenFolderController{},
			),
		),

		beego.NSNamespace("/kindergarten_life",
			beego.NSInclude(
				&controllers.KindergartenLifeController{},
			),
		),

		beego.NSNamespace("/kinship",
			beego.NSInclude(
				&controllers.KinshipController{},
			),
		),

		beego.NSNamespace("/life_picture",
			beego.NSInclude(
				&controllers.LifePictureController{},
			),
		),

		beego.NSNamespace("/member",
			beego.NSInclude(
				&controllers.MemberController{},
			),
		),

		beego.NSNamespace("/migrations",
			beego.NSInclude(
				&controllers.MigrationsController{},
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

		beego.NSNamespace("/permission_route",
			beego.NSInclude(
				&controllers.PermissionRouteController{},
			),
		),

		beego.NSNamespace("/role",
			beego.NSInclude(
				&controllers.RoleController{},
			),
		),

		beego.NSNamespace("/role_permission",
			beego.NSInclude(
				&controllers.RolePermissionController{},
			),
		),

		beego.NSNamespace("/route",
			beego.NSInclude(
				&controllers.RouteController{},
			),
		),

		beego.NSNamespace("/slide_show",
			beego.NSInclude(
				&controllers.SlideShowController{},
			),
		),

		beego.NSNamespace("/student",
			beego.NSInclude(
				&controllers.StudentController{},
			),
		),

		beego.NSNamespace("/student_ examination",
			beego.NSInclude(
				&controllers.StudentExaminationController{},
			),
		),

		beego.NSNamespace("/student_kindergarten",
			beego.NSInclude(
				&controllers.StudentKindergartenController{},
			),
		),

		beego.NSNamespace("/teacher",
			beego.NSInclude(
				&controllers.TeacherController{},
			),
		),

		beego.NSNamespace("/teacher_title",
			beego.NSInclude(
				&controllers.TeacherTitleController{},
			),
		),

		beego.NSNamespace("/teachers_show",
			beego.NSInclude(
				&controllers.TeachersShowController{},
			),
		),

		beego.NSNamespace("/user_permission",
			beego.NSInclude(
				&controllers.UserPermissionController{},
			),
		),

		beego.NSNamespace("/user_role",
			beego.NSInclude(
				&controllers.UserRoleController{},
			),
		),
	)
	beego.AddNamespace(ns)
}
