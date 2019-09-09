package cache

import (
	"errors"
	"time"
)

var (
	ErrCacheMiss = errors.New("cache miss")
)

type Client interface {
	Delete(k string) error
	Flush() error
	Get(k string, v interface{}) error
	Set(k string, v interface{}, expire time.Duration) error
}
