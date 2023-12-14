package crack

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"
)

// Background返回一个非空的Context。 它永远不会被取消，没有值，也没有期限。
// 它通常在main函数，初始化和测试时使用，并用作传入请求的顶级上下文。
var ctx = context.Background()

func rediscon(cancel context.CancelFunc, host, user, passwd string, port, timeout int) {
	client := redis.NewClient(&redis.Options{
		Addr:            fmt.Sprintf("%s:%d", host, port),
		Username:        user,
		Password:        passwd,
		DB:              0,
		DialTimeout:     time.Duration(timeout) * time.Second,
		MinRetryBackoff: time.Duration(timeout) * time.Second,
		ReadTimeout:     time.Duration(timeout) * time.Second,
	})
	_, err := client.Ping(ctx).Result()
	if err == nil {
		end(host, user, passwd, port, "Redis")
		cancel()
	}
}
