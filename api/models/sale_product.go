package models

type SaleProductPrimaryKey struct {
	Id string `json:"id"`
}

type CreateSaleProduct struct {
	ProductID string `json:"product_name"`
	Count     int    `json:"count"`
}

type SaleProduct struct {
	Id                string `json:"id"`
	ProductID         string `json:"product_name"`
	Price             int    `json:"price"`
	Discount          int    `json:"discount"`
	DiscountType      string `json:"discount_type"`
	PriceWithDiscount int    `json:"price_with_discount"`
	DiscountPrice     int    `json:"discount_price"`
	Count             int    `json:"count"`
	Total             int    `json:"total"`
	CreatedAt         string `json:"created_at"`
	UpdatedAt         string `json:"updated_at"`
}

type UpdateSaleProduct struct {
	Id           string `json:"id"`
	ProductID    string `json:"product_name"`
	Price        int    `json:"price"`
	Discount     int    `json:"discount"`
	DiscountType string `json:"discount_type"`
	Count        int    `json:"count"`
}

type SaleProductGetListRequest struct {
	Offset int    `json:"offset"`
	Limit  int    `json:"limit"`
	Search string `json:"search"`
}

type SaleProductGetListResponse struct {
	Count        int            `json:"count"`
	SaleProducts []*SaleProduct `json:"saleProducts"`
}
