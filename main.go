package main
import (
  "net/http"
	"strings"
	"github.com/sandbox/workers"
)

type Queues struct {
	jobs chan int
	increment int
}

func (queues *Queues) sayHello(w http.ResponseWriter, r *http.Request) {
  message := r.URL.Path
  message = strings.TrimPrefix(message, "/")
  switch r.Method {
	case http.MethodGet:
		queues.jobs <- queues.increment
		queues.increment++
    // message = "Hello " + message
    // w.Write([]byte(message)) 
  case http.MethodPost:
    w.Write([]byte("Testing post"))
  default:
    http.Error(w, "Invalid method", 405)
  }
}

// func sayHello(w http.ResponseWriter, r *http.Request) {
//   message := r.URL.Path
//   message = strings.TrimPrefix(message, "/")
//   switch r.Method {
// 	case http.MethodGet:
// 		jobs <- 29
//     // message = "Hello " + message
//     // w.Write([]byte(message)) 
//   case http.MethodPost:
//     w.Write([]byte("Testing post"))
//   default:
//     http.Error(w, "Invalid method", 405)
//   }
// }

func main() {
	q := workers.Start()
	jobs := &Queues{jobs: q.Jobs, increment: 0}
	results := q.Results
	http.HandleFunc("/", jobs.sayHello)
  if err := http.ListenAndServe(":8080", nil); err != nil {
    panic(err)
  }
}