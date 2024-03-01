package redis

import (
	"app/config"
	"app/storage"

	"github.com/go-redis/redis"
)

type Cache struct {
	cachedb *redis.Client
	product *productRepo
}

func NewConnectionRedis(cfg *config.Config) (storage.CacheI, error) {
	var client = redis.NewClient(
		&redis.Options{
			Addr:     cfg.RedisHost + cfg.RedisPort,
			Password: cfg.RedisPassword,
			DB:       cfg.RedisDB,
		},
	)

	_, err := client.Ping().Result()
	if err != nil {
		return nil, err
	}

	return &Cache{
		cachedb: client,
	}, nil
}

func (c *Cache) Close() {
	c.cachedb.Close()
}

func (c *Cache) Product() storage.ProductRepoCacheI {

	if c.product == nil {
		c.product = NewProductRepo(c.cachedb)
	}

	return c.product
}
