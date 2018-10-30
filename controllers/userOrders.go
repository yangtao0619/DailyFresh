package controllers

import (
	"github.com/astaxie/beego"
	"dailyfresh/models"
	"github.com/astaxie/beego/orm"
	"math"
)

type UserCenterOrderController struct {
	beego.Controller
}

//展示该用户的所有订单信息
func (c *UserCenterOrderController) ShowUserOrders() {
	userName := GetUser(&c.Controller)
	if userName == "" {
		beego.Error("用户未登录")
		c.Redirect("/", 302)
		return
	}
	var user models.User
	user.Name = userName
	newOrm := orm.NewOrm()
	readErr := newOrm.Read(&user, "Name")
	if readErr != nil {
		beego.Error("读取用户信息错误")
		c.Redirect("/", 302)
		return
	}
	//读取订单表
	var orderInfos []models.OrderInfo
	newOrm.QueryTable("OrderInfo").Filter("User", user).All(&orderInfos)
	//页码的发送
	count := len(orderInfos)
	pageSize := 2
	pageCount := math.Ceil(float64(count) / float64(pageSize))
	beego.Info("pageCount is", pageCount)
	pageIndex, err := c.GetInt("pageIndex", 1)
	beego.Info("pageIndex is", pageIndex)
	if err != nil {
		beego.Error("获取当前页码失败,err:", err)
	}
	pageDataMap := listPageTool(int(pageCount), pageIndex)
	pageDataMap["pageIndex"] = pageIndex
	c.Data["pageBuffer"] = pageDataMap["pageBuffer"].([]int)
	c.Data["prePage"] = pageDataMap["prePage"]
	c.Data["nextPage"] = pageDataMap["nextPage"]
	c.Data["pageIndex"] = pageDataMap["pageIndex"]

	//根据所有订单信息查询对应的订单商品集合
	newOrm.QueryTable("OrderInfo").Filter("User", user).Limit(pageSize, (pageIndex-1)*pageSize).All(&orderInfos)
	var newInfos []models.OrderInfo
	for _, orderInfo := range orderInfos {
		var orderGoods []*models.OrderGoods
		newOrm.QueryTable("OrderGoods").RelatedSel("GoodsSKU", "OrderInfo").Filter("OrderInfo",
			orderInfo).All(&orderGoods)
		for _, orderG := range orderGoods {
			beego.Info("orderGoods:", *orderG.GoodsSKU)
		}
		orderInfo.OrderGoods = orderGoods
		beego.Info("orderInfo time", orderInfo.Time.String())
		newInfos = append(newInfos, orderInfo)
	}
	//查询订单商品表
	c.Data["orderInfos"] = newInfos
	c.Layout = "userCenterLayout.html"
	c.TplName = "user_center_order.html"
}
