package response

type Response struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
	Meta    *Meta  `json:"meta,omitempty"`
}

type Meta struct {
	Pagination *Pagination `json:"pagination,omitempty"`
}

type Pagination struct {
	CurrentPage int   `json:"current_page"`
	PerPage     int   `json:"per_page"`
	TotalItems  int64 `json:"total_items"`
	TotalPages  int   `json:"total_pages"`
	HasNext     bool  `json:"has_next"`
}

type PaginationParams struct {
	Page    int
	PerPage int
	Total   int64
}

type ErrorData struct {
	Code    string        `json:"code"`
	Details []ErrorDetail `json:"details,omitempty"`
}

type ErrorDetail struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}
