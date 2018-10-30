package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"dailyfresh/models"
	"github.com/gomodule/redigo/redis"
	"strconv"
)

type UserAddCartController struct {
	beego.Controller
}

func (c *UserAddCartController) HandleUserAddCart() {
	//处理用户添加商品到购物车的请求
	resp := make(map[string]interface{})
	defer c.ServeJSON()
	skuId, err := c.GetInt("skuId")
	if err != nil {
		beego.Error("get skuId err:", err)
		resp["res"] = 1
		resp["errMsg"] = "获取商品Id失败"
		c.Data["json"] = resp
		return
	}
	goodsNum, err := c.GetInt("goodsNum")
	if err != nil {
		beego.Error("get goodsNum err:", err)
		resp["res"] = 1
		resp["errMsg"] = "获取数量信息失败"
		c.Data["json"] = resp
		return
	}
	//拿到skuId的时候需要存储用户id,商品id和商品数量
	userName := GetUser(&c.Controller)
	if userName == "" {
		beego.Error("用户未登录!")
		resp["res"] = 2
		resp["errMsg"] = "用户未登录!"
		c.Redirect("/", 302)
		c.Data["json"] = resp
		return
	}
	newOrm := orm.NewOrm()
	var user models.User
	user.Name = userName
	readErr := newOrm.Read(&user, "Name")
	if readErr != nil {
		beego.Error("用户不存在!")
		resp["res"] = 3
		resp["errMsg"] = "用户不存在!"
		c.Data["json"] = resp
		return
	}
	//检查要添加进数据库的商品数量是否超出库存量
	var goodsSku models.GoodsSKU
	goodsSku.Id = skuId
	readErr = newOrm.Read(&goodsSku)
	if readErr != nil {
		beego.Error("商品不存在!")
		resp["res"] = 4
		resp["errMsg"] = "商品不存在!"
		c.Data["json"] = resp
		return
	}
	conn, err := redis.Dial("tcp", "192.168.1.19:6379")
	if err != nil {
		beego.Error("连接数据库失败!")
		resp["res"] = 6
		resp["errMsg"] = "连接数据库失败!"
		c.Data["json"] = resp
		return
	}
	defer conn.Close()
	//添加之前要得到之前的sku数量,最后做累加
	reply, err := conn.Do("hget", "cart_"+strconv.Itoa(user.Id), skuId)
	preNum, _ := redis.Int(reply, err)
	if goodsNum+preNum > goodsSku.Stock {
		beego.Error("库存不足!")
		resp["res"] = 5
		resp["errMsg"] = "库存不足!"
		c.Data["json"] = resp
		return
	}
	//将数据存储到redis中
	conn.Do("hset", "cart_"+strconv.Itoa(user.Id), skuId, goodsNum+preNum)
	beego.Info("key is", "cart_"+strconv.Itoa(user.Id), "skuId", skuId)
	resp["res"] = 8
	resp["errMsg"] = ""
	//还需要查询当前的存储的总商品数返回给视图
	count := GetGoodsCount(conn, user.Id)
	resp["cartcount"] = count
	c.Data["json"] = resp
}
func GetGoodsCount(conn redis.Conn, userId int) int {
	reply, err := conn.Do("hgetall", "cart_"+strconv.Itoa(userId))
	numMap, err := redis.IntMap(reply, err)
	if err != nil {
		beego.Error("获取购物车总数失败")
		return 0
	}
	beego.Info("numMap is ", numMap)
	var count int
	for _, val := range numMap {
		count += val
	}
	return count
}
