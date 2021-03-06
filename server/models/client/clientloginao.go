package client

import (
	"crypto/rand"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	. "github.com/fishedee/language"
	. "library/models/common"
	"strings"
	// "time"
)

type ClientLoginAoModel struct {
	BaseModel
	ClientAo ClientAoModel
	ClientDb ClientDbModel
}

func (this *ClientLoginAoModel) Login(client Client) {
	sess, err := this.Session.SessionStart()
	if err != nil {
		panic("session启动失败")
	}
	defer sess.SessionRelease()

	v := this.ClientAo.GetByName(client.Username)

	fmt.Println("userData")
	fmt.Printf("%+v", v)

	if len(v) == 0 {
		Throw(1, "用户名不存在")
		return
	}

	hashAndSalt := strings.Split(v[0].Password, ":")
	password := hashAndSalt[0]
	salt := hashAndSalt[1]
	hash := sha1.New()
	passwordSha1Byte := hash.Sum([]byte(client.Password + salt))
	passwordSha1 := hex.EncodeToString(passwordSha1Byte)

	if password != passwordSha1 {
		Throw(1, "密码错误")

	}

	sess.Set("clientId", v[0].ClientId)

}

func (this *ClientLoginAoModel) Logout() {
	sess, err := this.Session.SessionStart()
	if err != nil {
		panic("session启动失败！")
	}
	defer sess.SessionRelease()

	sess.Set("clientId", 0)
}

func (this *ClientLoginAoModel) CheckMustLogin() Client {
	client := this.IsLogin()
	if client.ClientId == 0 {
		Throw(10001, "用户未登录！")
	}
	return client
}

func (this *ClientLoginAoModel) IsLogin() Client {
	sess, err := this.Session.SessionStart()
	if err != nil {
		panic("session启动失败")
	}
	defer sess.SessionRelease()

	clientId := sess.Get("clientId")
	clientIdInt, ok := clientId.(int)

	fmt.Println("clientId", clientId, "clientIdInt", clientIdInt)

	if ok && clientIdInt >= 10000 {
		return this.ClientAo.Get(clientIdInt)
	} else {
		return Client{}
	}
}

func (this *ClientLoginAoModel) Register(username, password, password2 string) {

	// sessDb := this.DB.NewSession()
	// defer sessDb.Close()
	// sessDb.Begin()

	if password != password2 {
		Throw(1, "确认密码不正确")
		return
	}

	v := this.ClientAo.GetByName(username)
	// v := this.ClientDb.GetByNameForTrans(sessDb, username)

	// time.Sleep(100 * time.Millisecond)

	if len(v) > 0 {
		Throw(1, "用户名已存在，请重新注册其他用户名字")
		return
	}

	//生成随机字符串
	k := make([]byte, 5)
	if _, err := rand.Read(k); err != nil {
		panic(err)
	}

	salt := hex.EncodeToString(k)

	// fmt.Println("salt", salt)
	hash := sha1.New()
	passwordSha1Byte := hash.Sum([]byte(password + salt))
	passwordSha1 := hex.EncodeToString(passwordSha1Byte) + ":" + salt

	// fmt.Println("passwordSha1", passwordSha1)

	// clientId := this.ClientDb.AddForTrans(sessDb, Client{
	// Username: username,
	// Password: passwordSha1,
	// })

	clientId := this.ClientDb.AddOnce(Client{
		Username: username,
		Password: passwordSha1,
	})

	// sessDb.Commit()

	sess, err := this.Session.SessionStart()
	if err != nil {
		panic("session启动失败！")
	}

	defer sess.SessionRelease()

	sess.Set("clientId", clientId)

}
