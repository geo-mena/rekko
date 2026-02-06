package response

import "net/http"

func Success(w http.ResponseWriter, statusCode int, message string, data any) {
	write(w, statusCode, Response{
		Status:  true,
		Message: message,
		Data:    data,
	})
}

func SuccessWithPagination(w http.ResponseWriter, statusCode int, message string, data any, params PaginationParams) {
	totalPages := 0
	if params.PerPage > 0 {
		totalPages = int((params.Total + int64(params.PerPage) - 1) / int64(params.PerPage))
	}

	write(w, statusCode, Response{
		Status:  true,
		Message: message,
		Data:    data,
		Meta: &Meta{
			Pagination: &Pagination{
				CurrentPage: params.Page,
				PerPage:     params.PerPage,
				TotalItems:  params.Total,
				TotalPages:  totalPages,
				HasNext:     params.Page < totalPages,
			},
		},
	})
}

func MessageOnly(w http.ResponseWriter, statusCode int, message string) {
	write(w, statusCode, Response{
		Status:  true,
		Message: message,
	})
}

func NoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}
