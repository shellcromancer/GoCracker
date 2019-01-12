package main

import (
    //    "fmt"
    "github.com/dstindiess/GoCracker/md5crypt"
    "log"
    "sync"
    "sync/atomic"
    "time"
)

var producers sync.WaitGroup
var consumers sync.WaitGroup
var alphabet = []rune{'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z'}
var StartTime = time.Now()
var cracks, pws int64

const CRACKING_LEN = 4

func Produce(ch chan<- string, character_set []rune, k int) {
    defer producers.Done()

    for i := 0; i < k; i++ {
        produce(ch, character_set, "", i+1)
    }
}

func produce(ch chan<- string, character_set []rune, prefix string, k int) {
    if k == 0 {
        ch <- prefix
        pws += 1
        return
    }
    for _, runeValue := range character_set {
        newPrefix := prefix + string(runeValue)
        produce(ch, alphabet, newPrefix, k-1)
    }
}

func consume(ch <-chan string) {
    defer consumers.Done()

    for pass := range ch {
        //Hash, check
        hash := md5crypt.Hash([]byte(pass), []byte("hfT7jp2q"), []byte("$1$"))
        //hash := "abcdef"
        atomic.AddInt64(&cracks, 1)

        switch hash {
        case "Vd693H7jroUcmcZV3RJ1S/":
            success(pass)
        case "/fBukIHL391IspS.gX/Eh1":
            success(pass)
        case "/qzSQ8SeCQEdSg47A7VPJ/":
            success(pass)
        case "8rU1qXqPJfSiwL8uts982.":
            success(pass)
        case "yKkGOHLs7BZiNuh03um670":
            success(pass)
        }
    }
}

func success(pw string) {
    elapsed := time.Since(StartTime)
    log.Println("[!] Success: The hash corresponds with ", string(pw))
    log.Println("[!] Cracking took ", elapsed)
    log.Println("[!] Total passwords tried: ", atomic.LoadInt64(&cracks))
}

func main() {
    ch := make(chan string, 321272406) // Buffered Channel

    for i := 0; i < 4; i++ {
        consumers.Add(1)
        go consume(ch)
    }

    var runeCharArr []rune
    for i := 0; i < 24; i += 6 {
        runeCharArr = nil
        for j := 0; j < 6; j++ {
            runeCharArr = append(runeCharArr, alphabet[i+j])
        }
        producers.Add(1)
        go Produce(ch, runeCharArr, CRACKING_LEN)
    }
    runeCharArr = nil

    runeCharArr = append(runeCharArr, alphabet[24])
    runeCharArr = append(runeCharArr, alphabet[25])
    producers.Add(1)
    Produce(ch, runeCharArr, CRACKING_LEN)

    producers.Wait()
    close(ch)
    elapsed := time.Since(StartTime)
    log.Println("Passwords Generated: ", pws)
    log.Println("Generation took ", elapsed)

    consumers.Wait()
}
