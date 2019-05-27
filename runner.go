package main

import (
    "fmt"
    "net/http"
    "io/ioutil"
    "os"
    "time"
    // "strconv"
    // "encoding/binary"
    )

func main() {
	start := time.Now()
	for j := 0; j < 10; j++ {
    response, err := http.Get("http://localhost:8080/" + "15kj3ba412")
    if err != nil {
        fmt.Printf("%s", err)
        os.Exit(1)
    } else {
        defer response.Body.Close()
        contents, err := ioutil.ReadAll(response.Body)
        if err != nil {
            fmt.Printf("%s", err)
            os.Exit(1)
        }
        // i := int64(binary.LittleEndian.Uint64(contents))
        fmt.Printf("Completed request %v and result was %v\n", j, string(contents))
    }
	}
	t := time.Now()
	fmt.Println(t.Sub(start))
}