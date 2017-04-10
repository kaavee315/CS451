// For this programming assignment, you will build a remotely accessible Dictionary in Go. 
// struct Word {
//                    String word;
//                    String meaning;
//                    Word* synonyms;\
//                     Enum {Noun, Verb, Adjective} - denoting a type. 
// };

// The dictionary will export the following procedures:
// func insertWord(String word, String meaning, StringList synonyms, Type type) error - The error returns either AlreadyExists or OtherServerSideError
// func removeWord(Word) error - The error contains either UnknownWord or OtherServerSideError.    Make this procedure take 5 seconds using sleep. RemoveWord should also remove this as the synonym from all other words that has this word as its synonym.    
// func lookupWord(String, Word*) error - The error contains either UnknownWord or OtherServerSideError.
// The server should be able to accept either TCP connections @ port 5000 and HTTP Connections @ port 6000 - in other words, the server must support both. 

// You will need to write three clients:
// a. A synchronous client using TCP in Go to    take user input to lookup, delete and insert words.    The removeWord will block for 5 seconds because of the server implementation. In each case the client must print out either Success Inserting word, Success Removing word or the word, its meaning, its type and its synonym strings. 
// b. A asynchronous client using HTTP in Go that does the same as a. Demo the async nature of the call you make by printing the your name and roll number in a loop till such time you block for the return value.
// c. A synchronous TCP client in Python that talks to this server using JSON encodings to do the same things as a. 

// Submit one TAR file that has:
// a. server.go
// b. sync-client.go
// c. async-client.go
// d. sync-python-client.py
// e. README - explaining precisely how to run all this.    We will blindly follow your instructions and if something does not work, you're sunk, so be careful with this. Paths are critical and any assumptions you have made about paths on your local machine may or may not work on the test machine. Assume that Go is available in /usr/local/go and Python in /usr/bin/python 

// A useful link you may want to go through before starting out is:
// https://jan.newmarch.name/go/



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
    "net/http"
    // "os"
    // "fmt"
    "errors"
    "time"
)

type Type_word int
    const (
        Noun Type_word    = iota
        Verb
        Adjective
    )

type Word struct{
        Word_s string
        Meaning string
        Synonyms []Word 
        Type_w Type_word
}

var (
    dict []Word
)

type Dictionary int

// func insertWord(String word, String meaning, StringList synonyms, Type type) error - The error returns either AlreadyExists or OtherServerSideError 
func (t *Dictionary)InsertWord(word_obj1 Word, reply *Word) error{
    found := false
    for _, ele:= range dict {
        if word_obj1.Word_s==ele.Word_s{
            found=true
        }
    }
    if found{
        return errors.New("AlreadyExists")
    }
    word_obj := Word(word_obj1)
    dict = append(dict, word_obj)
    return nil
}

// func removeWord(Word) error - The error contains either UnknownWord or OtherServerSideError.    Make this procedure take 5 seconds using sleep. RemoveWord should also remove this as the synonym from all other words that has this word as its synonym. 
func (t *Dictionary)RemoveWord(word1 string, reply *Word) error{
    start := time.Now()
    found := false
    // println("yoyo0 ", len(dict))
    for i:=0;i<len(dict);i++ {
        // fmt.Println("yoyo in main for ", i, dict[i])
        if word1==dict[i].Word_s {
            if i<(len(dict)-1){
                dict[len(dict)-1], dict[i] = dict[i], dict[len(dict)-1]
            }
            dict=dict[:len(dict)-1]
            i=i-1
            found=true
        } else {
            for j:=0;j<len(dict[i].Synonyms);j++ {
                if word1==dict[i].Synonyms[j].Word_s {
                    if j<(len(dict[i].Synonyms)-1){
                        dict[i].Synonyms[len(dict[i].Synonyms)-1], dict[i].Synonyms[j] = dict[i].Synonyms[j], dict[i].Synonyms[len(dict[i].Synonyms)-1]
                    }
                    dict[i].Synonyms=dict[i].Synonyms[:len(dict[i].Synonyms)-1]
                    j=j-1
                }
            }
        }
    }
    time.Sleep((5*time.Second) - time.Since(start))
    if !found {
        return errors.New("UnknownWord")
    }
    return nil
}

// func lookupWord(String, Word*) error - The error contains either UnknownWord or OtherServerSideError.
func (t *Dictionary)LookupWord(word1 string, word_obj *Word) error{
    found := false
    for _, ele:= range dict {
        if word1==ele.Word_s {
            *word_obj = Word(ele)
            found=true
        }
    }
    if !found {
        return errors.New("UnknownWord")
    }
    return nil
}

func tcp_handler() {

    dictionary := new(Dictionary)
    // Listen for incoming tcp packets on specified port.
    l, e := net.Listen("tcp", ":5000")
    if e != nil {
        log.Fatal("listen error:", e)
    }
    rpc.Register(dictionary)

    // This statement links rpc server to the socket, and allows rpc server to accept
    // rpc request coming from that socket.
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
    go tcp_handler()
    dictionary := new(Dictionary)
    word:= Word{"kaavee", "karan", nil, Noun}
    dictionary.InsertWord(word, nil) 
    word= Word{"ritesha", "ritesh", nil, Verb}
    dictionary.InsertWord(word, nil)
    word= Word{"hi", "bye", nil, Verb}
    dictionary.InsertWord(word, nil)
    rpc.Register(dictionary)
    rpc.HandleHTTP()
    l, e := net.Listen("tcp", ":6000")
    if e != nil {
        log.Fatal("listen error:", e)
    }
    http.Serve(l, nil)
}
