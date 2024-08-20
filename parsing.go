package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/mmcdole/gofeed"
)

const (
	FETCH_WORKERS_AMOUNT = 10
	FETCH_TIME_LIMIT     = 5 * time.Second

	UNREACHABLE_ERROR = "couldn't reach"
	TIME_LIMIT_ERROR  = "time to fetch exceeded"
	XML_FORMAT_ERROR  = "could not parse XML data"
)

// Self-explanatory. "Name" field is the user alias of the source.
type FeedError struct {
	What string
	Name string
}

func FeedErrorNew(what, name string) FeedError {
	return FeedError{what, name}
}

func (err FeedError) Error() string {
	return fmt.Sprintf("feed error: %s (\"%s\")", err.What, err.Name)
}

// FetchAndParse makes requests to the sourses URLs and returns results of the calls.
// This function uses multi-tread worker pool
func FetchAndParse(sources []Source) []ParseResult {
	sourcesChannel := make(chan Source)
	resultChannel := make(chan ParseResult, len(sources))

	wg := sync.WaitGroup{}
	ctx, cancel := context.WithTimeout(context.Background(), FETCH_TIME_LIMIT)
	defer cancel()

	sender := func() {
		for _, source := range sources {
			sourcesChannel <- source
		}
		close(sourcesChannel)
	}

	worker := func() {
		wg.Add(1)
		defer wg.Done()

		fp := gofeed.NewParser()

		for source := range sourcesChannel {
			select {
			case <-ctx.Done():
				resultChannel <- ParseResult{source.Name, nil, FeedErrorNew(TIME_LIMIT_ERROR, source.Name)}

			default:
				feed, err := fp.ParseURLWithContext(source.Url, ctx)
				// if err != nil {}
				resultChannel <- ParseResult{source.Name, feed, err}
			}
		}
	}

	go sender()
	for i := 0; i < FETCH_WORKERS_AMOUNT; i++ {
		go worker()
	}

	wg.Wait()

	results := make([]ParseResult, 0, len(sources))
	for i := 0; i < len(sources); i++ {
		res := <-resultChannel
		results = append(results, res)
	}

	return results
}
