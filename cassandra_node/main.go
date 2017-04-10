package main

import (
    // "bufio"
    "log"
    // "strings"
    // "runtime"
    // "sync"
    "net"
    "net/rpc"
    "net/rpc/jsonrpc"
    // "net/http"
    "os"
    "fmt"
    "errors"
    "hash/fnv"
    // "time"
)

var (
    successor_ip string = ""
    successor_port string = ""
    predecessor_ip string = ""
    predecessor_port string = ""
    store map[string]string = make(map[string]string)
)

type KeyVal struct{
        Key string
        Val string
}

type KeySpace int

func hash(s string) uint32 {
        h := fnv.New32a()
        h.Write([]byte(s))
        return h.Sum32()
}

func (t *KeySpace)findSuccessor()


// func insertWord(string word, string meaning, stringList synonyms, Type type) error - The error returns either AlreadyExists or OtherServerSideError 
func (t *KeySpace)Insert(keyVal KeyVal, reply *Word) error{
    store[keyVal.key] = keyVal.val
    return nil
}

// func removeWord(Word) error - The error contains either UnknownWord or OtherServerSideError.    Make this procedure take 5 seconds using sleep. RemoveWord should also remove this as the synonym from all other words that has this word as its synonym. 
func (t *KeySpace)Remove(key string, reply *string) error{
    _, ok := store[key]
    if !ok {
        return errors.New("Key not found")
    }
    delete(store, key)
    return nil
}

// func lookupWord(string, Word*) error - The error contains either UnknownWord or OtherServerSideError.
func (t *KeySpace)Get(key string, val *string) error{
    v, ok := store[key]
    if !ok {
        return errors.New("Key not found")
    }
    *val = v
    return nil
}

func findSuccessor


func as_server_for_others() {
    var l net.Listener
    var e error
    if os.Args[1]=="create" {
        l, e = net.Listen("tcp", ":" + os.Args[2])
    } else {
        l, e = net.Listen("tcp", ":" + os.Args[3])
    }
    if e != nil {
        log.Fatal("listen error:", e)
    }
    for {
        conn, err := l.Accept()
        if err != nil {
          log.Printf("accept error: %s", conn)
          continue
        }
        log.Printf("connection started: %v", conn.RemoteAddr())
        go jsonrpc.ServeConn(conn)
    }
}


func main() {
    keySpace1 := new(KeySpace)
    rpc.Register(keySpace1)

    if os.Args[1]!="create" {
        conn, err := net.Dial("tcp", os.Args[1]+":"+os.Args[2])
        if err != nil {
            log.Fatal("Connectiong:", err)
        }
        client := jsonrpc.NewClient(conn)
        err = client.Call("keySpace1.findSuccessor", nil, nil)
        if err != nil {
            fmt.Println("error:", err)
        } else {
            fmt.Println("Word Removed")
        }
    }
    go as_server_for_others()
    fmt.Println("hihi")
}
