package bechmarker

import (
	"context"
	"github.com/go-errors/errors"
	"net/http"
	"parser/core"
	"parser/serp"
	"sync"
	"sync/atomic"
	"time"
)

var maxResponseTime time.Duration

func SetResponseTime(responseTime int) {
	maxResponseTime = time.Duration(responseTime) * time.Second
}

func MakeRequest(ctx context.Context, url string, count *int64) {
	ch := make(chan struct{}, 1)
	go func() {
		resp, err := http.Get(url)
		if err == nil && resp.StatusCode == http.StatusOK {
			ch <- struct{}{}
		}
	}()

	select {
	case <-ctx.Done():
		return
	case <-ch:
		atomic.AddInt64(count, int64(1))
		return
	}
}

func loadURL(item serp.ResponseItem, n int) {
	ctx1 := context.Background()
	ctxt, cancel := context.WithTimeout(ctx1, maxResponseTime)
	defer cancel()

	var count int64

	var wg sync.WaitGroup
	wg.Add(n)

	for i := 0; i < n; i++ {
		go func(ctxt context.Context, i int) {
			defer wg.Done()
			MakeRequest(ctxt, item.Url, &count)
		}(ctxt, i)
	}
	wg.Wait()

	if int(count) == n {
		loadURL(item, n*2+10)
	} else {
		core.Put(item.Host, int(count))
	}
}

func Benchmark(listOfURL []serp.ResponseItem) {
	var wg sync.WaitGroup
	wg.Add(len(listOfURL))
	for i := 0; i < len(listOfURL); i++ {
		go func(i int) {
			defer wg.Done()
			_, err := core.Get(listOfURL[i].Host)
			if errors.Is(err, core.ErrorNoSuchKey) {
				loadURL(listOfURL[i], 1)
			}
		}(i)
	}
	wg.Wait()
}
