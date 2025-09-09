package cache

import (
	"container/list"

	models "test-task/internal/models"
)

type Cache struct {
	capacity  int
	cacheMap  map[string]*list.Element
	cacheList *list.List
}

func CreateCache(capacity int) *Cache {
	return &Cache{
		capacity:  capacity,
		cacheMap:  make(map[string]*list.Element),
		cacheList: list.New(),
	}
}

func (cache *Cache) Add(order *models.Order) {

	if existingElement, exist := cache.cacheMap[order.OrderUID]; exist {
		existingElement.Value = order
		cache.cacheList.MoveToFront(existingElement)
		return
	}

	element := cache.cacheList.PushFront(order)
	cache.cacheMap[order.OrderUID] = element

	if len(cache.cacheMap) > cache.capacity {
		cache.removeOldest()
	}
}

func (cache *Cache) Get(order_uid string) (order *models.Order, exist bool) {

	if element, exist := cache.cacheMap[order.OrderUID]; exist {
		cache.cacheList.MoveToFront(element)
		return element.Value.(*models.Order), true
	}
	return nil, false
}

func (cache *Cache) removeOldest() {

	oldestElement := cache.cacheList.Back()
	if oldestElement != nil {
		delete(cache.cacheMap, oldestElement.Value.(*models.Order).OrderUID)
		cache.cacheList.Remove(oldestElement)
	}
}

/* func FillLatestData(cache *Cache) {

	latestOrders, err := models.GetOrders(cache.Capacity)
	if err != nil {
		log.Println(err)
		return
	}

	for _, order := range latestOrders {
		cache.Add(order.Order_uid, order)
	}
} */
