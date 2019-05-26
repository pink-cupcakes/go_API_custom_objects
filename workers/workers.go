// In this example we'll look at how to implement
// a _worker pool_ using goroutines and channels.

package workers

import "fmt"
// import "time"

// Here's the worker, of which we'll run several
// concurrent instances. These workers will receive
// work on the `jobs` channel and send the corresponding
// results on `results`. We'll sleep a second per job to
// simulate an expensive task.
func worker(id int, q struct{Jobs chan int; Results chan int}) {
    for j := range q.Jobs {
        fmt.Println("worker", id, "started  job", j)
        // time.Sleep(time.Second)
        fmt.Println("worker", id, "finished job", j)
        q.Results <- j * 2
    }
}

func Start() struct{Jobs chan int; Results chan int} {
    // In order to use our pool of workers we need to send
    // them work and collect their results. We make 2
    // channels for this.
    type Queues struct {
        Jobs    chan int
        Results chan int
    }
    q := Queues{make(chan int, 100000), make(chan int, 100000)}
    // jobs := make(chan int, 100000)
    // results := make(chan int, 100000)

    // This starts up 3 workers, initially blocked
    // because there are no jobs yet.
    for w := 1; w <= 10000; w++ {
        go worker(w, q)
    }
		
    // Here we send 5 `jobs` and then `close` that
    // channel to indicate that's all the work we have.
    // for j := 1; j <= 5; j++ {
    // 	jobs <- j
    // }

    return q

    // close(jobs)

    // Finally we collect all the results of the work.
    // for a := 1; a <= 5; a++ {
    //     <-results
    // }
}