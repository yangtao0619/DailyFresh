package routers

import (
	"dailyfresh/controllers"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
)

func init() {
	//设置过滤器
	beego.InsertFilter("/goods/*", beego.BeforeExec, HandleGoodsFunc)
	beego.InsertFilter("/user/*", beego.BeforeExec, HandleUserFunc)

	beego.Router("/", &controllers.IndexController{},"get:ShowIndex")
	beego.Router("/register", &controllers.RegisterController{}, "get:ShowRegister;post:HandleRegister")
	beego.Router("/login", &controllers.LoginController{}, "get:ShowLogin;post:HandleLogin")
	beego.Router("/active", &controllers.ActiveController{}, "get:ActiveUser")
	beego.Router("/user/showUserInfo", &controllers.UserCenterInfoController{}, "get:ShowUserInfo")
	beego.Router("/user/showUserOrder", &controllers.UserCenterOrderController{}, "get:ShowUserOrders")
	beego.Router("/user/addDefaultAddr", &controllers.UserCenterAddrController{}, "get:ShowAddAddress;post:HandleAddAddress")
	beego.Router("/user/logout", &controllers.LogoutController{}, "get:Logout")
	beego.Router("/showGoodsDetail", &controllers.GoodsDetailController{}, "get:ShowDetail")
	beego.Router("/showGoodsList", &controllers.GoodsListController{}, "get:ShowGoodsList")
}

var HandleGoodsFunc = func(c *context.Context) {
	//先查看是否已经登录,没有的话就重定向到登录界面
	userName := c.Input.Session("username")
	if userName == nil {
		beego.Error("需要登录后进行此操作")
		c.Redirect(302, "/")
	} else {
		return
	}
}

var HandleUserFunc = func(c *context.Context) {
	//先查看是否已经登录,没有的话就重定向到登录界面
	userName := c.Input.Session("username")
	if userName == nil {
		beego.Error("需要登录后进行此操作")
		c.Redirect(302, "/login")
	} else {
		return
	}
}
