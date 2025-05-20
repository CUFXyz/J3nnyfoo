package storage

import (
	"sync"
)

type Cache struct {
	Cache map[string]string
	Mu    sync.Mutex
}

func (c *Cache) WriteCache(token, email string) error {

	if email == "" || token == "" {
		return ErrEmptyInput
	}
	c.Mu.Lock()
	defer c.Mu.Unlock()
	if _, ok := c.Cache[token]; ok {
		return ErrAlreadyExists
	}
	c.Cache[token] = email
	return nil
}

func (c *Cache) GetValue(val string) (string, error) {
	c.Mu.Lock()
	defer c.Mu.Unlock()
	resource, ok := c.Cache[val]
	if !ok {
		return "", ErrNotFound
	}
	return resource, nil
}

func NewCache() *Cache {
	var c Cache
	c.Cache = make(map[string]string)
	return &c
}
