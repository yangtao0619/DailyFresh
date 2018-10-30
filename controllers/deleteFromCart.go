package controllers

import (
	"github.com/astaxie/beego"
	"github.com/gomodule/redigo/redis"
	"dailyfresh/models"
	"github.com/astaxie/beego/orm"
	"strconv"
)

type DeleteFromCartController struct {
	beego.Controller
}

func (c *DeleteFromCartController) HandleDeleteFromCart() {
	resp := make(map[string]interface{})
	defer c.ServeJSON()
	skuId, err := c.GetInt("skuId")
	if err != nil {
		resp["errId"] = 0
		resp["reeMsg"] = "获取skuid失败"
		beego.Error("获取skuid失败")
		c.Data["json"] = resp
		return
	}
	conn, err := redis.Dial("tcp", "192.168.1.19:6379")
	if err != nil {
		resp["errId"] = 2
		resp["reeMsg"] = "连接redis数据库失败"
		beego.Error("连接redis数据库失败")
		c.Data["json"] = resp
		return
	}
	userName := GetUser(&c.Controller)
	if userName == "" {
		resp["errId"] = 3
		resp["reeMsg"] = "用户未登录"
		beego.Error("用户未登录")
		c.Data["json"] = resp
		return
	}
	var user models.User
	user.Name = userName
	newOrm := orm.NewOrm()
	readErr := newOrm.Read(&user, "Name")
	if readErr != nil {
		resp["errId"] = 4
		resp["reeMsg"] = "用户不存在"
		beego.Error("用户不存在")
		c.Data["json"] = resp
		return
	}
	_, err = conn.Do("hdel", "cart_"+strconv.Itoa(user.Id), skuId)
	if err != nil{
		resp["errId"] = 5
		resp["reeMsg"] = "删除数据失败"
		beego.Error("删除数据失败")
		c.Data["json"] = resp
		return
	}
	resp["errId"] = 6
	resp["reeMsg"] = ""
	c.Data["json"] = resp
}
