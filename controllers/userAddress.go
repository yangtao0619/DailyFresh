package controllers

import (
	"github.com/astaxie/beego"
	"dailyfresh/models"
	"github.com/astaxie/beego/orm"
	"errors"
)

type UserCenterAddrController struct {
	beego.Controller
}

func (c *UserCenterAddrController) ShowAddAddress() {
	username := GetUser(&c.Controller)
	//获取数据,将默认地址填充到页面上
	//北京市 海淀区 东北旺西路8号中关村软件园 （李思 收） 182****7528
	var user models.User
	user.Name = username
	newOrm := orm.NewOrm()
	newOrm.Read(&user, "Name")

	//查询地址表中id=user.Id的记录,只有一条
	var address models.Address
	newOrm.QueryTable("Address").Filter("User__Id", user.Id).Filter("Isdefault", true).One(&address)
	c.Data["address"] = address.Addr + "  (" + address.Receiver + " 收)  " + address.Phone
	beego.Info(address, user.Id)
	c.Layout = "userCenterLayout.html"
	c.TplName = "user_center_site.html"
}

func (c *UserCenterAddrController) HandleAddAddress() {
	//获取数据
	userName := c.GetSession("username")
	receiver := c.GetString("receiver")
	detailAddr := c.GetString("detailAddr")
	mailCode := c.GetString("mailCode")
	phoneNumber := c.GetString("phoneNumber")
	//校验数据
	if receiver == "" || detailAddr == "" || mailCode == "" || phoneNumber == "" {
		beego.Error("输入项不能为空")
		c.Data["errMsg"] = "输入项不能为空,请检查后重新提交"
		c.Redirect("/user/addDefaultAddr",302)
		return
	}
	//处理数据,将数据写入数据库中
	var user models.User
	var address models.Address
	user.Name = userName.(string)
	newOrm := orm.NewOrm()
	newOrm.Read(&user, "Name")
	address.Phone = phoneNumber
	address.Receiver = receiver
	address.Zipcode = mailCode
	address.Addr = detailAddr
	address.User = &user
	address.Isdefault = true
	//需要将原来的默认地址改成非默认的
	err := changeDefaultToFalse(newOrm, user)
	if err != nil {
		beego.Error(err)
		c.Data["errMsg"] = err
		c.TplName = "user_center_site.html"
		c.Layout = "userCenterLayout.html"
		return
	}
	_, err = newOrm.Insert(&address)
	if err != nil {
		beego.Error("数据插入失败")
		c.Data["errMsg"] = "数据插入失败,请稍后重试"
		c.Layout = "userCenterLayout.html"
		c.TplName = "user_center_site.html"
		return
	}
	beego.Info("插入地址数据成功")
	//返回视图
	c.Redirect("/user/addDefaultAddr", 302)
}
func changeDefaultToFalse(newOrm orm.Ormer, user models.User) error {
	var address models.Address
	newOrm.QueryTable("Address").Filter("User__Id", user.Id).Filter("Isdefault", true).One(&address)
	address.Isdefault = false
	_, err := newOrm.Update(&address, "Isdefault")
	if err != nil {
		return errors.New("重置默认地址失败")
	}
	return nil
}
