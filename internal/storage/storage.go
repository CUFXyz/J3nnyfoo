package storage

import (
	"sync"
)

type Cache struct {
	Cache map[string]string
	Mu    sync.Mutex
}

func (c *Cache) WriteCache(email, token string) error {
	c.Mu.Lock()
	defer c.Mu.Unlock()
	if email == "" || token == "" {
		return ErrEmptyInput
	}

	if _, ok := c.Cache[email]; ok {
		return ErrAlreadyExists
	}
	c.Cache[email] = token
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

func (c *Cache) LookUp(value string) (bool, error) {
	for _, link := range c.Cache {
		if link == value {
			return true, ErrAlreadyExists
		}
	}
	return false, nil
}
