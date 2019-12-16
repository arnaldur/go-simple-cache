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
	Nearby(k string, lon, lat, radius float64) ([]Location, error)
	GeoAdd(k string, locations ...Location) error
}

type Location struct {
	Name      string
	Longitude float64
	Latitude  float64
	Distance  float64
}
