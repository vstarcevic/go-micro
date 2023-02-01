package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

const maxCount = 5
const expiration = 10

func RateLimiterMiddleware(redis *redis.Client) func(http.Handler) http.Handler {
	f := func(h http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ip := GetRealIP(r)
			if userBlocked(redis, ip) {
				w.Header().Set("Content-Type", "text/plain; charset=utf-8")
				w.Header().Set("X-Content-Type-Options", "nosniff")
				w.WriteHeader(http.StatusTooManyRequests)
				fmt.Fprint(w, "Too many attempts")

			} else {
				h.ServeHTTP(w, r)
				return
			}
		}
		return http.HandlerFunc(fn)
	}
	return f
}

func GetRealIP(r *http.Request) string {
	IPAddress := r.Header.Get("X-Real-IP")
	if IPAddress == "" {
		IPAddress = r.Header.Get("X-Forwarder-For")
	}
	if IPAddress == "" {
		IPAddress = r.RemoteAddr
	}
	return IPAddress
}

func userBlocked(rdb *redis.Client, key string) bool {

	val, err := rdb.Get(ctx, key).Result()
	if err != nil {
		val = "0"
	}

	currentRequests, err := strconv.Atoi(val)
	if err != nil {
		panic(err)
	}

	if currentRequests > maxCount {
		return true
	} else {
		rdb.Incr(ctx, key)
		ttl := rdb.TTL(ctx, key)
		if ttl.Val() == -1 {
			rdb.Expire(ctx, key, expiration*time.Second)
		}
		return false

	}

}
