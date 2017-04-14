package main

import (
    "bufio"
    "log"
    "strings"
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

type Address struct{
        Ip string
        Port string
}

var (
    successor Address = Address{"",""}
    predecessor Address = Address{"",""}
    own_Address = Address{"",""}
    store map[string]string = make(map[string]string)
    onlyOne bool
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

func (t *KeySpace)findSuccessor(string key, reply *Address) error{
    if successor==own_Address {
        store[keyVal.Key] = keyVal.Val
        return nil
    }
    else {

    }
}

// func insertWord(string word, string meaning, stringList synonyms, Type type) error - The error returns either AlreadyExists or OtherServerSideError 

func (t *KeySpace)Insert(keyVal KeyVal, reply *string) error{
    store[keyVal.Key] = keyVal.Val
    return nil
}

func (t *KeySpace)callInsert(keyVal KeyVal, reply *string) error{
    if successor==own_Address {
        store[keyVal.Key] = keyVal.Val
        return nil
    }
    else {

    }
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

func (t *KeySpace)callRemove(key string, reply *string) error{
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

func (t *KeySpace)callGet(key string, val *string) error{
    v, ok := store[key]
    if !ok {
        return errors.New("Key not found")
    }
    *val = v
    return nil
}




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
    //./main create [portToListen]
    // ./main [ip_someNode] [port_someNode] [portToListen]
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
    } else { 
        addrs,_ := net.InterfaceAddrs()
        successor = Address{addrs[1].String(),os.Args[2]}
        own_Address = Address{addrs[1].String(),os.Args[2]}
    }
    go as_server_for_others()
    fmt.Println("hihi")

    for true {
        string_return := ""
        reader := bufio.NewReader(os.Stdin)
        fmt.Print("Enter Command(Insert/Remove/Get):") 
        text, _ := reader.ReadString('\n')
        text = strings.TrimSpace(text)
        // fmt.Println("hihi",text)
        if text=="Insert" { 
            var keyVal_obj KeyVal
            fmt.Print("Enter Key:") 
            text, _ := reader.ReadString('\n')
            text = strings.TrimSpace(text)
            keyVal_obj.Key = text
            fmt.Print("Enter Val:") 
            text, _ = reader.ReadString('\n')
            text = strings.TrimSpace(text)
            keyVal_obj.Val = text
            fmt.Println(keyVal_obj)
            err := keySpace1.Insert(keyVal_obj, &string_return)
            if err != nil {
                fmt.Println("error:", err)
            } else {
                fmt.Println("KeyVal Inserted")
            }
        } else if text=="Remove" {
            fmt.Print("Enter key:")
            text, _ := reader.ReadString('\n')
            text = strings.TrimSpace(text)

            err := keySpace1.Remove(text, &string_return)
            if err != nil {
                fmt.Println("error:", err)
            } else {
                fmt.Println("Key Removed")
            }
            
        } else if text=="Get" {
            fmt.Print("Enter key_string:")
            text, _ := reader.ReadString('\n')
            text = strings.TrimSpace(text)
            err := keySpace1.Get(text, &string_return)
            if err != nil {
                fmt.Println("error:", err)
            } else if string_return==""{
                fmt.Println("error:No val")
            } else {
                fmt.Println(string_return)
            }
        }
    }
}
