package controllers

import (
	"github.com/astaxie/beego"
	"dailyfresh/models"
	"github.com/astaxie/beego/orm"
)

type SearchGoodsController struct {
	beego.Controller
}

func (c *SearchGoodsController) HandleSearchName() {
	GetUser(&c.Controller)
	searchName := c.GetString("searchName")
	beego.Info("searchName is", searchName)
	//拿到搜索的关键字之后搜索数据库,将满足条件的数据返回给视图
	var searchResults []models.GoodsSKU
	newOrm := orm.NewOrm()
	newOrm.QueryTable("GoodsSKU").Filter("Name__icontains", searchName).All(&searchResults)
	//将搜索到的数据返回给视图
	c.Data["searchResults"] = searchResults
	ShowTypeAndGetUserInfo(&c.Controller)
	c.TplName = "searchResult.html"
}
