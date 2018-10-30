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
	//默认展示页
	beego.Router("/", &controllers.IndexController{}, "get:ShowIndex")
	//注册页面
	beego.Router("/register", &controllers.RegisterController{}, "get:ShowRegister;post:HandleRegister")
	//登录页面
	beego.Router("/login", &controllers.LoginController{}, "get:ShowLogin;post:HandleLogin")
	//激活
	beego.Router("/active", &controllers.ActiveController{}, "get:ActiveUser")
	//用户中心
	beego.Router("/user/showUserInfo", &controllers.UserCenterInfoController{}, "get:ShowUserInfo")
	//用户订单
	beego.Router("/user/showUserOrder", &controllers.UserCenterOrderController{}, "get:ShowUserOrders")
	//用户购物车
	beego.Router("/user/showUserCart", &controllers.UserCartController{}, "get:ShowUserCart")
	//详情页添加商品到购物车
	beego.Router("/user/addCart", &controllers.UserAddCartController{}, "post:HandleUserAddCart")
	//添加默认地址
	beego.Router("/user/addDefaultAddr", &controllers.UserCenterAddrController{}, "get:ShowAddAddress;post:HandleAddAddress")
	//退出登录
	beego.Router("/user/logout", &controllers.LogoutController{}, "get:Logout")
	//展示商品详情
	beego.Router("/showGoodsDetail", &controllers.GoodsDetailController{}, "get:ShowDetail")
	//展示商品列表
	beego.Router("/showGoodsList", &controllers.GoodsListController{}, "get:ShowGoodsList")
	//搜索商品
	beego.Router("/searchGoods", &controllers.SearchGoodsController{}, "post:HandleSearchName")
	//更新购物车的商品数量
	beego.Router("/user/updateCart", &controllers.UpdateCartController{}, "post:HandleUpdateCart")
	//从购物车中删除商品
	beego.Router("/user/deleteGoodsFromCart", &controllers.DeleteFromCartController{}, "post:HandleDeleteFromCart")
	//展示订单的操作
	beego.Router("/user/showOrder", &controllers.OrderShowController{}, "post:HandleShowOrderCart")
	//提交订单的操作
	beego.Router("/user/postOrder", &controllers.OrderPostController{}, "post:HandlePostOrderCart")
	//付款请求
	beego.Router("/user/goPay", &controllers.PayController{}, "get:GoPay")
	//支付成功跳转的界面
	beego.Router("/user/payOk", &controllers.PayOkController{}, "get:PayOk")
	//短信发送业务
	beego.Router("/sendMessage", &controllers.MessageSendController{}, "get:SendMessage")



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
