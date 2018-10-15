package controllers

import (
	"github.com/astaxie/beego"
	"regexp"
	"github.com/astaxie/beego/orm"
	"dailyfresh/models"
	"github.com/astaxie/beego/utils"
	"strconv"
	"errors"
)

type RegisterController struct {
	beego.Controller
}

func (c *RegisterController) ShowRegister() {
	c.TplName = "register.html"
}

func (c *RegisterController) HandleRegister() {
	//获取数据
	userName := c.GetString("user_name")
	pwd := c.GetString("pwd")
	cpwd := c.GetString("cpwd")
	email := c.GetString("email")
	beego.Info(userName, pwd, cpwd, email)
	//校验数据
	if userName == "" || pwd == "" || cpwd == "" || email == "" {
		beego.Error("注册信息不正确,请重新输入")
		c.TplName = "register.html"
		c.Data["errmsg"] = "注册信息不正确,请重新输入"
		return
	}
	//校验密码
	if cpwd != pwd {
		beego.Error("两次密码输入不一致")
		c.TplName = "register.html"
		c.Data["errmsg"] = "两次密码输入不一致"
		return
	}
	//校验邮箱
	compile, _ := regexp.Compile("^[A-Za-z0-9\u4e00-\u9fa5]+@[a-zA-Z0-9_-]+(\\.[a-zA-Z0-9_-]+)+$")
	s := compile.FindString(email)
	if s == "" {
		beego.Error("邮箱格式不正确")
		c.TplName = "register.html"
		c.Data["errmsg"] = "邮箱格式不正确"
		return
	}
	//查询用户名是否重复
	newOrm := orm.NewOrm()
	var user models.User
	user.Name = userName
	user.PassWord = pwd
	user.Email = email
	readErr := newOrm.Read(&user, "Name")
	if readErr != orm.ErrNoRows {
		beego.Error("用户名已经存在")
		c.TplName = "register.html"
		c.Data["errmsg"] = "用户名已经存在"
		return
	}
	//如果以上的数据校验成功,进行处理数据
	_, err := newOrm.Insert(&user)
	if err != nil {
		beego.Error("数据插入失败")
		c.TplName = "register.html"
		c.Data["errmsg"] = "数据插入失败"
		return
	}
	err = sendRegisterEmail(&user)
	if err != nil {
		c.TplName = "register.html"
		return
	}
	//返回视图,注册成功之后进入登录界面
	c.Redirect("/login", 302)
}
func sendRegisterEmail(user *models.User) error {
	config := `{"username":"geeksi@163.com","password":"20200101geeksi","host":"smtp.163.com","port":25}`
	temail := utils.NewEMail(config)
	temail.To = []string{user.Email} //指定收件人邮箱地址，就是用户在注册时填写的邮箱地址
	temail.From = "geeksi@163.com"   //指定发件人的邮箱地址，这里我们使用的QQ邮箱。
	temail.Subject = "天天生鲜用户激活"      //指定邮件的标题
	/*指定邮件的内容。该内容发送到用户的邮箱中以后，该用户打开邮箱，可以将该URL地址复制到地址栏中，敲回车键，就会向该指定的URL地址发送请求，
	 我们在该地址对应的方法中，接收该用户的ID,然后根据该Id,查询出用户的信息后，将其对应的一个属性，Active设置为true,表明用户已经激活了，
	 那么用户就可以登录了。*/
	temail.HTML = "复制该连接到浏览器中激活：http://192.168.1.19:8080/active?id=" + strconv.Itoa(user.Id)
	//发送邮件
	sendErr := temail.Send()
	if sendErr != nil {
		beego.Error("注册邮件发送失败")
		return errors.New("注册邮件发送失败")
	}

	return nil

}
