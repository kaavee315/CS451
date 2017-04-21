package main

import (
    // "bufio"
    // "log"
    // "strings"
    // "runtime"
    // "sync"
    // "net"
    // "net/rpc"
    // "net/rpc/jsonrpc"
    // "net/http"
    // "os"
    // "fmt"
    "errors"
    // "hash/fnv"
    // "time"
)

type KeyVal struct{
        Key string
        Val string
}

type KeySpace int

// func insertWord(string word, string meaning, stringList synonyms, Type type) error - The error returns either AlreadyExists or OtherServerSideError 

func (t *KeySpace)Insert(keyVal KeyVal, reply *string) error{
    store[keyVal.Key] = keyVal.Val
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