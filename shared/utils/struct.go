package utils

type PageAttr struct {
	Page     int `json:"page"`
	Size     int `json:"size"`
	Total    int `json:"total"`
	LastPage int `json:"last_page"`
}
