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
)

type KeyVal struct{
    Key string
    Val string
}

func main() {

    num_clients,_ := strconv.Atoi(os.Args[1])
    num_servers,_ := strconv.Atoi(os.Args[3])

    clients := make([]*rpc.Client, num_clients)
    conns := make([]net.Conn, num_clients)

    for client := 0; client < num_clients; client++ {
        server := "127.0.0.1:" + strconv.Itoa(5000 + client)
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

        for i := 0; i < num_servers; i++ {

            go func() {

                fmt.Println("Reading Key: ", i)

                var val_str string
                err := clients[rand.Intn(num_clients)].Call("KeySpace.CallGet", strconv.Itoa(i), &val_str)

                if err != nil {
                    fmt.Println("error:", err)
                } else {
                    fmt.Println("Read Value: ", val_str)
                }

            }()

        }

        elapsed := time.Since(start)
        log.Printf("num_servers read requests took %s", elapsed)

    } else if os.Args[2] == "write" {

        start := time.Now()


        for i := 0; i < num_servers; i++ {

            go func() {
                var keyval KeyVal
                
                keyval.Key = strconv.Itoa(i)
                keyval.Val = strconv.Itoa(i)

                fmt.Println("Inserting Key: ", keyval.Key)
                fmt.Println("Inserting Val: ", keyval.Val)
                var reply_str string
                err := clients[rand.Intn(num_clients)].Call("KeySpace.CallInsert", keyval, &reply_str)

                if err != nil {
                    fmt.Println("error:", err)
                }
            }()

        }

        elapsed := time.Since(start)
        log.Printf("num_servers write requests took %s", elapsed)

    }
}
