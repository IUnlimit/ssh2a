package cache

import (
	"errors"
	"github.com/bluele/gcache"
	log "github.com/sirupsen/logrus"
	"time"
)

type IPCountCache struct {
	Tag string
	// string: int
	holder gcache.Cache
	// ip: timestamp(time.Time)
	Passed gcache.Cache
}

func NewIPCountCache(tag string, size int, expiration time.Duration) *IPCountCache {
	return &IPCountCache{
		Tag:    tag,
		holder: gcache.New(size).LRU().Expiration(expiration).Build(),
		Passed: gcache.New(64).LRU().Expiration(expiration).Build(),
	}
}

func (c *IPCountCache) UpdateIPStatus(ip string, status bool) error {
	count, err := c.holder.Get(ip)
	if err != nil && !errors.Is(err, gcache.KeyNotFoundError) {
		return err
	}
	if status { // access
		if count != nil {
			c.holder.Remove(ip)
			log.Warnf("[%s] IP %s was accessed after %d failures", c.Tag, ip, count)
		}
		err := c.Passed.Set(ip, time.Now())
		if err != nil {
			return err
		}
		return nil
	}
	if count == nil {
		count = 0
	}
	err = c.holder.Set(ip, 1+count.(int))
	if err != nil {
		return err
	}
	return nil
}

func (c *IPCountCache) CheckPassed(ip string) bool {
	passedTime, err := c.Passed.Get(ip)
	if errors.Is(err, gcache.KeyNotFoundError) {
		return false
	}
	log.Infof("Connect via SSH from a verified IP(%s) at %s", ip, passedTime)
	return true
}
