package book

import (
	. "github.com/fishedee/language"
	. "library/models/common"
	"strconv"
)

type BookDbModel struct {
	BaseModel
}

func (this *BookDbModel) Search(where Book, limit CommonPage) Books {
	db := this.DB.NewSession()
	defer db.Close()

	if limit.PageIndex == 0 && limit.PageSize == 0 {
		return Books{
			Count: 0,
			Data:  []Book{},
		}
	}

	if where.Bid != 0 {
		db = db.And("bid = ?", where.Bid)
	}
	if where.Bname != "" {
		db = db.And("bname like ?", "%"+where.Bname+"%")
	}

	data := []Book{}
	var err error
	err = db.OrderBy("createTime desc").Limit(limit.PageSize, limit.PageIndex).Find(&data)
	if err != nil {
		panic(err)
	}

	count, err := db.Count(&where)
	if err != nil {
		panic(err)
	}

	return Books{
		Count: int(count),
		Data:  data,
	}
}
func (this *BookDbModel) Get(id int) Book {
	var books []Book
	err := this.DB.Where("bid = ?", id).Find(&books)
	if err != nil {
		panic(err)
	}
	if len(books) == 0 {
		Throw(1, "不存在该数据"+strconv.Itoa(id))
	}
	return books[0]
}

func (this *BookDbModel) Add(book Book) {
	_, err := this.DB.Insert(&book)
	if err != nil {
		panic(err)
	}
}

func (this *BookDbModel) Mod(id int, book Book) {
	_, err := this.DB.Where("bid=?", id).Update(&book)
	if err != nil {
		panic(err)
	}
}

func (this *BookDbModel) Del(id int) {
	var book Book
	_, err := this.DB.Where("Bid = ?", id).Delete(&book)
	if err != nil {
		panic(err)
	}
}
