package chassis

import (
	"time"

	"gorm.io/gorm"
)

//Page page
type Page struct {
	List   interface{} `json:"list,omitempty"`
	Total  uint        `json:"total,omitempty"`
	Offset uint        `json:"offset,omitempty"`
	Index  uint        `json:"page_index,omitempty"`
	Size   uint        `json:"page_size,omitempty"`
	Pages  uint        `json:"pages,omitempty"`
}

//Pagination 新建分页查询
type Pagination struct {
	Offset    uint        `json:"offset,omitempty"`
	Limit     uint        `json:"limit,omitempty"`
	Condition interface{} `json:"condition,omitempty"`
}

//NewPage new page
func newPage(data interface{}, index, size, count uint) *Page {
	var pages uint
	if count%size == 0 {
		pages = count / size
	} else {
		pages = count/size + 1
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
func NewPagination(db *gorm.DB, model interface{}, pageIndex, pageSize uint) *Page {
	var count int64
	db.Count(&count)
	if count > 0 && count > int64(pageIndex*pageSize) {
		db.Limit(int(pageSize)).
			Offset(int(pageIndex * pageSize)).
			Find(model)
		return newPage(model, pageIndex, pageSize, uint(count))
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
	DeletedAt *time.Time `sql:"index" json:"deleted_at,omitempty"` //deleted time
}

//ComplexBaseDO gorm model composed Model add Addition
type ComplexBaseDO struct {
	BaseDO
	Version  uint   `json:"version"` //version opt lock
	Addition string `json:"addition,omitempty"`
}
