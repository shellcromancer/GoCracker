package main

import (
    "crypto/md5"
    //    "fmt"
    //    "io"
    "log"
    "os"
    "sync"
    "time"
)

var wg sync.WaitGroup

var magic = "$1$"
var salt = "hfT7jp2q"
var md5cryptTarget = "8rU1qXqPJfSiwL8uts982."
var md5CryptSwaps = [16]int{12, 6, 0, 13, 7, 1, 14, 8, 2, 15, 9, 3, 5, 10, 4, 11}
var alphabet = []rune{'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z'}
var StartTime = time.Now()

const itoa64 = "./0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

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
        produce(ch, character_set, newPasswd, k-1)
    }
}

func md5crypt(password, salt, magic []byte) []byte {
    //Initialization
    intermediate := md5.New()
    alternate := md5.New()

    intermediate.Write(password)
    intermediate.Write(magic)
    intermediate.Write(salt)

    alternate.Write(password)
    alternate.Write(salt)
    alternate.Write(password)

    for i, mixin := 0, alternate.Sum(nil); i < len(password); i++ {
        intermediate.Write([]byte{mixin[i%16]})
    }

    for i := len(password); i != 0; i >>= 1 {
        if i&1 == 0 {
            intermediate.Write([]byte{password[0]})
        } else {
            intermediate.Write([]byte{0})
        }
    }

    final := intermediate.Sum(nil)
    // Loop/ Stetching
    for i := 0; i < 1000; i++ {
        hasher := md5.New()

        if i&1 == 0 {
            hasher.Write(final)
        } else {
            hasher.Write(password)
        }
        if i%3 != 0 {
            hasher.Write(salt)
        }
        if i%7 != 0 {
            hasher.Write(password)
        }
        if i&1 == 0 {
            hasher.Write(password)
        } else {
            hasher.Write(final)
        }
        final = hasher.Sum(nil)
    }

    //Finalization
    result := make([]byte, 0, 22)
    v := uint(0)
    bits := uint(0)
    for _, i := range md5CryptSwaps {
        v |= (uint(final[i]) << bits)
        for bits = bits + 8; bits > 6; bits -= 6 {
            result = append(result, itoa64[v&0x3f])
            v >>= 6
        }
    }

    return append(result, itoa64[v&0x3f])
}

func consume(ch <-chan string) {
    defer wg.Done()

    for pass := range ch {
        //Hash, check
        hash := md5crypt([]byte(pass), []byte(salt), []byte(magic))
        if string(hash) == md5cryptTarget {
            elapsed := time.Since(StartTime)
            log.Printf("[!] Success: The hash corresponds with %s.\n", string(pass))
            log.Printf("Cracking took %s", elapsed)
            os.Exit(0)
        }
    }
}

func main() {
    ch := make(chan string, 100) // Buffered Channel

    for i := 0; i < 5; i++ {
        wg.Add(1)
        go consume(ch)
    }

    Produce(ch, alphabet, 6)

    wg.Wait()
}
