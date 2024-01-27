package db

import (
	"log"
	"sync"
)

type Cache struct {
	mutex  *sync.RWMutex
	item   map[int64]Order
	dbinst *Postgres
}

type Item struct {
	Value Order
}

func CacheInit(pg *Postgres) *Cache {
	c := Cache{}
	c.NewCache(pg)
	return &c
}

func (c *Cache) NewCache(pg *Postgres) {
	log.Println("Init cache")
	c.dbinst = pg
	pg.SetCache(c)
	c.mutex = &sync.RWMutex{}

	item, err := pg.GetCacheStateInDB()
	if err != nil {
		log.Println("Error to set cache state from database")
	}
	c.item = item
}

func (c *Cache) SetOrderInCache(order Order, orderId int64) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.item[orderId] = order
	log.Println("Order added in cache")

	c.dbinst.SaveCacheInDB(orderId)
}

func (c *Cache) GetOrderInCache(orderId int64) (*OrderDto, error) {
	var orderDto = &OrderDto{}
	c.mutex.RLock()

	item, found := c.item[orderId]

	c.mutex.RUnlock()
	if !found {
		log.Println("Miss cache, start found in DB")
		order, err := c.dbinst.GetOrderById(orderId)

		if err != nil {
			log.Println("Id is not exists")
			return orderDto, err
		}
		log.Println("Order found in DB")
		c.SetOrderInCache(order, orderId)

		orderDto = c.CreateDto(order)
		return orderDto, err
	}
	log.Println("Order found in cache")

	orderDto = c.CreateDto(item)
	return orderDto, nil
}

func (c *Cache) CreateDto(order Order) *OrderDto {
	var orderDto = &OrderDto{}
	orderDto.OrderUID = order.OrderUID
	orderDto.TrackNumber = order.TrackNumber
	orderDto.CustomerId = order.CustomerId
	orderDto.DeliveryService = order.DeliveryService
	orderDto.DateCreated = order.DateCreated

	return orderDto
}
