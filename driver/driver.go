package driver

import (
	"time"
)

type Driver interface {
	// Create(key string, dir bool, value string, expireTime time.Time, unique bool) error
	Delete(key string, dir, recursive bool) error
	Get(key string, recursive bool, sort bool) []interface{}
	Set(key string, dir bool, value string, expireTime time.Time) error
	// Update() error
}
