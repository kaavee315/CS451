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
    "bufio"
    "log"
    "strings"
    // "runtime"
    // "sync"
    // "net"
    "net/rpc"
    // "net/http"
    "os"
    "fmt"
    // "errors"
    // "time"
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

func main() {
    // Asynchronous call
    client, err := rpc.DialHTTP("tcp", os.Args[1]+":6000")
    if err != nil {
        log.Fatal("dialing:", err)
    }

    for true {
        word := new(Word)
        reader := bufio.NewReader(os.Stdin)
        fmt.Print("Enter Command(Insert/Remove/Lookup):") 
        text, _ := reader.ReadString('\n')
        text = strings.TrimSpace(text)
        // fmt.Println("hihi",text)
        if text=="Insert" {
            var word_obj Word 
            fmt.Print("Enter Word:") 
            text, _ := reader.ReadString('\n')
            text = strings.TrimSpace(text)
            word_obj.Word_s = text
            text = strings.TrimSpace(text)
            fmt.Print("Enter Meaning:") 
            text, _ = reader.ReadString('\n')
            text = strings.TrimSpace(text)
            word_obj.Meaning = text
            no_synonym:=false
            for !no_synonym {
                fmt.Print("Enter Synonym(If none left enter 0):")
                text, _ := reader.ReadString('\n')
                text = strings.TrimSpace(text)
                if text!="0" {
                    word_obj.Synonyms = append(word_obj.Synonyms, Word{Word_s:text})
                } else {
                    no_synonym=true
                }
            }
            fmt.Print("Enter Type(Noun, Adjective, Verb):") 
            text, _ = reader.ReadString('\n')
            text = strings.TrimSpace(text)
            if text=="Noun" {
                word_obj.Type_w = Noun
            } else if text=="Verb" {
                word_obj.Type_w = Verb
            } else {
                word_obj.Type_w = Adjective
            }
            fmt.Println(word_obj)
            InsertCall := client.Go("Dictionary.InsertWord", word_obj, word, nil)
            done := false
            for !done {
                select {
                    case replyCall := <-InsertCall.Done:
                            if replyCall.Error==nil {
                                fmt.Println("Word Inserted")
                            } else {
                                fmt.Println(replyCall.Error)
                            }             
                            done=true
                    default:
                            fmt.Println("Karan:130050019")
                }
            }
        } else if text=="Remove" {
            fmt.Print("Enter word_string:")
            text, _ := reader.ReadString('\n')
            text = strings.TrimSpace(text)

            removeCall := client.Go("Dictionary.RemoveWord", text, word, nil)
            done := false
            for !done {
                select {
                    case replyCall := <-removeCall.Done:
                            if replyCall.Error==nil {
                                fmt.Println("Word Removed")
                            } else {
                                fmt.Println(replyCall.Error)
                            }         
                            done=true
                    default:
                            fmt.Println("Karan:130050019")
                }
            }
        } else if text=="Lookup" {
            fmt.Print("Enter word_string:")
            text, _ := reader.ReadString('\n')
            text = strings.TrimSpace(text)
            lookUpCall := client.Go("Dictionary.LookupWord", text, word, nil)
            done := false
            for !done {
                select {
                    case replyCall := <-lookUpCall.Done:
                            if replyCall.Error==nil {
                                fmt.Println(word)
                            } else if word==nil{
                				fmt.Println("error:No word")
            				} else {
                                fmt.Println(replyCall.Error)
                            }
                            done=true
                    default:
                            fmt.Println("Karan:130050019")
                }
            }
        }


        // lookUpCall := client.Go("Dictionary.LookupWord", "kaavee", word, nil)
        // <-lookUpCall.Done
        // fmt.Println(word.Word_s,":",word.Meaning)
    }
}
