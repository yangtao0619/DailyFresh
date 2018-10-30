package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

type PayOkController struct {
	beego.Controller
}

func (c *PayOkController) PayOk() {
	//获取orderid
	orderId := c.GetString("out_trade_no")
	if orderId == "" {
		beego.Error("获取订单id失败")
		c.Redirect("/user/showUserOrder", 302)
		return
	}

	//更新数据库
	newOrm := orm.NewOrm()
	newOrm.QueryTable("OrderInfo").Filter("OrderId", orderId).Update(orm.Params{"Orderstatus": 2})
	c.Redirect("/user/showUserOrder", 302)
}
