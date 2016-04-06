package book

import (
	"time"
)

type Book struct {
	Bid        int `xorm:"autoincr"`
	Bname      string
	Username   string
	CreateTime time.Time `xorm:"created"`
	ModifyTime time.Time `xorm:"updated"`
}

type Books struct {
	Count int
	Data  []Book
}
