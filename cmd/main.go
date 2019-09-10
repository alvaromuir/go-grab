package main

import (
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/alvaromuir/go-grab/actions"
)


func worker(tasksCh <-chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		task, ok := <-tasksCh
		if !ok {
			return
		}
		d := time.Duration(task) * time.Millisecond
		time.Sleep(d)
		url := fmt.Sprintf("<http://url/to/download/from/>%s.ext", strconv.Itoa(task+2000))
		// fmt.Println("processing task", task)
		fmt.Println("processing url", url)
		actions.DownloadFromURL(url, "</path/to/download/to>")
	}
}

func pool(wg *sync.WaitGroup, workers, tasks int) {
	tasksCh := make(chan int)

	for i := 0; i < workers; i++ {
		go worker(tasksCh, wg)
	}

	for i := 0; i < tasks; i++ {
		tasksCh <- i
	}

	close(tasksCh)
}

func main() {
	var wg sync.WaitGroup
	wg.Add(100)
	go pool(&wg, 100, 2000)
	wg.Wait()
}