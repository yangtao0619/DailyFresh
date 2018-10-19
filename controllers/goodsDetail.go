package controllers

import (
	"github.com/astaxie/beego"
	"dailyfresh/models"
	"github.com/astaxie/beego/orm"
	"github.com/gomodule/redigo/redis"
	"errors"
	"strconv"
)

type GoodsDetailController struct {
	beego.Controller
}

func (c *GoodsDetailController) ShowDetail() {
	//展示商品详情页面,获取传递的id值
	id, err := c.GetInt("id")
	//校验数据
	if err != nil {
		beego.Error("get int err:", err)
		c.Redirect("/", 302)
		return
	}
	newOrm := orm.NewOrm()

	//处理数据,根据id找到对应的sku
	var sku models.GoodsSKU
	newOrm.QueryTable("GoodsSKU").RelatedSel("Goods", "GoodsType").Filter("Id", id).One(&sku)
	beego.Info("sku is", sku)
	c.Data["goods"] = sku
	GetUser(&c.Controller)
	ShowTypeAndGetUserInfo(&c.Controller)
	//获取同类型的新品推荐
	var newGoods []models.GoodsSKU
	newOrm.QueryTable("GoodsSKU").RelatedSel("GoodsType").Filter("GoodsType", sku.GoodsType).
		OrderBy("Time").Limit(2, 0).All(&newGoods)
	c.Data["newGoods"] = newGoods
	err = recordView(c, &newOrm, id)
	if err != nil {
		beego.Error(err)
	}
	c.TplName = "detail.html"
}
func recordView(c *GoodsDetailController, newOmer *orm.Ormer, goodsId int) error {
	//将用户的浏览记录写入redis
	userName := c.GetSession("username")
	var user models.User
	//需要先判断用户是否已经登录,只有登录之后才需要写入

	if userName != nil {
		user.Name = userName.(string)
		//key是用户的id,value是浏览的商品的id
		conn, err := GetRedisConnect()
		if err != nil {
			return errors.New("redis连接出错")
		}
		readErr := (*newOmer).Read(&user, "Name")
		if readErr != nil {
			beego.Error("readerr:", readErr)
		}
		beego.Info("insert user id is", user.Id, "username is", userName)
		reply, err := conn.Do("lrem", "history"+strconv.Itoa(user.Id), 0, goodsId)
		insertResult, _ := redis.Bool(reply, err)
		if !insertResult {
			err = errors.New("删除数据失败")
		}
		reply, err = conn.Do("lpush", "history"+strconv.Itoa(user.Id), goodsId)
		insertResult, _ = redis.Bool(reply, err)
		if !insertResult {
			return errors.New("插入数据失败")
		}
	}
	beego.Info("goodId is", goodsId, "userId is", user.Id)
	return nil
}

func ShowTypeAndGetUserInfo(c *beego.Controller) {
	//查找类型数据
	var types []models.GoodsType
	ormer := orm.NewOrm()
	ormer.QueryTable("GoodsType").All(&types)
	c.Data["types"] = types
	c.Layout = "goodLayout.html"
}
