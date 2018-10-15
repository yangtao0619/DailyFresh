package controllers

import "github.com/astaxie/beego"

type UserCenterOrderController struct {
	beego.Controller
}

//展示该用户的所有订单信息
func (c *UserCenterOrderController) ShowUserOrders() {
	userName := c.GetSession("username")
	c.Data["username"] = userName
	c.Layout = "layout.html"
	c.TplName = "user_center_order.html"
}
