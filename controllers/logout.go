package controllers

import "github.com/astaxie/beego"

type LogoutController struct {
	beego.Controller
}

func (c *LogoutController) Logout() {
	//登出的时候要删除session
	c.DelSession("username")
	c.Redirect("/login", 302)
}
