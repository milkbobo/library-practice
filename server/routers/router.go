package routers

import (
	. "library/controllers"
)

func init() {
	InitRoute("/main", &MainController{})
}
