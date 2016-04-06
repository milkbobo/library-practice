package client

import (
	// . "github.com/fishedee/language
	. "library/models/common"
)

type ClientAoModel struct {
	BaseModel
	ClientDb ClientDbModel
}

func (this *ClientAoModel) Search(where Client, limit CommonPage) Clients {
	return this.ClientDb.Search(where, limit)
}

func (this *ClientAoModel) Get(id int) Client {
	return this.ClientDb.Get(id)
}

func (this *ClientAoModel) GetByIds(ids []int) []Client {
	return this.ClientDb.GetByIds(ids)
}

func (this *ClientAoModel) Add(data Client) int {
	return this.ClientDb.Add(data)
}

func (this *ClientAoModel) Mod(id int, data Client) {
	this.ClientDb.Mod(id, data)
}

func (this *ClientAoModel) Del(id int) {
	this.ClientDb.Del(id)
}

// func (this *ClientAoModel) AddOnce(data Client) int {
// 	if data.OpenId == "" {
// 		Throw(1, "不合法的openId"+data.OpenId)
// 	}
// 	clients := this.ClientDb.GetByOpenId(data.OpenId)
// 	if len(clients) != 0 {
// 		this.ClientDb.Mod(data.ClientId, data)
// 		return clients[0].ClientId
// 	} else {
// 		return this.ClientDb.Add(data)
// 	}
// }
