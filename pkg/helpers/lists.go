package helpers

func NewListConfig(sortField string, isAsceding bool, pageSize int, page int, andLogic bool, expressionFilter []ExpressionFilter) *ListConfig {
	config := &ListConfig{
		SortField:         sortField,
		IsAscending:       isAsceding,
		PageSize:          pageSize,
		Page:              page,
		AndLogic:          andLogic,
		ExpressionFilters: expressionFilter,
	}
	return config
}

type ListConfig struct {
	SortField         string             // The column name of struct that we want to sort by them
	IsAscending       bool               // Ascending sort or Descending  ?
	PageSize          int                // The count of data that we need
	Page              int                // The count of data that we want to skip
	AndLogic          bool               // And(&) between ExpressionFilters ?
	ExpressionFilters []ExpressionFilter // Filters - If we need filters
}

type ExpressionFilter struct {
	PropertyName string             // The column name of struct that we want to filter on this
	Value        interface{}        // Value of that column name
	Comparison   ComparisonOperator // Comparison
}

type ComparisonOperator uint8

const (
	Equal ComparisonOperator = iota
	LessThan
	LessThanOrEqual
	GreaterThan
	GreaterThanOrEqual
	NotEqual
	Contains
	StartsWith
	EndsWith
)

type Pages[s interface{}] struct {
	Page     int
	PageSize int
	Content  []s
	Size     int
}
