package info

import (
	"math"
)

// 分页数据
type Page struct {
	CurPage     int // 当前页码
	TotalPage   int // 总页
	PerPageSize int
	Count       int // 总记录数
	List        interface{}
}

func NewPage(page, perPageSize, count int, list interface{}) Page {
	totalPage := 0
	if count > 0 {
		totalPage = int(math.Ceil(float64(count) / float64(perPageSize)))
	}
	return Page{page, totalPage, perPageSize, count, list}
}
