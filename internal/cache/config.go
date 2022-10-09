package cache

import (
	"os"
	"strconv"
)

// RedisConfig is struct for storing data about path to Redis Storage
type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

// Load is useful for loading RedisConfig data
func (r *RedisConfig) Load() {
	r.Addr = os.Getenv("REDIS_ADDR")
	r.Password = os.Getenv("REDIS_PASSWORD")
	DB, err := strconv.ParseInt(os.Getenv("REDIS_DB"), 10, 32)
	if err != nil {
		panic(err)
	}
	r.DB = int(DB)
}
