package main

import (
	"context"
	"github.com/nettyrnp/url-shortener/util"
)

func main() {
	ctx := context.Background()
	app, err := NewApp(ctx)
	util.Die(err)
	util.Die(app.Run(ctx))
}
