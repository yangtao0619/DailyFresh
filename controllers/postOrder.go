package controllers

import (
	"github.com/astaxie/beego"
	"dailyfresh/models"
	"time"
	"github.com/astaxie/beego/orm"
	"strconv"
	"github.com/gomodule/redigo/redis"
	"strings"
)

type OrderPostController struct {
	beego.Controller
}

func (c *OrderPostController) HandlePostOrderCart() {
	resp := make(map[string]interface{})
	defer c.ServeJSON()
	userName := GetUser(&c.Controller)
	if userName == "" {
		resp["respCode"] = 1
		resp["errMsg"] = "用户未登录"
		beego.Error("用户未登录")
		c.Data["json"] = resp
		return
	}
	var user models.User
	user.Name = userName
	newOrm := orm.NewOrm()
	newOrm.Begin()
	readErr := newOrm.Read(&user, "Name")
	if readErr != nil {
		resp["respCode"] = 2
		resp["errMsg"] = "账号不存在"
		beego.Error("账号不存在")
		c.Data["json"] = resp
		return
	}
	//提交订单的请求发到这里,写入数据到数据库中
	skuIdSlice := strings.Split(strings.Replace(strings.Replace(c.GetString("skuIds"), "[", "", -1),
		"]", "", -1), " ")
	beego.Info("skuIdSlice ", skuIdSlice)
	addrId, err := c.GetInt("addrId")
	payMethod, err := c.GetInt("payMethod")
	goodsAmount, err := c.GetInt("goodsAmount")
	goodsPrice, err := c.GetInt("goodsPrice")
	transitPrice, err := c.GetInt("transitPrice")
	if err != nil {
		resp["respCode"] = 3
		resp["errMsg"] = "前台数据错误"
		beego.Error("前台数据错误")
		c.Data["json"] = resp
		return
	}
	var addr models.Address
	addr.Id = addrId
	readErr = newOrm.Read(&addr)
	if readErr != nil {
		resp["respCode"] = 4
		resp["errMsg"] = "读取地址错误"
		beego.Error("读取地址错误")
		c.Data["json"] = resp
		return
	}
	//订单数据写入订单表
	var orderInfo models.OrderInfo
	orderInfo.OrderId = time.Now().Format("20060102150405") + strconv.Itoa(user.Id)
	orderInfo.User = &user
	orderInfo.Address = &addr
	orderInfo.PayMethod = payMethod
	orderInfo.TotalCount = goodsAmount
	orderInfo.Orderstatus = 1 //未支付
	orderInfo.TotalPrice = goodsPrice
	orderInfo.TransitPrice = transitPrice
	newOrm.Insert(&orderInfo)
	//订单商品数据写入订单商品表
	conn, err := redis.Dial("tcp", "192.168.1.19:6379")
	if err != nil {
		resp["respCode"] = 5
		resp["errMsg"] = "redis连接错误"
		beego.Error("redis连接错误")
		c.Data["json"] = resp
		return
	}
	defer conn.Close()
	for _, skuId := range skuIdSlice {
		var goodsSku models.GoodsSKU
		goodsSkuId, err := strconv.Atoi(skuId)
		if err != nil {
			resp["respCode"] = 5
			resp["errMsg"] = "读取商品错误"
			beego.Error("读取商品错误,err:", err)
			c.Data["json"] = resp
			return
		}
		goodsSku.Id = goodsSkuId
		//查询商品数量
		newOrm.Read(&goodsSku)
		i := 3
		for i > 0 {
			reply, err := conn.Do("hget", "cart_"+strconv.Itoa(user.Id), goodsSku.Id)
			count, err := redis.Int(reply, err)
			if err != nil {
				resp["respCode"] = 6
				resp["errMsg"] = "获取商品数量失败"
				beego.Error("获取商品数量失败")
				c.Data["json"] = resp
				return
			}
			var orderGoods models.OrderGoods
			orderGoods.OrderInfo = &orderInfo
			//在修改商品的数量之前需要判断此时的库存是否充足
			if goodsSku.Stock < count {
				newOrm.Rollback()
				resp["respCode"] = 7
				resp["errMsg"] = "商品库存不足"
				beego.Error("商品库存不足")
				c.Data["json"] = resp
				return
			}
			preCount := goodsSku.Stock
			//查询商品的数量
			orderGoods.Count = count
			orderGoods.Price = goodsSku.Price
			//从redis删除购物车数据

			orderGoods.GoodsSKU = &goodsSku
			_, err = newOrm.Insert(&orderGoods)
			if err != nil {
				beego.Error("insert err:", err)
				return
			}
			//删除库存
			goodsSku.Stock -= count
			goodsSku.Sales += count
			//newOrm.Update(&orderGoods)
			num, _ := newOrm.QueryTable("GoodsSKU").Filter("Id", goodsSku.Id).Filter("Stock", preCount).
				Update(orm.Params{"Stock": goodsSku.Stock, "Sales": goodsSku.Sales})
			beego.Info("num is", num)
			if num == 0 {
				if i > 0 {
					i -= 1
					continue
				}
				newOrm.Rollback()
				resp["respCode"] = 9
				resp["errMsg"] = "库存改变,请重新提交"
				beego.Error("库存改变,请重新提交")
				c.Data["json"] = resp
				return
			} else {
				conn.Do("hdel", "cart_"+strconv.Itoa(user.Id), skuId)
				break
			}
		}

	}
	newOrm.Commit()
	resp["respCode"] = 8
	resp["errMsg"] = ""
	c.Data["json"] = resp
}
