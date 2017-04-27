package main

import (
    // "bufio"
    "log"
    // "strings"
    "net"
    "net/rpc/jsonrpc"
    "net/rpc"
    "os"
    "fmt"
    "math/rand"
    "strconv"
    "time"
    "sync"
)

type KeyVal struct{
    Key string
    Val string
}

func write(client *rpc.Client, i int, wg *sync.WaitGroup) {
    defer wg.Done()
    var keyval KeyVal
    
    keyval.Key = strconv.Itoa(i)
    keyval.Val = strconv.Itoa(i)

    fmt.Println("Inserting Key: ", keyval.Key)
    fmt.Println("Inserting Val: ", keyval.Val)
    var reply_str string
    err := client.Call("KeySpace.CallInsert", keyval, &reply_str)

    if err != nil {
        fmt.Println("error:", err)
    }

    print("done")
    return 
}

func main() {

    num_clients,_ := strconv.Atoi(os.Args[1])
    num_servers,_ := strconv.Atoi(os.Args[3])

    clients := make([]*rpc.Client, num_clients)
    conns := make([]net.Conn, num_clients)

    for client := 0; client < num_clients; client++ {
        server := "127.0.0.1:" + strconv.Itoa(5000 + rand.Intn(num_servers))
        var err error
        conns[client], err = net.Dial("tcp", server)

        if err != nil {
            log.Fatal("Connection: ", err)
        }

        clients[client] = jsonrpc.NewClient(conns[client])
    }
    
    if os.Args[2] == "read" {
        
        // var keyval KeyVal

        start := time.Now()

        var wg sync.WaitGroup

        n,_:= strconv.Atoi(os.Args[4])

        for i := 0; i < n; i++ {
            wg.Add(1)

            func() {
                defer wg.Done()

                fmt.Println("Reading Key: ", i)

                var val_str string
                err := clients[rand.Intn(num_clients)].Call("KeySpace.CallGet", strconv.Itoa(i), &val_str)

                if err != nil {
                    fmt.Println("error:", err)
                } else {
                    fmt.Println("Read Value: ", val_str)
                }

            }()
            wg.Wait()

        }

        elapsed := time.Since(start)
        log.Printf("num_servers read requests took %s", elapsed)

    } else if os.Args[2] == "write" {

        start := time.Now()

        var wg sync.WaitGroup
        n,_:= strconv.Atoi(os.Args[4])
        for i := 0; i < n ; i++ {
            wg.Add(1)
            write(clients[rand.Intn(num_clients)], i, &wg)

        }
        wg.Wait()

        elapsed := time.Since(start)
        log.Printf("num_servers write requests took %s", elapsed)

    }
}
