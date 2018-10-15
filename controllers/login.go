package controllers

import (
	"github.com/astaxie/beego"
	"dailyfresh/models"
	"github.com/astaxie/beego/orm"
	"time"
	"encoding/base64"
)

type LoginController struct {
	beego.Controller
}

func (c *LoginController) ShowLogin() {
	//获取数据,校验数据
	cookieName := c.Ctx.GetCookie("username")
	name, _ := base64.StdEncoding.DecodeString(cookieName)
	cookieName = string(name)
	if cookieName == "" {
		c.Data["checked"] = ""
		c.Data["username"] = ""
	} else {
		c.Data["checked"] = "checked"
		c.Data["username"] = cookieName

	}
	c.TplName = "login.html"
}

func (c *LoginController) HandleLogin() {
	//获取数据,校验数据
	username := c.GetString("username")
	pwd := c.GetString("pwd")
	checked := c.GetString("checked")
	beego.Info(username, pwd, checked)
	if username == "" || pwd == "" {
		beego.Error("用户名或者密码不能为空")
		c.Data["errmsg"] = "用户名或者密码不能为空"
		return
	}
	//查询用户是否存在
	var user models.User
	newOrm := orm.NewOrm()
	user.Name = username
	readErr := newOrm.Read(&user, "Name")
	if readErr == orm.ErrNoRows {
		beego.Error("用户名不存在")
		c.Data["errmsg"] = "用户名不存在"
		c.TplName="login.html"
		return
	}
	//用户名存在的时候,比较密码
	if user.PassWord != pwd {
		beego.Error("密码错误!")
		c.Data["errmsg"] = "密码错误,请重新输入!"
		c.TplName="login.html"
		return
	}
	//用户名和密码都正确的时候,判断用户是否激活
	if user.Active != true {
		beego.Error("该用户没有激活!")
		c.Data["errmsg"] = "该用户没有激活,请激活后重试!"
		c.TplName="login.html"
		return
	}
	//如果都校验成功了,需要检测是否需要记住用户名
	if checked == "on" {
		//写入cookie
		c.Ctx.SetCookie("username", base64.StdEncoding.EncodeToString([]byte(username)), time.Second*3600)
	} else {
		c.Ctx.SetCookie("username", "", -1)
	}
	//将用户名存到session
	c.SetSession("username", username)
	c.Redirect("/", 302)
}
