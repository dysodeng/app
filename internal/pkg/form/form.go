package form

// Id id
type Id struct {
	Id uint64 `form:"id" json:"id"`
}

// Ids id列表
type Ids struct {
	Ids []uint64 `form:"ids" json:"ids"`
}

// SearchPagination 搜索分页
type SearchPagination struct {
	Pagination
	Keyword
}

// Keyword 关键词
type Keyword struct {
	Keyword string `form:"keyword" json:"keyword"`
}

// Pagination 分页
type Pagination struct {
	Page     int  `form:"page" json:"page"`
	PageSize int  `form:"page_size" json:"page_size"`
	NoPage   bool `form:"no_page" json:"no_page"`
}

func (pagination *Pagination) Valid() bool {
	if pagination.NoPage {
		return true
	}
	if pagination.Page <= 0 || pagination.PageSize <= 0 {
		return false
	}
	return true
}

func (pagination *Pagination) CheckOrDefault() {
	if pagination.Page <= 0 {
		pagination.Page = 1
	}
	if pagination.PageSize <= 0 {
		pagination.PageSize = 10
	}
}

func (pagination *Pagination) Default() {
	pagination.Page = 1
	pagination.PageSize = 10
}

func (pagination *Pagination) Offset() int {
	return (pagination.Page - 1) * pagination.PageSize
}
