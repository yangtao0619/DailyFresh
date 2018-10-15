package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"dailyfresh/models"
)

type ActiveController struct {
	beego.Controller
}

func (c *ActiveController) ActiveUser() {
	//处理注册的逻辑
	id, err := c.GetInt("id")
	if err != nil {
		beego.Error("get id err:", err)
		c.TplName = "register.html"
		return
	}
	beego.Info("get user id", id)
	//查询用户,并写入数据
	newOrm := orm.NewOrm()
	var user models.User
	user.Id = id
	readErr := newOrm.Read(&user)
	if readErr != nil {
		beego.Error("参数不正确")
		c.TplName = "register.html"
		return
	}
	user.Active = true
	_, err = newOrm.Update(&user)
	if err != nil {
		beego.Error("激活失败")
		return
	}
	//激活成功,重定向到登录界面
	c.Redirect("/login",302)
}
