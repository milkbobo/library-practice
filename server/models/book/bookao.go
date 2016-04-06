package book

import (
	// . "github.com/fishedee/language"
	. "library/models/common"
	// "time"
)

type BookAoModel struct {
	BaseModel
	BookDb BookDbModel
}

func (this *BookAoModel) Search(book Book, page CommonPage) Books {
	return this.BookDb.Search(book, page)
}

func (this *BookAoModel) Get(id int) Book {
	return this.BookDb.Get(id)
}

func (this *BookAoModel) Add(book Book) {
	// this.checkInput(book)

	this.BookDb.Add(book)
}

func (this *BookAoModel) Mod(id int, book Book) {
	// this.checkInput(book)

	this.BookDb.Mod(id, book)
}

func (this *BookAoModel) Del(id int) {

	this.BookDb.Del(id)

}
