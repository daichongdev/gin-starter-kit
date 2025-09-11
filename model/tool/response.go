package tool

// APIResponse 统一API响应结构
type APIResponse struct {
	Status  string      `json:"status"`  // "success" 或 "error"
	Message string      `json:"message"` // 响应消息
	Data    interface{} `json:"data"`    // 响应数据，成功时包含具体数据，失败时为null
}

// PaginationResponse 分页响应结构
type PaginationResponse struct {
	Status  string         `json:"status"`
	Message string         `json:"message"`
	Data    interface{}    `json:"data"`
	Meta    PaginationMeta `json:"meta"`
}

type PaginationMeta struct {
	CurrentPage int   `json:"current_page"`
	PerPage     int   `json:"per_page"`
	Total       int64 `json:"total"`
	TotalPages  int   `json:"total_pages"`
	HasNext     bool  `json:"has_next"`
	HasPrev     bool  `json:"has_prev"`
}

// SuccessResponse 响应构造函数
func SuccessResponse(message string, data interface{}) APIResponse {
	return APIResponse{
		Status:  "success",
		Message: message,
		Data:    data,
	}
}

func ErrorResponse(message string) APIResponse {
	return APIResponse{
		Status:  "error",
		Message: message,
		Data:    nil,
	}
}

func PaginationSuccessResponse(message string, data interface{}, meta PaginationMeta) PaginationResponse {
	return PaginationResponse{
		Status:  "success",
		Message: message,
		Data:    data,
		Meta:    meta,
	}
}
