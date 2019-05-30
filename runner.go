package main

import (
    "fmt"
    "net/http"
    "io/ioutil"
    // "os"
    "time"
    "bytes"
    // "strconv"
    // "encoding/binary"
    )

func makeHttpPostReq(url string, ClickTimeID string){

    client := http.Client{}

    var jsonprep string = `{"CTID":"`+ClickTimeID+`","CTCompanyID": "1412asdf2001","Flag":true}`

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

}

func main() {
	start := time.Now()
	for j := 0; j < 1000; j++ {
    // response, err := http.Get("http://localhost:8080/" + "15kj3ba412")
        makeHttpPostReq("http://localhost:8080/", "15kj3ba412")
    // if err != nil {
    //     fmt.Printf("%s", err)
    //     os.Exit(1)
    // } else {
    //     defer response.Body.Close()
    //     contents, err := ioutil.ReadAll(response.Body)
    //     if err != nil {
    //         fmt.Printf("%s", err)
    //         os.Exit(1)
    //     }
    //     // i := int64(binary.LittleEndian.Uint64(contents))
    //     fmt.Printf("Completed request %v and result was %v\n", j, string(contents))
    // }
	}
	t := time.Now()
	fmt.Println(t.Sub(start))
}