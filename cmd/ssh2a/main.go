package ssh2a

import (
	"context"
	"github.com/IUnlimit/ssh2a/internal/cache"
	"github.com/IUnlimit/ssh2a/internal/http"
	"github.com/IUnlimit/ssh2a/internal/ssh"
)

func Run() {
	ctx := context.Background()
	cache.InitCache(ctx)
	go ssh.Init(10220)
	go http.Init(10221)
	select {}
}
