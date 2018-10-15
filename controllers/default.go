package controllers

import (
	"github.com/astaxie/beego"
)

type MainController struct {
	beego.Controller
}

func (c *MainController) Get() {
	userName := c.GetSession("username")
	c.Data["username"]=userName
	c.Layout="layout.html"
	c.TplName = "index.html"
}
