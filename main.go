package main

import (
	"flag"
	"fmt"
	"log"
	"github.com/valyala/fasthttp"
	"learn/cache"
)

var (
	segs = flag.Int("segs", 255, "Total buckets number in Cache")
	addr = flag.String("addr", ":8080", "TCP address to listen to")
	compress = flag.Bool("compress", false, "Whether to enable transparent response compression")
	kv = cache.New(*segs)
)

func main() {
	flag.Parse()
	h := requestHandler
	if *compress {
		h = fasthttp.CompressHandler(h)
	}

	if err := fasthttp.ListenAndServe(*addr, h); err != nil {
		log.Fatalf("Error in ListenAndServe: %s", err)
	}
}

func requestHandler(ctx *fasthttp.RequestCtx) {
	args := ctx.QueryArgs()
	cmd := string(args.Peek("cmd"))
	key := string(args.Peek("key"))
	if cmd == "set" {
		val := args.GetUintOrZero("val")
		fmt.Fprintf(ctx, "set %s=%d!\n\n", key, val)
		kv.Set(key, val)
	} else {
		fmt.Fprintf(ctx, "%s is %d\n\n", key, kv.Get(key))
	}
	ctx.SetContentType("text/plain; charset=utf8")
}