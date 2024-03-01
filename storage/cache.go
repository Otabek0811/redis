package storage

import "app/api/models"

type CacheI interface {
	Close()
	Product() ProductRepoCacheI
}

type ProductRepoCacheI interface {
	CreateList(*models.ProductGetListResponse) error
	GetList() (*models.ProductGetListResponse, error)
	Delete() error
	Exists() (bool, error)
}
