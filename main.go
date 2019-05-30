package main
import (
  "net/http"
	"strings"
	"fmt"
	// "encoding/binary"
	// "strconv"
	// "log"
	// "os"
	"io/ioutil"
	// "encoding/json"
	"github.com/sandbox/workers"
)

type Queues struct {
	jobs 				chan []byte
	results 		chan []byte
	increment 	int
}

func (queues *Queues) sayHello(w http.ResponseWriter, r *http.Request) {
  message := r.URL.Path
	message = strings.TrimPrefix(message, "/")
		switch r.Method {
		case http.MethodPost:
			reqBody, err := ioutil.ReadAll(r.Body)
			if err != nil {
				fmt.Println(err)
			}
			queues.jobs <- reqBody
			queues.increment++
			select {
			case result := <-queues.results:
				w.Write([]byte(result))
			}
			
		case http.MethodGet:
			reqBody, err := ioutil.ReadAll(r.Body)
			if err != nil {
				fmt.Println(err)
			}
			queues.jobs <- reqBody
			queues.increment++
			select {
			case result := <-queues.results:
				w.Write([]byte(result))
			}
		default:
			http.Error(w, "Invalid method", 405)
		}
}

func main() {
	q := workers.Start()
	jobs := &Queues{jobs: q.Jobs, results: q.Results, increment: 0}
	http.HandleFunc("/", jobs.sayHello)
  if err := http.ListenAndServe(":8080", nil); err != nil {
    panic(err)
  }
}