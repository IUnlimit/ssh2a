package cache

import (
	"context"
	"github.com/allegro/bigcache/v3"
	"log"
	"strconv"
	"time"
)

// ip: request count
var tempAddrCache *bigcache.BigCache

// ip: access timestamp
var accessIpMap map[string]int64

func InitCache(ctx context.Context) {
	var err error
	tempAddrCache, err = bigcache.New(ctx, bigcache.DefaultConfig(3*time.Minute))
	if err != nil {
		log.Fatalf("Error occurred when initial cache: %v", err)
	}
	accessIpMap = make(map[string]int64)
}

// AddTempCache increases the count of IP attempts and
// returns the current number of attempts
func AddTempCache(ip string) int {
	bytes, err := tempAddrCache.Get(ip)
	if err != nil {
		_ = tempAddrCache.Set(ip, []byte(strconv.Itoa(1)))
		return 1
	}
	count, _ := strconv.Atoi(string(bytes))
	count++
	_ = tempAddrCache.Set(ip, []byte(strconv.Itoa(count)))
	return count
}

func ContainsTempCache(ip string) bool {
	bytes, err := tempAddrCache.Get(ip)
	return err == nil && bytes != nil
}

func InvalidTempCache(ip string) {
	_ = tempAddrCache.Delete(ip)
}

func SetAccessCache(ip string) {
	accessIpMap[ip] = time.Now().UnixMilli()
}

func ContainsAccessMap(ip string) bool {
	_, exist := accessIpMap[ip]
	return exist
}
