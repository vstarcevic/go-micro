package middleware

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()
var maxCountPerMinute = 5

func RateLimiterMiddleware(redis *redis.Client) func(http.Handler) http.Handler {
	f := func(h http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			if r.Method == "GET" || r.Method == "POST" {
				ip := GetRealIP(r)

				if !userCanMakeRequest(redis, ip) {
					//w.Header().Add("HeyHeyHey", "HeyHeyHey2") // NOTE THIS LINE
					http.Error(w, "Too many attempts slow down", http.StatusForbidden)
				} else {
					h.ServeHTTP(w, r)
					return
				}
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

func getRedisKey(rdb *redis.Client, key string) (int, error) {

	val, err := rdb.Get(ctx, key).Result()
	if err != nil {
		return 0, err
	}

	return strconv.Atoi(val)

}

func userCanMakeRequest(rdb *redis.Client, key string) bool {

	val, err := rdb.Get(ctx, key).Result()
	if err != nil {
		val = "0"
	}

	currentRequests, err := strconv.Atoi(val)
	if err != nil {
		panic(err)
	}

	if currentRequests > maxCountPerMinute {
		return false
	} else {
		rdb.Incr(ctx, key)
		ttl := rdb.TTL(ctx, key)
		if ttl.Val().Seconds() != -1 {
			rdb.Expire(ctx, key, 10*time.Second)
		}
		return true

	}

}
