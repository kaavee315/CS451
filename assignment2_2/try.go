// For this programming assignment, you will build a remotely accessible Dictionary in Go. 
// struct Word {
//          String word;
//          String meaning;
//          Word* synonyms;\
//           Enum {Noun, Verb, Adjective} - denoting a type. 
// };

// The dictionary will export the following procedures:
// func insertWord(String word, String meaning, StringList synonyms, Type type) error - The error returns either AlreadyExists or OtherServerSideError
// func removeWord(Word) error - The error contains either UnknownWord or OtherServerSideError.  Make this procedure take 5 seconds using sleep. RemoveWord should also remove this as the synonym from all other words that has this word as its synonym.  
// func lookupWord(String, Word*) error - The error contains either UnknownWord or OtherServerSideError.
// The server should be able to accept either TCP connections @ port 5000 and HTTP Connections @ port 6000 - in other words, the server must support both. 

// You will need to write three clients:
// a. A synchronous client using TCP in Go to  take user input to lookup, delete and insert words.  The removeWord will block for 5 seconds because of the server implementation. In each case the client must print out either Success Inserting word, Success Removing word or the word, its meaning, its type and its synonym strings. 
// b. A asynchronous client using HTTP in Go that does the same as a. Demo the async nature of the call you make by printing the your name and roll number in a loop till such time you block for the return value.
// c. A synchronous TCP client in Python that talks to this server using JSON encodings to do the same things as a. 

// Submit one TAR file that has:
// a. server.go
// b. sync-client.go
// c. async-client.go
// d. sync-python-client.py
// e. README - explaining precisely how to run all this.  We will blindly follow your instructions and if something does not work, you're sunk, so be careful with this. Paths are critical and any assumptions you have made about paths on your local machine may or may not work on the test machine. Assume that Go is available in /usr/local/go and Python in /usr/bin/python 

// A useful link you may want to go through before starting out is:
// https://jan.newmarch.name/go/



package server

import (
  // "bufio"
  // "log"
  // "runtime"
  // "sync"
  "log"
  "strings"
  "net/http"
  "net"
  "os"
  "fmt"
  "errors"
  "time"
)

type type_word int
  const (
    Noun type_word  = iota
    Verb
    Adjective
  )

type Word struct{
    word string
    meaning string
    synonyms []Word 
    type_w type_word
}

var (
  dict []Word
  s = "a"
    
)

// func insertWord(String word, String meaning, StringList synonyms, Type type) error - The error returns either AlreadyExists or OtherServerSideError 
func insertWord(word1, meaning1 string, synonyms1 []Word,type_w1 type_word ) error{
  found := false
  for _, ele:= range dict {
    if word1==ele.word{
      found=true
    }
  }
  if found{
    return errors.New("AlreadyExists")
  }
  word_obj := Word{word: word1, meaning:meaning1, synonyms:synonyms1 ,type_w:type_w1}
  dict = append(dict, word_obj)
  return nil
}

// func removeWord(Word) error - The error contains either UnknownWord or OtherServerSideError.  Make this procedure take 5 seconds using sleep. RemoveWord should also remove this as the synonym from all other words that has this word as its synonym. 
func removeWord(word1 string) error{
  start := time.Now()
  found := false
  for i, ele:= range dict {
    if word1==ele.word {
      if i<(len(dict)-1){
        dict[len(dict)-1], dict[i] = dict[i], dict[len(dict)-1]
      }
      dict=dict[:len(dict)-1]
      found=true
    } else {
      for j,ele2:= range ele.synonyms {
        if word1==ele2.word {
          if j<(len(ele.synonyms)-1){
            ele.synonyms[len(ele.synonyms)-1], ele.synonyms[i] = ele.synonyms[i], ele.synonyms[len(ele.synonyms)-1]
          }
          ele.synonyms=ele.synonyms[:len(ele.synonyms)-1]
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
func lookupWord(word1 string, word_obj *Word) error{
  found := false
  for _, ele:= range dict {
    if word1==ele.word {
      *word_obj = Word(ele)
      found=true
    }
  }
  if !found {
    return errors.New("UnknownWord")
  }
  return nil
}

func checkError(err error) {
  if err != nil {
    fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
    os.Exit(1)
  }
}

func handleHTTPClient(w http.ResponseWriter, r *http.Request) {
  r.ParseForm()  // parse arguments, you have to call this by yourself
  fmt.Println(r.Form)  // print form information in server side
  fmt.Println("path", r.URL.Path)
  fmt.Println("scheme", r.URL.Scheme)
  fmt.Println(r.Form["url_long"])
  for k, v := range r.Form {
      fmt.Println("key:", k)
      fmt.Println("val:", strings.Join(v, ""))
  }
  fmt.Fprintf(w, "Hello astaxie!") // send data to client side
}

// The server should be able to accept either TCP connections @ port 5000 and HTTP Connections @ port 6000 - in other words, the server must support both. 
func http_server(){
  http.HandleFunc("/", handleHTTPClient) // set router
  err := http.ListenAndServe(":6000", nil) // set listen port
  if err != nil {
      log.Fatal("ListenAndServe: ", err)
  }
}

func handleClient(conn net.Conn) {
  // close connection on exit
  defer conn.Close()

  var buf [512]byte
  for {
    // read upto 512 bytes
    n, err := conn.Read(buf[0:])
    if err != nil {
      return
    }

    // write the n bytes read
    _, err2 := conn.Write(buf[0:n])
    if err2 != nil {
      return
    }
  }
}

func tcp_server(){
  const (
    CONN_HOST = "localhost"
    CONN_PORT = "5000"
    CONN_TYPE = "tcp"
  ) 
  listener, err := net.Listen(CONN_TYPE, CONN_HOST+":"+CONN_PORT)
  defer listener.Close()
  if err != nil {
      fmt.Println("Error listening:", err.Error())
      os.Exit(1)
  }
  for {
    conn, err := listener.Accept()
    if err != nil {
      continue
    }
    go handleClient(conn)
  }
}

func main() {
  go tcp_server()
  http_server()
  // insertWord("kaavee","karan",nil,Verb)
  // removeWord("kaavee")
}
