package main

import (
	_ "dailyfresh/routers"
	"github.com/astaxie/beego"
	_ "dailyfresh/models"
)

func main() {
	beego.Run()
}

