package storage

import (
	"sync"
)

type Cache struct {
	Cache map[string]string
	Mu    sync.Mutex
}

func (c *Cache) WriteCache(token string, email string) error {
	c.Mu.Lock()
	defer c.Mu.Unlock()
	if email == "" || token == "" {
		return ErrEmptyInput
	}

	if _, ok := c.Cache[token]; ok {
		return ErrAlreadyExists
	}
	c.Cache[token] = email
	return nil
}

func (c *Cache) GetValue(token string) (string, error) {
	c.Mu.Lock()
	defer c.Mu.Unlock()
	email, ok := c.Cache[token]
	if !ok {
		return "", ErrNotFound
	}
	return email, nil
}

func NewCache() *Cache {
	var c Cache
	c.Cache = make(map[string]string)
	return &c
}
