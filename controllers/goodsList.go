package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"dailyfresh/models"
	"math"
	"github.com/gomodule/redigo/redis"
)

type GoodsListController struct {
	beego.Controller
}

func (c *GoodsListController) ShowGoodsList() {
	//展示类型列表数据
	ShowTypeAndGetUserInfo(&c.Controller)
	userName := GetUser(&c.Controller)
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
	count := showCartNum(err, userName, newOrm, c)
	c.Data["goodsCount"] = count
	c.TplName = "list.html"
}

func showCartNum(err error, userName string, newOrm orm.Ormer, c *GoodsListController) int {
	conn, err := redis.Dial("tcp", "192.168.1.19:6379")
	if err != nil {
		beego.Error("redis连接错误")
		c.Redirect("/", 302)
		return 0
	}
	defer conn.Close()
	var user models.User
	user.Name = userName
	readErr := newOrm.Read(&user, "Name")
	if readErr != nil {
		beego.Error("查询数据失败")
		return 0
	}
	count := GetGoodsCount(conn, user.Id)
	return count
}
func showCurrentPageGoods(ormer *orm.Ormer, c *GoodsListController, id int, sort string) {
	//根据类型找出所有的sku
	var allSku []models.GoodsSKU
	var count int64
	var qs orm.QuerySeter
	if sort == "" {
		//默认排序
		qs = (*ormer).QueryTable("GoodsSKU").RelatedSel("Goods", "GoodsType").
			Filter("GoodsType__Id", id)
		count, _ = qs.Count()
	} else if sort == "price" {
		//按价格排序
		qs = (*ormer).QueryTable("GoodsSKU").RelatedSel("Goods", "GoodsType").
			Filter("GoodsType__Id", id).OrderBy("Price")
		count, _ = qs.Count()
	} else if sort == "popularity" {
		//按人气排序
		qs = (*ormer).QueryTable("GoodsSKU").RelatedSel("Goods", "GoodsType").
			Filter("GoodsType__Id", id).OrderBy("Sales")
		count, _ = qs.Count()
		//Filter("GoodsType__Id", id).OrderBy("Sales").All(&allSku)
	}
	//查询出来的数据需要进行分页
	pageSize := 2
	pageCount := math.Ceil(float64(count) / float64(pageSize))
	beego.Info("pageCount is", pageCount)
	pageIndex, err := c.GetInt("pageIndex", 1)
	beego.Info("pageIndex is", pageIndex)
	if err != nil {
		beego.Error("获取当前页码失败,err:", err)
	}
	pageDataMap := listPageTool(int(pageCount), pageIndex)
	pageDataMap["pageIndex"] = pageIndex
	qs.Limit(pageSize, (pageIndex-1)*pageSize).All(&allSku)
	beego.Info("pageBuffer is", pageDataMap["pageBuffer"].([]int), "pageindex is", pageIndex)
	sendDataToView(c, allSku, sort, pageDataMap)
}

/*
传递数据给视图
 */
func sendDataToView(c *GoodsListController, allSku []models.GoodsSKU, sort string, pageDataMap map[string]interface{}) {
	c.Data["allSku"] = allSku
	c.Data["sort"] = sort
	c.Data["pageBuffer"] = pageDataMap["pageBuffer"].([]int)
	c.Data["prePage"] = pageDataMap["prePage"]
	c.Data["nextPage"] = pageDataMap["nextPage"]
	c.Data["pageIndex"] = pageDataMap["pageIndex"]
	c.Data["sort"] = sort
}
func showNewGoods(c *GoodsListController, newOrm *orm.Ormer, typeId int) {
	var newGoods []models.GoodsSKU
	(*newOrm).QueryTable("GoodsSKU").RelatedSel("Goods", "GoodsType").Filter("GoodsType__Id", typeId).
		Limit(2, 0).OrderBy("Time").All(&newGoods)
	beego.Info("newGoods are", newGoods)
	c.Data["typeId"] = typeId
	c.Data["newGoods"] = newGoods
}
