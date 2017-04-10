package main

import (
    "fmt"
    "net"
    )

func main() {
    // ifaces, _ := net.Interfaces()
    // for _, i := range ifaces {
    //     addrs, _ := i.Addrs()
    //     fmt.Println("addrs - ",addrs)
    //     for _, addr := range addrs {
    //         var ip net.IP
    //         switch v := addr.(type) {
    //         case *net.IPNet:
    //                 ip = v.IP
    //         case *net.IPAddr:
    //                 ip = v.IP
    //         }
    //         fmt.Println("ip - ",ip)
    //     }
    // }
    addrs,_ := net.InterfaceAddrs() 
    fmt.Println(addrs[1].String())
}