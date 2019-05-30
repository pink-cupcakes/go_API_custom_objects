package main

import (
	"fmt"
	"net/http"
	"sort"
	"time"
	"strconv"
	"bytes"
	"io/ioutil"
	"math/rand"
)

// a struct to hold the result from each request including an index
// which will be used for sorting the results after they come in
type result struct {
	index int
	res   http.Response
	err   error
}


// boundedParallelGet sends requests in parallel but only up to a certain
// limit, and furthermore it's only parallel up to the amount of CPUs but
// is always concurrent up to the concurrency limit
func boundedParallelGet(urls []string, concurrencyLimit int) []result {

	// this buffered channel will block at the concurrency limit
	semaphoreChan := make(chan struct{}, concurrencyLimit)

	// this channel will not block and collect the http request results
	resultsChan := make(chan *result)

	// make sure we close these channels when we're done with them
	defer func() {
		close(semaphoreChan)
		close(resultsChan)
	}()

	// keen an index and loop through every url we will send a request to
	for i, url := range urls {

		// start a go routine with the index and url in a closure
		go func(i int, url string) {

			// this sends an empty struct into the semaphoreChan which
			// is basically saying add one to the limit, but when the
			// limit has been reached block until there is room
			semaphoreChan <- struct{}{}

			// send the request and put the response in a result struct
			// along with the index so we can sort them later along with
			// any error that might have occoured
			client := http.Client{}


			var jsonprep string = `{"CTID":"` + strconv.Itoa(rand.Int()) + `1","CTCompanyID": "1412asdf2001","Flag":true}`
	
			var jsonStr = []byte(jsonprep)
	
			req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
			req.Header.Set("Content-Type", "application/json")
	
			resp, err := client.Do(req)
			if err != nil {
				fmt.Println("Unable to reach the server.")
			} else {
				body, _ := ioutil.ReadAll(resp.Body)
				fmt.Println("body=", string(body))
			}
			result := &result{i, *resp, err}
			if err != nil {
				fmt.Println("COWABUNGA")
			}
			fmt.Println(i)

			// now we can send the result struct through the resultsChan
			resultsChan <- result

			// once we're done it's we read from the semaphoreChan which
			// has the effect of removing one from the limit and allowing
			// another goroutine to start
			<-semaphoreChan
			return
		}(i, url)
	}

	// make a slice to hold the results we're expecting
	var results []result

	// start listening for any results over the resultsChan
	// once we get a result append it to the result slice
	for {
		result := <-resultsChan
		results = append(results, *result)

		// if we've reached the expected amount of urls then stop
		if len(results) == len(urls) {
			break
		}
	}

	// let's sort these results real quick
	sort.Slice(results, func(i, j int) bool {
		return results[i].index < results[j].index
	})

	// now we're done we return the results
	return results
}

// we'll use the init function to set up the benchmark
// by making a slice of 100 URLs to send requets to
var urls []string

func init() {
	for i := 0; i < 1000; i++ {
		urls = append(urls, "http://localhost:8080/" + strconv.Itoa(i))
	}
}

// the main function sets up an anonymous benchmark func
// that will time how long it takes to get all the URLs
// at the specified concurrency level
//
// and you should see something like the following printed
// depending on how fast your computer and internet is
//
// 5 bounded parallel requests: 100/100 in 5.533223255
// 10 bounded parallel requests: 100/100 in 2.5115351219
// 25 bounded parallel requests: 100/100 in 1.189462884
// 50 bounded parallel requests: 100/100 in 1.17430002
// 75 bounded parallel requests: 100/100 in 1.001383863
// 100 bounded parallel requests: 100/100 in 1.3769354
func main() {
	benchmark := func(urls []string, concurrency int) string {
		startTime := time.Now()
		results := boundedParallelGet(urls, concurrency)
		seconds := time.Since(startTime).Seconds()
		tmplate := "%d bounded parallel requests: %d/%d in %v"
		return fmt.Sprintf(tmplate, concurrency, len(results), len(urls), seconds)
	}

	// fmt.Println(benchmark(urls, 10))
	// fmt.Println(benchmark(urls, 25))
	// fmt.Println(benchmark(urls, 50))
	fmt.Println(benchmark(urls, 1000))
}