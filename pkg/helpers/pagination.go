package helpers

type Pagination struct {
	PageSize  int    `json:"page_size"`
	Page      int    `json:"page"`
	SortField string `json:"sort_field"`
	Filter    string `json:"filter"`
}
