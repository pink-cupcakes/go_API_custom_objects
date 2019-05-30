package main
import (
  "net/http"
	// "strings"
	// "fmt"
	// "encoding/binary"
	// "strconv"
	// "log"
	// "os"
	// "io/ioutil"
	// "encoding/json"
	"github.com/sandbox/workers"
)

type Queues struct {
	jobs 				chan *http.Request
	results 		chan []byte
	increment 	int
}

func (queues *Queues) sayHello(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			queues.jobs <- r
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