package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// A Ttype represents a meaninglessness of our life
type Ttype struct {
	id int
	cT string // время создания
}

func main() {
	timeout := time.Duration(2) * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	superChan := make(chan Ttype, 100)
	go func() {
		taskCreturer(ctx, superChan)
		close(superChan)
	}()

	doneTasks := make(map[int]error)
	undoneTasks := make(map[int]error)
	var wg1 sync.WaitGroup
	var lock sync.Mutex

	for t := range superChan {
		wg1.Add(1)
		go taskWorkerAndSorter(&wg1, &lock, doneTasks, undoneTasks, t)
	}
	wg1.Wait()

	fmt.Println("Errors:")
	for id, err := range undoneTasks {
		fmt.Println(id, err)
	}

	fmt.Println("Done tasks:")
	for id := range doneTasks {
		fmt.Println(id)
	}
}

func taskWorkerAndSorter(wg1 *sync.WaitGroup, lock *sync.Mutex, doneTasks map[int]error, undoneTasks map[int]error, a Ttype) {
	tt, err := time.Parse(time.RFC3339, a.cT)
	if tt.After(time.Now().Add(-20*time.Second)) && err == nil {
		lock.Lock()
		doneTasks[a.id] = err
		lock.Unlock()
	} else {
		lock.Lock()
		undoneTasks[a.id] = fmt.Errorf("task id %d, error %s", a.id, err)
		lock.Unlock()
	}
	wg1.Done()
}

func taskCreturer(ctx context.Context, a chan Ttype) {
	var wg sync.WaitGroup
	for {
		select {
		case <-ctx.Done():
		default:
			b := time.Now()
			wg.Add(1)
			go func() {
				defer wg.Done()
				ct := b.Format(time.RFC3339)
				if b.Nanosecond()%2 > 0 {
					ct = "Some error occured"
				}
				a <- Ttype{
					id: int(b.Unix()),
					cT: ct,
				}
			}()
			continue
		}
		break
	}
	wg.Wait()
}
