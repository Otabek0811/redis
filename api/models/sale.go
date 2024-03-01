package models

type SalePrimaryKey struct {
	Id string `json:"id"`
}

type CreateSale struct {
	UserID string `json:"user_id"`
	Total  int    `json:"total"`
	Count  int    `json:"count"`
}

type Sale struct {
	Id        string `json:"id"`
	UserID    string `json:"user_id"`
	Total     int    `json:"total"`
	Count     int    `json:"count"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type UpdateSale struct {
	Id       string `json:"id"`
	Total int `json:"total"`
	Count int `json:"count"`
}

type SaleGetListRequest struct {
	Offset int    `json:"offset"`
	Limit  int    `json:"limit"`
	Search string `json:"search"`
}

type SaleGetListResponse struct {
	Count int     `json:"count"`
	Sales []*Sale `json:"sales"`
}
