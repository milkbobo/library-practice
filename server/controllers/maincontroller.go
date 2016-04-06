package controllers

import (
	// . "github.com/fishedee/language"
	// . "github.com/fishedee/web"
	"fmt"
	. "library/models/book"
	. "library/models/client"
	. "library/models/common"
)

type MainController struct {
	BaseController
	ClientLoginAo ClientLoginAoModel
	ClientAo      ClientAoModel
	BookAo        BookDbModel
}

//主页
func (this *MainController) Index_index() interface{} {

	//检查权限
	this.ClientLoginAo.CheckMustLogin()

	//业务逻辑
	books := this.BookAo.Search(Book{}, CommonPage{
		PageSize:  50,
		PageIndex: 0,
	})
	fmt.Printf("%+v", books)
	return books
}

//注册页面
func (this *MainController) Signup_register() interface{} {

	return 0
}

//注册操作
func (this *MainController) Register_redirect() interface{} {

	//检查输入
	var client struct {
		Username  string
		Password  string
		Password2 string
	}
	this.CheckPost(&client)

	fmt.Println(client)

	//用户注册
	this.ClientLoginAo.Register(client.Username, client.Password, client.Password2)

	return "/index/index"

}

//登陆页面
func (this *MainController) Signin_login() interface{} {

	return 0
}

//登陆操作
func (this *MainController) Login_redirect() interface{} {

	//检查输入
	var client Client
	this.CheckPost(&client)

	//登录
	this.ClientLoginAo.Login(client)

	return "/index/index"

}

//登出操作
func (this *MainController) Out_login() interface{} {

	this.ClientLoginAo.Logout()

	return 0

}

//添加书本页面
func (this *MainController) Addbook_add() interface{} {
	return 0
}

//添加书本操作
func (this *MainController) Add_redirect() interface{} {
	//检查输入
	book := Book{}
	this.CheckPost(&book)

	//检查权限
	this.ClientLoginAo.CheckMustLogin()

	//业务逻辑
	this.BookAo.Add(book)

	return "/index/index"

}

//修改页面
func (this *MainController) Alter_edit() interface{} {

	//检查输入
	book := Book{}
	this.CheckGet(&book)

	//检查权限
	this.ClientLoginAo.CheckMustLogin()

	//业务逻辑
	books := this.BookAo.Get(book.Bid)
	fmt.Printf("%+v", books)

	return books
}

//修改操作
func (this *MainController) Edit_redirect() interface{} {

	//检查输入
	book := Book{}
	this.CheckPost(&book)

	//登录
	this.ClientLoginAo.CheckMustLogin()

	//业务逻辑
	this.BookAo.Mod(book.Bid, book)

	return "/index/index"

}

//删除操作
func (this *MainController) Del_redirect() interface{} {

	//检查输入
	book := Book{}
	this.CheckGet(&book)

	//登录
	this.ClientLoginAo.CheckMustLogin()

	//业务逻辑
	this.BookAo.Del(book.Bid)

	return "/index/index"

}
