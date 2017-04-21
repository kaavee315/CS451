package main

import (
    // "bufio"
    // "log"
    // "strings"
    // "runtime"
    // "sync"
    "net"
    // "net/rpc"
    // "net/rpc/jsonrpc"
    // "net/http"
    // "os"
    // "fmt"
    "errors"
    "hash/fnv"
    // "time"
    )

type Address struct{
        Ip, Port string
}

func (addr Address) to_string() string {
    return addr.Ip + ":" + addr.Port
}

func hash(s string) uint32 {
        h := fnv.New32a()
        h.Write([]byte(s))
        return h.Sum32()
}

func in_between(a string, b string, c string) bool{
    return (((hash(c)>hash(a)) && 
            (hash(b)>hash(a)) && 
            (hash(b)<=hash(c))) || 
        ((hash(c)<hash(a)) && 
            ((hash(b)>hash(a)) ||   
                (hash(b)<=hash(c)))))
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