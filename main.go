package main

import (
	"fmt"
	"crypto/md5"
    "io"
    "sync"
    "log"
    "time"
    "os"
)

var wg sync.WaitGroup
var target = "a753b48c75094ded57309c0b1e84b458" // "abcdef" x1000
var StartTime = time.Now()


func Produce(ch chan<- string, character_set []rune, k int) {
    produce(ch, character_set, "", k)
    close(ch)
}

func produce(ch chan<- string, character_set []rune, passwd string, k int) {
    if k == 0 {
        ch <- passwd
    	return
    }

    for _, runeValue := range character_set {
    	newPasswd := passwd + string(runeValue)
    	produce(ch, character_set, newPasswd, k - 1)
     }
}

func md5plus(text string, cost int) string {
	for i := 0; i < cost; i++ {
        hash := md5.New()
        io.WriteString(hash, text)
		text = fmt.Sprintf("%x", hash.Sum(nil))
	}
	return text
}

/*
func produce() {
}
*/

func consume(ch <-chan string) {
    defer wg.Done()
    //salt := ""

    for pass := range ch {
        //Hash, check
        hash := md5plus(pass, 1000)
        if hash == target {
            elapsed := time.Since(StartTime)
            log.Printf("[!] Success: The hash corresponds with %s.\n", string(pass))
            log.Printf("Cracking took %s", elapsed)
            os.Exit(0)        
        }
    }
}

func main() {
    ch := make(chan string, 100) // Buffered Channel
	var alphabet = []rune {'a','b','c','d','e','f','g','h','i','j','k','l','m','n','o','p','q','r','s','t','u','v','w','x','y','z'}

    for i := 0; i < 5; i++ {
        wg.Add(1)
        go consume(ch)
    }

    Produce(ch, alphabet, 6)
    wg.Wait()
}