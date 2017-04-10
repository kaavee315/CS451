// Instructions to run the script:-
// 1. Copy the file wordCount.go to a folder in your source directory in workspace of Go(let's say src/CS451/130050019_wordCount/wordCount.go for instructions ahead)
// 2. Compile using go install CS451/130050019_wordCount
// 3. Run the bin/130050019_wordCount $file_path $string from workspace of Go
// The output will be of format "Count = N", where N is the number of lines which contains $string in the file located at $file_path



package main

import (
  "bufio"
  "log"
  "fmt"
  "os"
  "strings"
  "runtime"
  "sync"
)

var (
	count int = 0
)


//get the parallelism to get number of threads to be created
func MaxParallelism() int {
    maxProcs := runtime.GOMAXPROCS(0)
    numCPU := runtime.NumCPU()
    if maxProcs < numCPU {
        return maxProcs
    }
    return numCPU
}

//worker thread handler
func thread_handler(channel_string, channel_done chan string, count_mux sync.Mutex, find string) {
	for{
		s,ok := <- channel_string
		if !ok {
			channel_done <- "done"
			return
		}
		if strings.Contains(s, find) {
			count_mux.Lock()
			count++
			count_mux.Unlock()
		}
	}
}

func main() {

  // open a file
  if file, err := os.Open(os.Args[1]); err == nil {
    // make sure it gets closed
    defer file.Close()


    find := os.Args[2]
    channel_string := make(chan string)
    channel_done := make(chan string)
    max_threads := MaxParallelism()
    var count_mux sync.Mutex

    //loop to spawn the threads
    for i:=1;i<max_threads;i++ {
    	go thread_handler(channel_string, channel_done, count_mux, find)
    }

    // create a new scanner and read the file line by line
    scanner := bufio.NewScanner(file)
    count=0
    for scanner.Scan() {
        s := string(scanner.Text())
        channel_string <- s
    }

    close(channel_string)

    // check for errors
    if err = scanner.Err(); err != nil {
      log.Fatal(err)
    }

    //check for the threads to close
    for i:=1;i<max_threads;i++ {
    	<- channel_done
    }

    fmt.Println("Count = ",count)

  } else {
    log.Fatal(err)
  }
}
