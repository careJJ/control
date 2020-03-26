package main

import (
	_ "shoptest/routers"
	"github.com/astaxie/beego"
	_"shoptest/models"
)

func main() {
	beego.Run()
}

