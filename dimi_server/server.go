package main

import (
	"errors"
	
	"log"
	"net"
	"net/http"
	"net/rpc"
	"net/rpc/jsonrpc"
	"time"
	"fmt"
	

)

type myRPCServer struct {
   *rpc.Server
}

func (r *myRPCServer) ServeHTTP(w http.ResponseWriter, req *http.Request) {
   log.Printf("established http connection with %s.\n",req.RemoteAddr)
   r.Server.ServeHTTP(w,req)
   log.Printf("closed http connection with %s",req.RemoteAddr)
}



func (r *myRPCServer) ServeCodec(codec rpc.ServerCodec, conn net.Conn ) {
	
	r.Server.ServeCodec(codec)
	log.Printf("closed tcp connection with %s",conn.RemoteAddr())
	
}

func (r *myRPCServer) HandleHTTP(rpcPath, debugPath string) {
    http.Handle(rpcPath, r)
}


func (r *myRPCServer) Register(rcvr interface{}) error {
	return r.Server.Register(rcvr)
}


type Type int
const ( 
     Noun Type = iota+1 
     Verb 
     Adjective  
) 

func (t Type) String() string {
    s:=""
    if t==Noun {
    	s+="Noun"
    } else if t==Verb {
    	s+="Verb"
	} else {
		s+="Adjective"
	}
	return s
}

type Word struct{
	Word string
	Meaning string
	Synonyms map[string]Word
	Wordtype Type
}

var Dictionary map[string]Word= make(map[string]Word)

type DictService int

func (t *DictService) InsertWord(insword Word, reply *int) error {
	if _, ok := Dictionary[insword.Word]; ok {
    	return errors.New("Already Present")
	} else {
		Dictionary[insword.Word]=insword
		for key := range insword.Synonyms {
			if _,ok := Dictionary[key]; ok {
				if _,ok := Dictionary[key].Synonyms[insword.Word]; !ok {
					Dictionary[key].Synonyms[insword.Word]= insword
				}
			}
		}
		
	}
	return nil
}

func (t *DictService) RemoveWord(remword string, reply *int) error {

	time.Sleep(5 * time.Second)
	if _, ok := Dictionary[remword]; ok {
    	delete(Dictionary, remword)
    	for key := range Dictionary {
    		if _,ok := Dictionary[key].Synonyms[remword]; ok {
    			delete(Dictionary[key].Synonyms,remword)
    		}

    	}
	} else {
		return errors.New("Word not present")
	}
	return nil
}

func (t *DictService) LookupWord(lookword string, reply *Word) error {
	if val, ok := Dictionary[lookword]; ok {
    	*reply = val
	} else {
		return errors.New("Word not present")
	}
	return nil
}

func main() {

	dictservice := new(DictService)
	server := &myRPCServer{rpc.NewServer()}
  
	server.Register(dictservice)
	server.HandleHTTP(rpc.DefaultRPCPath, rpc.DefaultDebugPath)

	go func() {
		l, e := net.Listen("tcp", ":6000")
		if e != nil {
			log.Fatal("http server listen error: ", e)
		} else {
			log.Printf("http server listening on port 6000.")
		}
		http.Serve(l, nil)

	}()

	go func() {
		l, e := net.Listen("tcp", ":5000")
		if e != nil {
			log.Fatal("tcp server listen error: ", e)
		} else {
			log.Printf("tcp server listening on port 5000.")
		}

		for {
			if conn, err := l.Accept(); err != nil {
				log.Fatal("tcp accept error: " + err.Error())
			} else {
				log.Printf("established tcp connection with with %s.",conn.RemoteAddr())
				go server.ServeCodec(jsonrpc.NewServerCodec(conn),conn)
			}
		}
	}()

	for{
	var choice int
	fmt.Scanf("%d\n", &choice)
	switch choice {

	case 1:
		for key,val := range Dictionary {
			fmt.Println(key)
			fmt.Println(val.Meaning)
			for k := range val.Synonyms {
				fmt.Println(k)
			}
			fmt.Println(val.Wordtype)
    		fmt.Println("Next word:")

    	}
    }
	}



}