package main

import (
	"github.com/suman7383/lru-cache/cache"
	"github.com/suman7383/lru-cache/server"
)

func main() {
	cache := cache.NewCache(3)
	srv := server.New(":3000", cache)

	srv.Start()
}
