package chassis

import (
	"time"

	"gorm.io/gorm"
	"gorm.io/plugin/soft_delete"
)

//Page page
type Page struct {
	List   interface{} `json:"list,omitempty"`
	Total  int64        `json:"total,omitempty"`
	Offset int        `json:"offset,omitempty"`
	Index  int        `json:"page_index,omitempty"`
	Size   int        `json:"page_size,omitempty"`
	Pages  int        `json:"pages,omitempty"`
}

//Pagination 新建分页查询
type Pagination struct {
	Offset    uint        `json:"offset,omitempty"`
	Limit     uint        `json:"limit,omitempty"`
	Condition interface{} `json:"condition,omitempty"`
}

//NewPage new page
func newPage(data interface{}, index, size int, count int64) *Page {
	var pages int
	if count%int64(size) == 0 {
		pages = int(count / int64(size))
	} else {
		pages = int(count/int64(size) + 1)
	}
	return &Page{
		List:   data,
		Total:  count,
		Size:   size,
		Offset: index * size,
		Index:  index,
		Pages:  pages,
	}
}

//NewPagination pagination query
func NewPagination(db *gorm.DB, model interface{}, pageIndex, pageSize int) *Page {
	var count int64
	db.Count(&count)
	if count > 0 && count > int64(pageIndex*pageSize) {
		db.Limit(int(pageSize)).
			Offset(int(pageIndex * pageSize)).
			Find(model)
		return newPage(model, pageIndex, pageSize, count)
	}
	return nil
}

//SampleBaseDO model with id pk
type SampleBaseDO struct {
	ID uint `gorm:"primary_key" json:"id"`
}

//Model gorm model
type BaseDO struct {
	ID        uint       `gorm:"primary_key" json:"id"`            // primary key
	CreatedAt time.Time  `json:"created_at,omitempty"`             // created time
	UpdatedAt time.Time  `json:"updated_at,omitempty"`             //updated time
	DeletedAt soft_delete.DeletedAt //deleted time
}

//ComplexBaseDO gorm model composed Model add Addition
type ComplexBaseDO struct {
	BaseDO
	Version  uint   `json:"version"` //version opt lock
	Addition string `json:"addition,omitempty"`
}
