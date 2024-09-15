package cache

import (
	"github.com/bluele/gcache"
	"time"
)

// string: int
var loginCache gcache.Cache

func init() {
	loginCache = gcache.New(64).Simple().Expiration(24 * time.Hour).Build()
}
