package redis

import (
	"context"
	"golang-test-task/internal/entities"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/mailru/easyjson"
)

// Client is for wrapping original redis.Client
type Client struct {
	client   *redis.Client
	duration time.Duration
}

// NewClientForTest is constructor for testing
func NewClientForTest(client *redis.Client) *Client {
	return &Client{client: client, duration: 5 * time.Minute}
}

// NewClient is constructor for Client
func NewClient(ctx context.Context, config Config) *Client {
	options := redis.Options{
		Addr:     config.Addr,
		Password: config.Password,
		DB:       config.DB,
	}
	client := redis.NewClient(&options)

	if _, err := client.Ping(ctx).Result(); err != nil {
		panic(err)
	}
	duration := 5 * time.Minute
	return &Client{client: client, duration: duration}
}

// FindItemValue gives cached item
func (c *Client) FindItemValue(ctx context.Context, key string) (*entities.APIAdItem, error) {
	v, err := c.client.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	if len(v) == 0 {
		return nil, nil
	}
	var itm entities.APIAdItem
	err = easyjson.Unmarshal([]byte(v), &itm)
	if err != nil {
		return nil, err
	}
	return &itm, nil
}

// SetItemValue setting item to cache
func (c *Client) SetItemValue(ctx context.Context, key string, item entities.APIAdItem) error {
	bs, err := easyjson.Marshal(item)
	if err != nil {
		return err
	}
	_, err = c.client.Set(ctx, key, string(bs), c.duration).Result()
	return err
}

// GetDuration gives duration of key
func (c *Client) GetDuration() time.Duration {
	return c.duration
}
