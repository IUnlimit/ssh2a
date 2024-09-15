package cache

import (
	"github.com/bluele/gcache"
	log "github.com/sirupsen/logrus"
	"time"
)

// string: int
var ipCache gcache.Cache

func init() {
	ipCache = gcache.New(512).Simple().Expiration(time.Hour).Build()
}

func UpdateIPStatus(ip string, status bool) error {
	count, err := ipCache.Get(ip)
	if err != nil {
		return err
	}
	if status { // access
		if count != nil {
			ipCache.Remove(ip)
			log.Warnf("IP %s was accessed after %d failures", ip, count)
		}
		return nil
	}
	err = ipCache.Set(ip, 1+count.(int))
	if err != nil {
		return err
	}
	return nil
}
