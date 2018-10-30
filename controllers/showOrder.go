package controllers

import (
	"github.com/astaxie/beego"
	"dailyfresh/models"
	"github.com/astaxie/beego/orm"
	"github.com/gomodule/redigo/redis"
	"strconv"
)

type OrderShowController struct {
	beego.Controller
}

func (c *OrderShowController) HandleShowOrderCart() {
	userName := GetUser(&c.Controller)
	if userName == "" {
		beego.Error("用户未登录")
		c.Redirect("/login", 302)
		return
	}
	var user models.User
	user.Name = userName
	newOrm := orm.NewOrm()
	readErr := newOrm.Read(&user, "Name")
	if readErr != nil {
		beego.Error("用户不存在")
		c.Redirect("/login", 302)
		return
	}
	//显示订单确认界面,需要获取的信息有默认的收货地址,支付方式,商品数量,金额,运费,总价格
	skuIdSlice := c.GetStrings("skuId")
	if len(skuIdSlice) == 0 {
		beego.Error("购物车数据错误")
		c.Redirect("/user/showUserCart", 302)
		return
	}
	//将前台需要的数据存储到map集合中
	skuDatas := make([]map[string]interface{}, 0)
	conn, err := redis.Dial("tcp", "192.168.1.19:6379")
	if err != nil {
		beego.Error("连接redis失败")
		c.Redirect("/user/showUserCart", 302)
		return
	}
	defer conn.Close()
	var goodsAmount, totalPrice int
	for index, skuId := range skuIdSlice {
		//循环从redis读取数据
		temp := make(map[string]interface{})
		reply, err := conn.Do("hget", "cart_"+strconv.Itoa(user.Id), skuId)
		number, err := redis.Int(reply, err)
		if err != nil {
			beego.Error("获取商品数量出错,err is:", err)
			c.Redirect("/user/showUserCart", 302)
			return
		}
		var goods models.GoodsSKU
		goods.Id, err = strconv.Atoi(skuId)
		if err != nil {
			beego.Error("转换商品数量出错")
			c.Redirect("/user/showUserCart", 302)
			return
		}
		readErr = newOrm.Read(&goods)
		if readErr != nil {
			beego.Error("读取商品信息出错")
			c.Redirect("/user/showUserCart", 302)
			return
		}

		price := goods.Price
		subTotal := price * number
		temp["goods"] = goods
		temp["number"] = number
		temp["subTotal"] = subTotal
		temp["index"] = index + 1
		goodsAmount += number
		totalPrice += subTotal
		skuDatas = append(skuDatas, temp)
	}
	//查询地址信息,返回地址列表
	var addrs []models.Address
	_, err = newOrm.QueryTable("Address").RelatedSel("User").Filter("User__Id", user.Id).All(&addrs)
	if err != nil {
		beego.Error("查询地址数据失败,err:", err)
		addrs = append(addrs, models.Address{Addr: "无默认地址,请添加"})
	}
	//循环的从数据库中读取数据
	transit := 10 //运费
	shouldPayPrice := totalPrice + transit
	c.Data["skuDatas"] = skuDatas
	c.Data["goodsAmount"] = goodsAmount       //总件数
	c.Data["totalPrice"] = totalPrice         //总价格
	c.Data["shouldPayPrice"] = shouldPayPrice //实付款
	c.Data["addrs"] = addrs                   //地址列表
	c.Data["transit"] = transit               //实付款
	c.Data["skuIdSlice"] = skuIdSlice               //所有的商品id
	c.Layout = "userCenterLayout.html"
	beego.Info("skuDatas:", skuDatas, "goodsAmount:",
		goodsAmount, "totalPrice:", totalPrice, "shouldPayPrice:", shouldPayPrice, "addrs", addrs)
	c.TplName = "place_order.html"
}
