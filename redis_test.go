package chassis

import (
	"context"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"

	"c6x.io/chassis/config"
)

func TestReadOptions(t *testing.T) {
	config.LoadFromEnvFile()

	var opts redis.Options
	readRedisOptions(&opts)
	assert.NotEmpty(t, opts)
	assert.NotNil(t, &opts)
	assert.Equal(t, "chassis-redis-ut:6379", opts.Addr)
}
func TestRedisClient(t *testing.T) {
	config.LoadFromEnvFile()
	ctx := context.Background()
	pong := RedisClient().Ping(ctx)
	t.Logf("%v", pong)
	RedisClient().Set(ctx, "test", "123", 5*time.Minute)

	val := RedisClient().Get(ctx, "test").Val()
	assert.NotEmpty(t, val)
	assert.Equal(t, "123", val)
}
