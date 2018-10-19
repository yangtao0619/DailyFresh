package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"dailyfresh/models"
)

type IndexController struct {
	beego.Controller
}

func GetUser(c *beego.Controller) string {
	userName := c.GetSession("username")
	var name string
	if userName == nil {
		name = ""
	} else {
		name = userName.(string)
	}
	c.Data["username"] = name
	return name
}

//展示首页的数据
func (c *IndexController) ShowIndex() {
	GetUser(&c.Controller)
	//需要查询四张表,分别是促销表,轮播商品表,首页分类商品展示表和商品类型表
	newOrm := orm.NewOrm()
	var goodsTypes []models.GoodsType
	newOrm.QueryTable("GoodsType").All(&goodsTypes)
	c.Data["goodsTypes"] = goodsTypes
	//查询促销表
	var promotionGoods []models.IndexPromotionBanner
	newOrm.QueryTable("IndexPromotionBanner").OrderBy("Index").All(&promotionGoods)
	c.Data["promotionGoods"] = promotionGoods
	//查询轮播表
	var bannerGoods []models.IndexGoodsBanner
	newOrm.QueryTable("IndexGoodsBanner").OrderBy("Index").All(&bannerGoods)
	c.Data["bannerGoods"] = bannerGoods
	beego.Info(bannerGoods)
	//需要拿到对应的类型名称来查询数据,根据类型分别查找图片商品和文字商品
	goodsSlice := make([]map[string]interface{}, len(goodsTypes))
	for index, goodsType := range goodsTypes {
		goodsMap := make(map[string]interface{})
		goodsMap["type"] = goodsType
		goodsSlice[index] = goodsMap
	}
	//根据type查找对应的IndexTypeGoodsBanner
	//用这个type查找
	var textGoods []models.IndexTypeGoodsBanner
	var picGoods []models.IndexTypeGoodsBanner
	for _, goodMap := range goodsSlice {
		goodType := goodMap["type"]
		newOrm.QueryTable("IndexTypeGoodsBanner").RelatedSel("GoodsType", "GoodsSKU").
			Filter("GoodsType", goodType).Filter("DisplayType", 0).OrderBy("Index").All(&textGoods)
		newOrm.QueryTable("IndexTypeGoodsBanner").RelatedSel("GoodsType", "GoodsSKU").
			Filter("GoodsType", goodType).Filter("DisplayType", 1).OrderBy("Index").All(&picGoods)
		goodMap["textGoods"] = textGoods
		goodMap["picGoods"] = picGoods
	}

	c.Data["goods"] = goodsSlice

	c.Layout = "userCenterLayout.html"
	c.TplName = "index.html"
}
