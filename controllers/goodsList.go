package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"dailyfresh/models"
)

type GoodsListController struct {
	beego.Controller
}

func (c *GoodsListController) ShowGoodsList() {
	//展示类型列表数据
	ShowTypeAndGetUserInfo(&c.Controller)
	GetUser(&c.Controller)
	//获取typeId执行排序
	id, err := c.GetInt("id")
	if err != nil {
		beego.Error("获取类型id失败")
		c.TplName = "list.html"
		return
	}
	//获取排序字段
	sort := c.GetString("sort")
	beego.Info("sort name is", sort)
	//获取typeid,查询最近上架的两条数据
	newOrm := orm.NewOrm()
	showNewGoods(c, &newOrm, id)
	//获取所有的分类商品数据进行显示
	showCurrentPageGoods(&newOrm, c, id, sort)
	c.TplName = "list.html"
}
func showCurrentPageGoods(ormer *orm.Ormer, c *GoodsListController, id int, sort string) {
	//根据类型找出所有的sku
	var allSku []models.GoodsSKU
	if sort == "" {
		//默认排序
		(*ormer).QueryTable("GoodsSKU").RelatedSel("Goods", "GoodsType").
			Filter("GoodsType__Id", id).All(&allSku)
	} else if sort == "price" {
		//按价格排序
		(*ormer).QueryTable("GoodsSKU").RelatedSel("Goods", "GoodsType").
			Filter("GoodsType__Id", id).OrderBy("Price").All(&allSku)
	} else if sort == "popularity" {
		//按人气排序
		(*ormer).QueryTable("GoodsSKU").RelatedSel("Goods", "GoodsType").
			Filter("GoodsType__Id", id).OrderBy("Sales").All(&allSku)
	}
	beego.Info("all sku in this type are", allSku)
	c.Data["allSku"] = allSku
	c.Data["sort"]=sort
}
func showNewGoods(c *GoodsListController, newOrm *orm.Ormer, typeId int) {
	var newGoods []models.GoodsSKU
	(*newOrm).QueryTable("GoodsSKU").RelatedSel("Goods", "GoodsType").Filter("GoodsType__Id", typeId).
		Limit(2, 0).OrderBy("Time").All(&newGoods)
	beego.Info("newGoods are", newGoods)
	c.Data["typeId"] = typeId
	c.Data["newGoods"] = newGoods
}
