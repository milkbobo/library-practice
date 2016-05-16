package routers

import (
	. "github.com/fishedee/web"
	. "library/controllers"
)

func init() {
	InitRoute("/index", &MainController{})
}
