package model

type Pagination struct {
	Page           int   `json:"page"`
	PageSize       int   `json:"pageSize"`
	TotalCount     int   `json:"totalCount"`
	TotalPages     int   `json:"totalPages"`
	Pages          []int `json:"pages"`
	StartItem      int   `json:"startItem"`
	EndItem        int   `json:"endItem"`
	PerformerIdStr int   `json:"performerIdStr"`
}
