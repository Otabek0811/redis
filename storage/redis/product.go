package redis

import (
	"app/api/models"
	"encoding/json"

	"github.com/go-redis/redis"
)

const PRODUCT_SLUG = "products"

type productRepo struct {
	cachedb *redis.Client
}

func NewProductRepo(redisDB *redis.Client) *productRepo {
	return &productRepo{
		cachedb: redisDB,
	}
}

func (c *productRepo) CreateList(req *models.ProductGetListResponse) error {

	body, err := json.Marshal(req)
	if err != nil {
		return err
	}

	err = c.cachedb.Set(PRODUCT_SLUG, body, 0).Err()
	if err != nil {
		return err
	}

	return nil
}

func (c *productRepo) GetList() (*models.ProductGetListResponse, error) {

	var resp *models.ProductGetListResponse

	data, err := c.cachedb.Get(PRODUCT_SLUG).Result()
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(data), &resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *productRepo) Delete() error {

	err := c.cachedb.Del(PRODUCT_SLUG).Err()
	if err != nil {
		return err
	}

	return nil
}

func (c *productRepo) Exists() (bool, error) {

	rowsAffected, err := c.cachedb.Exists(PRODUCT_SLUG).Result()
	if err != nil {
		return false, err
	}

	if rowsAffected <= 0 {
		return false, nil
	}

	return true, nil
}
