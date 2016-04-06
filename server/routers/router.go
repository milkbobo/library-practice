package routers

import (
	. "library/controllers"
)

func init() {
	InitRoute("/index", &MainController{})
}
