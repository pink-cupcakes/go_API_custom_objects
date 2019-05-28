package main
import (
  "net/http"
	"strings"
	"fmt"
	// "encoding/binary"
	// "strconv"
	// "log"
	"io/ioutil"
	"encoding/json"
	"github.com/sandbox/workers"
)

type Queues struct {
	jobs chan string
	results chan string
	increment int
}

func (queues *Queues) sayHello(w http.ResponseWriter, r *http.Request) {
  message := r.URL.Path
	message = strings.TrimPrefix(message, "/")
	// fmt.Println(message)
	// if i, err := strconv.Atoi(message); err == nil {
		// fmt.Println(i)
		switch r.Method {
		case http.MethodPost:
			var request map[string]interface{}
			reqBody, err := ioutil.ReadAll(r.Body)
			if err != nil {
				fmt.Println(err)
			}
			json.Unmarshal([]byte(reqBody), &request)
			for key, value := range request {
				fmt.Println(key, value.(string))
			}
			queues.jobs <- message
			queues.increment++
			select {
			case result := <-queues.results:
				// b := make([]byte, 8)
				// binary.LittleEndian.PutUint64(b, uint64(result))
				// i := int64(binary.LittleEndian.Uint64(b))
				// fmt.Println(i)
				w.Write([]byte(result))
			}
			
		case http.MethodGet:
			queues.jobs <- message
			queues.increment++
			select {
			case result := <-queues.results:
				// b := make([]byte, 8)
				// binary.LittleEndian.PutUint64(b, uint64(result))
				// i := int64(binary.LittleEndian.Uint64(b))
				// fmt.Println(i)
				w.Write([]byte(result))
			}
		default:
			http.Error(w, "Invalid method", 405)
		}
	// }
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