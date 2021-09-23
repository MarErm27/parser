package main

import (
	"flag"
	"parser/bechmarker"
	"parser/http/rest"
)

func main() {
	responseTime := flag.Int("responseTime", 3, "max site response time in seconds")
	flag.Parse()
	bechmarker.SetResponseTime(*responseTime)
	rest.Handle()
}
