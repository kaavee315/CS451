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
    "time"
)

type Address struct{
        Ip, Port string
}

func (addr Address) to_string() string {
    return addr.Ip + ":" + addr.Port
}

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

var (
    successor Address = Address{"",""}
    predecessor Address = Address{"",""}
    own_Address = Address{"",""}
    store map[string]string = make(map[string]string)
    onlyOne bool
    keySpace1 = new(KeySpace)
)

func in_between(a string, b string, c string) bool{
    return (((hash(c)>hash(a)) && 
            (hash(b)>hash(a)) && 
            (hash(b)<=hash(c))) || 
        ((hash(c)<hash(a)) && 
            ((hash(b)>hash(a)) ||   
                (hash(b)<=hash(c)))))
}



func (t *KeySpace)FindSuccessor(key string, reply *Address) error{
    if successor==own_Address {
        *reply = Address{successor.Ip,successor.Port}
    } else {
        if in_between(own_Address.to_string(),key,successor.to_string()) {
            *reply = successor 
        } else {
            conn, err := net.Dial("tcp", successor.Ip + ":" + successor.Port)
            if err != nil {
                return err
            }
            client := jsonrpc.NewClient(conn)
            err = client.Call("KeySpace.FindSuccessor", key, reply)
            if err != nil {
                return err
            }
            conn.Close()
        }
    }
    return nil
}

func (t *KeySpace)GetPredeccesor(nothing string, reply *Address) error{
    *reply = predecessor
    return nil
}

func (t *KeySpace)Notify(addr Address, reply *Address) error{
    if (predecessor.Ip=="" || in_between(predecessor.to_string(),addr.to_string(),own_Address.to_string())) {
        predecessor = addr
    }
    fmt.Println("predecessor changed to - ",predecessor.to_string())
    if (successor.to_string()==own_Address.to_string()){
        successor=addr
    }
    return nil
}

func Stabilize() error{
    if successor.to_string() == own_Address.to_string() {
		return nil
	}
	conn, err := net.Dial("tcp", successor.Ip + ":" + successor.Port)
    if err != nil {
        return err
    }
    client := jsonrpc.NewClient(conn)
    var succ_pred Address
    err = client.Call("KeySpace.GetPredeccesor", "", succ_pred)
    if err != nil {
        return err
    }
    conn.Close()
    if(in_between(own_Address.to_string(),succ_pred.to_string(),successor.to_string())) {
        successor = succ_pred
    }

    conn, err = net.Dial("tcp", successor.Ip + ":" + successor.Port)
    if err != nil {
        return err
    }
    client = jsonrpc.NewClient(conn)
    var reply Address
    err = client.Call("KeySpace.Notify", own_Address, reply)
    if err != nil {
        return err
    }
    conn.Close()
    return nil
}

// func insertWord(string word, string meaning, stringList synonyms, Type type) error - The error returns either AlreadyExists or OtherServerSideError 

func (t *KeySpace)Insert(keyVal KeyVal, reply *string) error{
    store[keyVal.Key] = keyVal.Val
    return nil
}

