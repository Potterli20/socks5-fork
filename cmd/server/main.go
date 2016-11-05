package main

import (
    "log"
    "net"
    //"time"
    "github.com/txthinking/socks5"
    "net/http"
    _ "net/http/pprof"
)

func main() {
    go func (){
        log.Println(http.ListenAndServe(":1094", nil))
    }()
    l, err := net.Listen("tcp", ":1090")
    if err != nil {
        log.Fatal(err)
    }
    defer func() {
        if err := l.Close(); err != nil {
            log.Println(err)
        }
    }()
    for {
        c, err := l.Accept()
        if err != nil {
            log.Fatal(err)
        }
        //if err := c.SetDeadline(time.Now().Add(60 * time.Second)); err != nil {
            //log.Println(err)
            //if err = c.Close(); err != nil {
                //log.Println(err)
            //}
            //continue
        //}
        s := socks5.NewServer(c)
        go func (){
            if err := s.Handle(); err != nil {
                log.Println(err)
            }
        }()
    }
}