package main

import (
    "fmt"
    "net/http"
    "io/ioutil"
		"os"
		"time"
    )

func main() {
	start := time.Now()
	for j := 0; j < 100000; j++ {
    response, err := http.Get("http://localhost:8080/")
    if err != nil {
        fmt.Printf("%s", err)
        os.Exit(1)
    } else {
        defer response.Body.Close()
        ioutil.ReadAll(response.Body)
        // fmt.Printf("%s\n", string(contents))
    }
	}
	t := time.Now()
	fmt.Println(t.Sub(start))
}