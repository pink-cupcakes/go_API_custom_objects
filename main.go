package main
import (
  "net/http"
	"strings"
	// "fmt"
	"encoding/binary"
	"strconv"
	"github.com/sandbox/workers"
)

type Queues struct {
	jobs chan int
	results chan int
	increment int
}

func (queues *Queues) sayHello(w http.ResponseWriter, r *http.Request) {
  message := r.URL.Path
	message = strings.TrimPrefix(message, "/")
	if i, err := strconv.Atoi(message); err == nil {
		// fmt.Println(i)
		switch r.Method {
		case http.MethodGet:
			queues.jobs <- i
			queues.increment++
			select {
			case result := <-queues.results:
				b := make([]byte, 8)
				binary.LittleEndian.PutUint64(b, uint64(result))
				// i := int64(binary.LittleEndian.Uint64(b))
				// fmt.Println(i)
				w.Write([]byte(b))
			}
			
		case http.MethodPost:
			w.Write([]byte("Testing post"))
		default:
			http.Error(w, "Invalid method", 405)
		}
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
	jobs := &Queues{jobs: q.Jobs, results: q.Results, increment: 0}
	http.HandleFunc("/", jobs.sayHello)
  if err := http.ListenAndServe(":8080", nil); err != nil {
    panic(err)
  }
}