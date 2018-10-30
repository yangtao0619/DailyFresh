package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"dailyfresh/models"
)

func ShowTypeAndGetUserInfo(c *beego.Controller) {
	//查找类型数据
	var types []models.GoodsType
	ormer := orm.NewOrm()
	ormer.QueryTable("GoodsType").All(&types)
	c.Data["types"] = types
	c.Layout = "goodLayout.html"
}

func listPageTool(pageCount int, pageIndex int) (pageData map[string]interface{}) {
	//根据传入的总页码和当前页码,返回页码切片,上一页和下一页页码,设定显示的总页码数为5
	pageData = make(map[string]interface{})
	pageBuffer := make([]int, 5)
	if pageCount <= 5 {
		pageBuffer = make([]int, pageCount)
		for i := 1; i <= pageCount; i++ {
			pageBuffer[i-1] = i
		}
	} else if pageIndex < 3 {
		pageBuffer = []int{1, 2, 3, 4, 5}
	} else if pageIndex > pageCount-2 {
		pageBuffer = []int{pageIndex - 4, pageIndex - 3, pageIndex - 2, pageIndex - 1, pageIndex}
	} else {
		pageBuffer = []int{pageIndex - 2, pageIndex - 1, pageIndex, pageIndex + 1, pageIndex + 2}
	}
	prePage := pageIndex - 1
	if prePage <= 1 {
		prePage = 1
	}

	nextPage := pageIndex + 1
	if nextPage >= pageCount {
		nextPage = pageCount
	}
	pageData["pageBuffer"] = pageBuffer
	pageData["prePage"] = prePage
	pageData["nextPage"] = nextPage
	return
}