func (t *KeySpace)callInsert(keyVal KeyVal, reply *string) error{
    if successor==own_Address {
        fmt.Println("successor == ownaddress = ",successor.to_string())
        err := keySpace1.Insert(keyVal, reply)
        return err
    } else {
        var to_send Address 
        err := keySpace1.FindSuccessor(keyVal.Key, &to_send)
        if err != nil {
            return err
        }
        if to_send==own_Address {
            fmt.Println("to_send == ownaddress = ",successor.to_string())
            err := keySpace1.Insert(keyVal, reply)
            return err
        } else {
            fmt.Println("to_send = ",successor.to_string())
            conn, err := net.Dial("tcp", to_send.Ip + ":" + to_send.Port)
            if err != nil {
                return err
            }
            client := jsonrpc.NewClient(conn)
            err = client.Call("KeySpace.Insert", keyVal, reply)
            if err != nil {
                return err
            }
            conn.Close()
        }
    }
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

func (t *KeySpace)callRemove(key string, reply *string) error{
    if successor==own_Address {
        fmt.Println("successor == ownaddress = ",successor.to_string())
        err := keySpace1.Remove(key,reply)
        return err
    } else {
        var to_send Address 
        err := keySpace1.FindSuccessor(key, &to_send)
        if err != nil {
            return err
        }
        if to_send==own_Address {
            fmt.Println("to_send == ownaddress = ",successor.to_string())
            err := keySpace1.Remove(key,reply)
            return err
        } else {
            fmt.Println("to_send = ",successor.to_string())
            conn, err := net.Dial("tcp", to_send.Ip + ":" + to_send.Port)
            if err != nil {
                return err
            }
            client := jsonrpc.NewClient(conn)
            err = client.Call("KeySpace.Remove", key, reply)
            if err != nil {
                return err
            }
            conn.Close()
        }
    }
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
    if successor==own_Address {
        fmt.Println("successor == ownaddress = ",successor.to_string())
        err := keySpace1.Get(key,val)
        return err
    } else {
        var to_send Address 
        err := keySpace1.FindSuccessor(key, &to_send)
        if err != nil {
            return err
        }
        if to_send==own_Address {
            fmt.Println("to_send == ownaddress = ",successor.to_string())
            err := keySpace1.Get(key,val)
            return err
        } else {
            fmt.Println("to_send = ",successor.to_string())
            conn, err := net.Dial("tcp", to_send.Ip + ":" + to_send.Port)
            if err != nil {
                return err
            }
            client := jsonrpc.NewClient(conn)
            err = client.Call("KeySpace.Get", key, val)
            if err != nil {
                return err
            }
            conn.Close()
        }
    }
    return nil
}

func externalIP() (string, error) {
    ifaces, err := net.Interfaces()
    if err != nil {
        return "", err
    }
    for _, iface := range ifaces {
        if iface.Flags&net.FlagUp == 0 {
            continue // interface down
        }
        if iface.Flags&net.FlagLoopback != 0 {
            continue // loopback interface
        }
        addrs, err := iface.Addrs()
        if err != nil {
            return "", err
        }
        for _, addr := range addrs {
            var ip net.IP
            switch v := addr.(type) {
            case *net.IPNet:
                ip = v.IP
            case *net.IPAddr:
                ip = v.IP
            }
            if ip == nil || ip.IsLoopback() {
                continue
            }
            ip = ip.To4()
            if ip == nil {
                continue // not an ipv4 address
            }
            return ip.String(), nil
        }
    }
    return "", errors.New("are you connected to the network?")
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
    keySpace2 := new(KeySpace)
    rpc.Register(keySpace2)
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

func callStabilize() {
    for t := range time.NewTicker(200000 * time.Nanosecond).C {
        Stabilize()
    }
}

func main() {
    //./main create [portToListen]
    // ./main [ip_someNode] [port_someNode] [portToListen]
    rpc.Register(keySpace1)

    if os.Args[1]!="create" {
        conn, err := net.Dial("tcp", os.Args[1]+":"+os.Args[2])
        if err != nil {
            log.Fatal("Connectiong:", err)
        }
        client := jsonrpc.NewClient(conn)
        ip,err := externalIP()
        own_Address = Address{ip,os.Args[3]}
        fmt.Println("address:- ",own_Address.to_string()," , Hash of address - ", hash(own_Address.to_string()))
        err = client.Call("KeySpace.FindSuccessor", own_Address.to_string(), &successor)
        if err != nil {
            log.Fatal("Successor not found error:", err)
        }
        fmt.Println("my successor - ",successor.to_string())
        conn.Close()
        conn, err = net.Dial("tcp", successor.Ip+":"+successor.Port)
         if err != nil {
            log.Fatal("Connectiong:", err)
        }
        client = jsonrpc.NewClient(conn)
        err = client.Call("KeySpace.Notify", own_Address, nil)
        if err != nil {
            log.Fatal("Successor not found error:", err)
        }
        conn.Close()

    } else { 
        ip,_ := externalIP()
        successor = Address{ip,os.Args[2]}
        own_Address = Address{ip,os.Args[2]}
        fmt.Println("address:- ",own_Address.to_string()," , Hash of address - ", hash(own_Address.to_string()))
    }
    go as_server_for_others()
    go callStabilize()
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
            fmt.Println("Hash of key - ",hash(text))
            fmt.Print("Enter Val:") 
            text, _ = reader.ReadString('\n')
            text = strings.TrimSpace(text)
            keyVal_obj.Val = text
            fmt.Println(keyVal_obj)
            err := keySpace1.callInsert(keyVal_obj, &string_return)
            if err != nil {
                fmt.Println("error:", err)
            } else {
                fmt.Println("KeyVal Inserted")
            }
        } else if text=="Remove" {
            fmt.Print("Enter key:")
            text, _ := reader.ReadString('\n')
            text = strings.TrimSpace(text)

            err := keySpace1.callRemove(text, &string_return)
            if err != nil {
                fmt.Println("error:", err)
            } else {
                fmt.Println("Key Removed")
            }
            
        } else if text=="Get" {
            fmt.Print("Enter key_string:")
            text, _ := reader.ReadString('\n')
            text = strings.TrimSpace(text)
            err := keySpace1.callGet(text, &string_return)
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
