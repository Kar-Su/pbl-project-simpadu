package utils

const DefaultPageSize = 10

type PaginationQuery struct {
	Page int `form:"page,default=1" example:"1"`
}

type PaginationMeta struct {
	Page       int   `json:"page" example:"1"`
	PerPage    int   `json:"per_page" example:"10"`
	TotalItems int64 `json:"total_items" example:"100"`
	TotalPages int   `json:"total_pages" example:"10"`
}

type PaginatedData[T any] struct {
	Items      T              `json:"items"`
	Pagination PaginationMeta `json:"pagination"`
}

func GetOffset(page int) int {
	if page <= 0 {
		page = 1
	}
	return (page - 1) * DefaultPageSize
}

func GetPage(page int) int {
	if page <= 0 {
		return 1
	}
	return page
}

func CountTotalPages(totalItems int64) int {
	if totalItems == 0 {
		return 1
	}
	pages := int(totalItems) / DefaultPageSize
	if int(totalItems)%DefaultPageSize != 0 {
		pages++
	}
	return pages
}

func BuildPaginatedResponse[T any](items T, page int, totalItems int64) PaginatedData[T] {
	return PaginatedData[T]{
		Items: items,
		Pagination: PaginationMeta{
			Page:       GetPage(page),
			PerPage:    DefaultPageSize,
			TotalItems: totalItems,
			TotalPages: CountTotalPages(totalItems),
		},
	}
}
