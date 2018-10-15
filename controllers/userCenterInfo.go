package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"dailyfresh/models"
)

type UserCenterInfoController struct {
	beego.Controller
}

func (c *UserCenterInfoController) ShowUserInfo() {
	//展示用户信息表
	userName := c.GetSession("username")
	//查询用户信息
	newOrm := orm.NewOrm()
	var user models.User
	user.Name = userName.(string)
	newOrm.Read(&user, "Name")
	//多表查询,查询该用户的默认地址
	var address models.Address
	newOrm.QueryTable("Address").Filter("User__Id", user.Id).Filter("Isdefault", true).One(&address)
	beego.Info("address is", address)
	beego.Info("user",user)
	c.Data["username"]=userName
	c.Data["address"]=address.Addr
	c.Data["phoneNumber"]=address.Phone
	c.Layout="layout.html"
	c.TplName = "user_center_info.html"
}
