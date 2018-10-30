package controllers

import (
	"github.com/astaxie/beego"
	"dailyfresh/models"
	"github.com/astaxie/beego/orm"
	"strconv"
	"github.com/gomodule/redigo/redis"
)

type UserCartController struct {
	beego.Controller
}

func (c *UserCartController) ShowUserCart() {
	//显示用户名
	userName := GetUser(&c.Controller)
	var user models.User
	user.Name = userName
	newOrm := orm.NewOrm()
	readErr := newOrm.Read(&user, "Name")
	if readErr != nil {
		beego.Error("读取用户信息失败")
		c.Redirect("/", 302)
		return
	}
	var dataSlice []map[string]interface{}
	//需要传递的数据,查询redis,有几个field切片就有多大
	conn, err := redis.Dial("tcp", "212.64.52.176:6379")
	if err != nil {
		beego.Error("redis连接错误")
		c.Redirect("/", 302)
		return
	}
	defer conn.Close()
	reply, err := conn.Do("hgetall", "cart_"+strconv.Itoa(user.Id))
	goodsInfoMap, err := redis.IntMap(reply, err)
	if err != nil {
		beego.Error("转换map失败")
		c.Redirect("/", 302)
		return
	}
	var totalPrice, totalNum int
	for skuId, skuNum := range goodsInfoMap {
		itemMap := make(map[string]interface{})
		//根据skuid找到sku
		var goodsSku models.GoodsSKU
		goodsSku.Id, err = strconv.Atoi(skuId)
		if err != nil {
			beego.Error("转换id失败")
			return
		}
		newOrm.Read(&goodsSku)
		price := goodsSku.Price
		subtotal := skuNum * price
		itemMap["goods"] = goodsSku
		itemMap["num"] = skuNum
		itemMap["subtotal"] = subtotal
		dataSlice = append(dataSlice, itemMap)

		totalNum += skuNum
		totalPrice += subtotal
	}
	c.Data["dataSlice"] = dataSlice
	c.Data["totalPrice"] = totalPrice
	c.Data["totalNum"] = totalNum

	beego.Info("dataSlice",dataSlice)
	beego.Info("totalPrice",totalPrice)
	beego.Info("totalNum",totalNum)
	ShowTypeAndGetUserInfo(&c.Controller)
	c.Layout = "goodLayout.html"
	c.TplName = "cart.html"
}
