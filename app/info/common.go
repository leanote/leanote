package info

import (
)


// 分页数据
type Page struct {
	CurPage int // 当前页码
	TotalPage int // 总页
	Count int // 总记录数
	List interface{}
}