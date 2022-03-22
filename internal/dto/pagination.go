package dto

type Pagination struct {
	Filter string
	Page   int32
	Limit  int32
	Offset int32
}
