package tool

import (
	"gorm.io/gorm"
)

// PaginationRequest 分页请求参数
type PaginationRequest struct {
	Page     int    `json:"page" form:"page" binding:"omitempty,min=1"`                   // 页码，默认为1
	PageSize int    `json:"page_size" form:"page_size" binding:"omitempty,min=1,max=100"` // 每页数量，默认为10，最大100
	OrderBy  string `json:"order_by" form:"order_by"`                                     // 排序字段
	Order    string `json:"order" form:"order" binding:"omitempty,oneof=asc desc"`        // 排序方向：asc 或 desc
}

// GetPage 获取页码，默认为1
func (p *PaginationRequest) GetPage() int {
	if p.Page <= 0 {
		return 1
	}
	return p.Page
}

// GetPageSize 获取每页数量，默认为10
func (p *PaginationRequest) GetPageSize() int {
	if p.PageSize <= 0 {
		return 10
	}
	if p.PageSize > 100 {
		return 100
	}
	return p.PageSize
}

// GetOffset 计算偏移量
func (p *PaginationRequest) GetOffset() int {
	return (p.GetPage() - 1) * p.GetPageSize()
}

// GetOrderBy 获取排序字段，默认为id
func (p *PaginationRequest) GetOrderBy() string {
	if p.OrderBy == "" {
		return "id"
	}
	return p.OrderBy
}

// GetOrder 获取排序方向，默认为desc
func (p *PaginationRequest) GetOrder() string {
	if p.Order == "" {
		return "desc"
	}
	return p.Order
}

// Paginate GORM分页Scope - 更优雅的分页实现
func (p *PaginationRequest) Paginate() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		page := p.GetPage()
		pageSize := p.GetPageSize()
		offset := (page - 1) * pageSize

		// 添加排序
		orderClause := p.GetOrderBy() + " " + p.GetOrder()

		return db.Offset(offset).Limit(pageSize).Order(orderClause)
	}
}

// PaginateResult 分页结果结构
type PaginateResult struct {
	Data interface{}    `json:"data"`
	Meta PaginationMeta `json:"meta"`
}

// NewPaginateResult 创建分页结果
func NewPaginateResult(data interface{}, pagination *PaginationRequest, total int64) *PaginateResult {
	pageSize := pagination.GetPageSize()
	currentPage := pagination.GetPage()
	totalPages := int(total+int64(pageSize)-1) / pageSize

	return &PaginateResult{
		Data: data,
		Meta: PaginationMeta{
			CurrentPage: currentPage,
			PerPage:     pageSize,
			Total:       total,
			TotalPages:  totalPages,
			HasNext:     currentPage < totalPages,
			HasPrev:     currentPage > 1,
		},
	}
}

// PaginationMetaEnhanced  分页元数据 - 增强版本
type PaginationMetaEnhanced struct {
	CurrentPage int   `json:"current_page"`
	PerPage     int   `json:"per_page"`
	Total       int64 `json:"total"`
	TotalPages  int   `json:"total_pages"`
	HasNext     bool  `json:"has_next"`
	HasPrev     bool  `json:"has_prev"`
	From        int   `json:"from"`
	To          int   `json:"to"`
}
