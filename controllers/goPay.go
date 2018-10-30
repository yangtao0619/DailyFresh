package controllers

import (
	"github.com/astaxie/beego"
	"github.com/smartwalle/alipay"
	"fmt"
)

type PayController struct {
	beego.Controller
}

func (c *PayController) GoPay() {
	//首先判断支付方式
	payMethod, _ := c.GetInt("payMethod")
	beego.Info("payMethod is", payMethod)
	switch payMethod {
	case 0:
		//货到付款
		fallthrough
	case 1:
		fallthrough
		//银联
	case 3:
		//支付宝
		aliPay(c)
	case 2:
		//微信
		tecentPay()
	}

}

//微信支付
func tecentPay() {

}
func aliPay(c *PayController) {

	var aliPublicKey = "MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKC" +
		"AQEA3H4ZenEmE+UWgsWM/kitKd5CfOP54aitNHn2CQIbATDStXmvQfg" +
		"j2oSdV/p3LvPJGN8RSZeUBusfrRYtvbfg4WNNdA6HJ4twciqOQS8U/Q" +
		"xPp7aUvE2wUVENzD/Dy4a1ZHC3+K2lIGr7NwaVncYEFGZnexN5lXwjp" +
		"VqwXSKNjsT3Jc3hTTdpto4UpDCmm89w7QuGOU2NfFE5CnKwSW9OfUqc9" +
		"/zqA6ZuQeov2nFPkmPdbE2P6QGALulg4zz2FohQYxB0E61slFc70NeVS" +
		"a5tkTfiXEBPspQONJjbEiZYb89FvoStpxL+DuyIFAxHEj2sPDYu435ek" +
		"zVdxAoP8jWsEQIDAQAB"
	// 可选，支付宝提供给我们用于签名验证的公钥，通过支付宝管理后台获取
	var privateKey = "MIIEpQIBAAKCAQEA3H4ZenEmE+UWgsWM/kitKd5CfOP54aitNHn2CQIbATDStXmv" +
		"Qfgj2oSdV/p3LvPJGN8RSZeUBusfrRYtvbfg4WNNdA6HJ4twciqOQS8U/QxPp7aU" +
		"vE2wUVENzD/Dy4a1ZHC3+K2lIGr7NwaVncYEFGZnexN5lXwjpVqwXSKNjsT3Jc3h" +
		"TTdpto4UpDCmm89w7QuGOU2NfFE5CnKwSW9OfUqc9/zqA6ZuQeov2nFPkmPdbE2P" +
		"6QGALulg4zz2FohQYxB0E61slFc70NeVSa5tkTfiXEBPspQONJjbEiZYb89FvoSt" +
		"pxL+DuyIFAxHEj2sPDYu435ekzVdxAoP8jWsEQIDAQABAoIBAQCC0WIGbklcNmgz" +
		"sEelurLai17BQHVKOEyDPPUHhTNGcpQhTY/4wONsy4+a2iSKO+ONGRPlqMQPksKZ" +
		"a/Y5gHYw4zzZ5aC0ipttcOgzrl5ygDJmXAJE8obwx/k6vH5LK6JFdEcCiOvWnwJr" +
		"NEHieNCE1fkBYZ2aXiu7+GF48H4yPHnuxDLTtPTAQVHKsEniLApulc+lio9R+9cg" +
		"qPFzdlyY4qlkgqKUbtFcJdKGdPR7nvythmKcat2x6PgrCiIVnDgY9GVXIRMdSuMi" +
		"rTS7qQoahp0phUKQxPzArOCi8KNgMzCUFwSOPlrhoVQlZeTrVtD1kZkcFwW0oCHu" +
		"vW13s6sdAoGBAPmTtxa2KEW/pZShyqOWQdt2S1QX42OZIKOZE4aYOLe0sBpeo+aY" +
		"5Ws+C55MJ9lmgGoXR0L4/94yYHYcon70K9hLOhCh6nM7o+lQ9B6AtZgtyaaDJeYp" +
		"2WZGsefZMxtuyR4u3nNLYhKZDfgZFDp1dZr+if2B7gElrnOHUc+h1qLzAoGBAOIq" +
		"xIq9kB18DZLLBWqTe1D2GcSDfIbjLPExMR1Z9xm6bvMVQXL63zTqFEnUV7Pv9zn4" +
		"YE65GbjrZIcR4ZVWB3/z/eBgSrtCIYos7gI2UyH2I/af4U9WNSN4kS/O7AUzSA3F" +
		"NH6E4wkFo75l+C947yWrBtrMhIxcIH/c2ZBs6U3rAoGBANfHhOaX+13CgpBtCdxB" +
		"zxLFxf8g4DJ+dB+9+4nFFlSOXiuOY7q9uqzr6fOk+FcYLjKLics3qVEc0RWNUFjf" +
		"FwFcmQlEVIXorKDOoyG0Ok0mWVAj16KV0CaDPNGtkmhHco8sCpw4MsTNm6xDUp/w" +
		"agvlwrxxl6taPugXuP4BeKdNAoGAXrZUSlmqIX7S3Fdi9EfAy53UGqSJoJ9AMd1M" +
		"2SLUxRR65BdRqkn+8VTZnDVtaPAkE0W9ZxpC+FqzZZEKbBRz3ZSbC7ynbxX5n7jD" +
		"D3Aajk1asCwyGZxbcnhKLMA1vNPF5+Ze3mDeBugys0hWj+LQG3Es1LHzDCiEf6dI" +
		"ASBq73MCgYEAvY4YbKIHpOk8EEOhP4Ij2fZJkdJoDpouYJehm4VJB+YMs0SVbYF2" +
		"EeBle1KIG7PwWd74BGz86vF+nY+P88fa6lnA7aOER4SI/be0O7rzXEtmlC1fzhv5" +
		"m5Tw2G64PkmMs/QfBPqdaxxOiBU8HN1QCPsnA4tLd2kDZaiwxbByhq8="
	// 必须，上一步中使用 RSA签名验签工具 生成的私钥
	var client = alipay.New("2016092000557229", aliPublicKey, privateKey, false)
	orderId := c.GetString("orderId")
	totalPrice := c.GetString("totalPrice")
	var p = alipay.AliPayTradePagePay{}
	p.NotifyURL = "http://192.168.1.19:8088/user/payOk"
	p.ReturnURL = "http://192.168.1.19:8088/user/payOk"
	p.Subject = "天天生鲜支付"
	p.OutTradeNo = orderId
	p.TotalAmount = totalPrice
	p.ProductCode = "FAST_INSTANT_TRADE_PAY"

	var url, err = client.TradePagePay(p)
	if err != nil {
		fmt.Println(err)
	}

	var payURL = url.String()
	fmt.Println(payURL)
	// 这个 payURL 即是用于支付的 URL，可将输出的内容复制，到浏览器中访问该 URL 即可打开支付页面。
	c.Redirect(payURL, 302)
}
