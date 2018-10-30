package controllers

import (
	"github.com/astaxie/beego"
	"dailyfresh/models"
	"github.com/astaxie/beego/orm"
	"github.com/gomodule/redigo/redis"
	"errors"
	"strconv"
)

type UpdateCartController struct {
	beego.Controller
}

func (c *UpdateCartController) HandleUpdateCart() {
	resp := make(map[string]interface{})
	defer c.ServeJSON()
	//判断用户是否登录
	userName := GetUser(&c.Controller)
	if userName == "" {
		resp["respCode"] = 1
		beego.Error("用户未登录")
		resp["errMsg"] = "用户未登录"
		c.Data["json"] = resp
		return
	}
	var user models.User
	user.Name = userName
	newOrm := orm.NewOrm()
	readErr := newOrm.Read(&user, "Name")
	//用户检测
	if readErr != nil {
		resp["respCode"] = 2
		beego.Error("用户不存在")
		resp["errMsg"] = "用户不存在"
		c.Data["json"] = resp
		return
	}
	//商品检测
	skuId, err := c.GetInt("skuId")
	var goodsSku models.GoodsSKU
	goodsSku.Id = skuId
	requestNum, err2 := c.GetInt("setNum")
	if err != nil || err2 != nil {
		resp["respCode"] = 3
		beego.Error("请求数据不正确")
		resp["errMsg"] = "请求数据不正确"
		c.Data["json"] = resp
		return
	}
	readErr = newOrm.Read(&goodsSku)
	if readErr != nil {
		resp["respCode"] = 4
		beego.Error("商品不存在")
		resp["errMsg"] = "商品不存在"
		c.Data["json"] = resp
		return
	}
	//库存量检测
	stock := goodsSku.Stock
	var setNum int
	if stock < requestNum {
		resp["respCode"] = 5
		resp["errMsg"] = "库存不足,库存量是:" + strconv.Itoa(stock)
		beego.Error("库存不足,库存量是:", stock)
		setNum = stock
	} else {
		resp["respCode"] = 6
		setNum = requestNum
	}
	//将数据存储到数据库中
	conn, err := redis.Dial("tcp", "192.168.1.19:6379")
	if err != nil {
		err = errors.New("redis连接错误")
		resp["respCode"] = 7
		resp["errMsg"] = "redis连接错误"
		beego.Error("redis连接错误")
		c.Data["json"] = resp
		return
	}
	defer conn.Close()
	//如果正常连接redis,就向redis中存入数据
	_, err = conn.Do("hset", "cart_"+strconv.Itoa(user.Id), skuId, setNum)
	beego.Info("cart_"+strconv.Itoa(user.Id), skuId, setNum)
	if err != nil {
		err = errors.New("redis存储错误")
		resp["respCode"] = 8
		resp["errMsg"] = "redis存储错误"
		beego.Error("redis存储错误")
		c.Data["json"] = resp
		return
	}
	resp["setNum"] = setNum
	//返回一个json,带有响应码和商品的数量
	c.Data["json"] = resp
	beego.Info("resp is", resp)
}
