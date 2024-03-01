package models

type MarketPrimaryKey struct {
	Id string `json:"id"`
}

type CreateMarket struct {
	Name        string   `json:"name"`
	Address     string   `json:"address"`
	PhoneNumber string   `json:"phone_number"`
	ProductIds  []string `json:"product_ids"`
}

type Market struct {
	Id          string     `json:"id"`
	Name        string     `json:"name"`
	Address     string     `json:"address"`
	PhoneNumber string     `json:"phone_number"`
	CreatedAt   string     `json:"created_at"`
	UpdatedAt   string     `json:"updated_at"`
	Products    []*Product `json:"products"`
}

type UpdateMarket struct {
	Id          string   `json:"id"`
	Name        string   `json:"name"`
	Address     string   `json:"address"`
	PhoneNumber string   `json:"phone_number"`
	ProductIds  []string `json:"product_ids"`
}

type MarketGetListRequest struct {
	Offset int    `json:"offset"`
	Limit  int    `json:"limit"`
	Search string `json:"search"`
}

type MarketGetListResponse struct {
	Count   int       `json:"count"`
	Markets []*Market `json:"markets"`
}
