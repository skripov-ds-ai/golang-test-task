package memcached

import (
	"golang-test-task/internal/entities"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/mailru/easyjson"
)

// Client is a wrapper for Memcached
type Client struct {
	client   *memcache.Client
	duration time.Duration
}

// NewClient is a constructor for Memcached wrapper
func NewClient(memcacheStrings ...string) *Client {
	m := memcache.New(memcacheStrings...)
	if err := m.Ping(); err != nil {
		panic(err)
	}
	m.Timeout = 100 * time.Millisecond
	m.MaxIdleConns = 100
	duration := 5 * time.Minute
	return &Client{client: m, duration: duration}
}

// FindItemValue gives cached item
func (c *Client) FindItemValue(key string) (*entities.APIAdItem, error) {
	v, err := c.client.Get(key)
	if err != nil {
		return nil, err
	}
	if v == nil {
		return nil, nil
	}
	var itm entities.APIAdItem
	err = easyjson.Unmarshal(v.Value, &itm)
	if err != nil {
		return nil, err
	}
	return &itm, nil
}

// SetItemValue setting item to cache
func (c *Client) SetItemValue(key string, item entities.APIAdItem) error {
	bs, err := easyjson.Marshal(item)
	if err != nil {
		return err
	}
	cacheItem := memcache.Item{Value: bs, Key: key, Expiration: int32(c.duration)}
	err = c.client.Set(&cacheItem)
	return err
}

// GetDuration gives duration of key
func (c *Client) GetDuration() time.Duration {
	return c.duration
}
