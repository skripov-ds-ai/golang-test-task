package redis

import (
	"os"
	"strconv"
)

// Config is struct for storing data about path to Redis Storage
type Config struct {
	Addr     string
	Password string
	DB       int
}

// Load is useful for loading Config data
func (r *Config) Load() {
	r.Addr = os.Getenv("REDIS_ADDR")
	r.Password = os.Getenv("REDIS_PASSWORD")
	DB, err := strconv.ParseInt(os.Getenv("REDIS_DB"), 10, 32)
	if err != nil {
		panic(err)
	}
	r.DB = int(DB)
}
