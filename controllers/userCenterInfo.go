package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"dailyfresh/models"
	"strconv"
	"github.com/gomodule/redigo/redis"
	"errors"
)

type UserCenterInfoController struct {
	beego.Controller
}

func (c *UserCenterInfoController) ShowUserInfo() {
	//展示用户信息表
	userName := GetUser(&c.Controller)
	//查询用户信息
	newOrm := orm.NewOrm()
	var user models.User
	user.Name = userName
	newOrm.Read(&user, "Name")
	//多表查询,查询该用户的默认地址
	var address models.Address
	newOrm.QueryTable("Address").Filter("User__Id", user.Id).Filter("Isdefault", true).One(&address)

	//进行浏览记录的展示.获取五条就可以
	err := showViewRecords(&user, c, &newOrm)
	if err != nil {
		beego.Error("search view record err:", err)
	}

	c.Data["address"] = address.Addr
	c.Data["phoneNumber"] = address.Phone
	c.Layout = "userCenterLayout.html"
	c.TplName = "user_center_info.html"
}
func showViewRecords(user *models.User, c *UserCenterInfoController, newOrm *orm.Ormer) (err error) {
	id := user.Id
	conn, err := redis.Dial("tcp", "212.64.52.176:6379")
	if err != nil {
		err = errors.New("redis连接错误")
	}
	defer conn.Close()
	beego.Info("id is", id)
	reply, err := conn.Do("lrange", "history"+strconv.Itoa(id), 0, 4)
	rangeResult, _ := redis.Ints(reply, err)
	var goods []models.GoodsSKU
	for _, goodsId := range rangeResult {
		//根据id查找对应的SKU
		var sku models.GoodsSKU
		(*newOrm).QueryTable("GoodsSKU").Filter("Id", goodsId).One(&sku)
		goods = append(goods, sku)
	}
	c.Data["goods"] = goods
	beego.Info("rangeResult is", rangeResult)
	return nil
}
