package main

import (
    "crypto/md5"
    "fmt"
    "gopkg.in/cheggaaa/pb.v1"
    "log"
    "sync"
    "sync/atomic"
    "time"
)

var alphabet = []rune{'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z'}
var wg sync.WaitGroup
var StartTime = time.Now()
var len6Target = "8rU1qXqPJfSiwL8uts982."
var len8Target = "yKkGOHLs7BZiNuh03um670"
var md5CryptSwaps = [16]int{12, 6, 0, 13, 7, 1, 14, 8, 2, 15, 9, 3, 5, 10, 4, 11}

const itoa64 = "./0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

func Produce(ch chan<- string, bar *pb.ProgressBar, cracks int64, character_set []rune, k int) {
    for i := 0; i < k; i++ {
	fmt.Println(i+1)
	produce(ch, character_set, "", i+1)
    }
    //close(ch)
}

func produce(ch chan<- string, character_set []rune, prefix string, k int) {
    if (k == 0) {
	ch <- prefix
	return
    }
    for _, runeValue := range character_set {
	newPrefix := prefix + string(runeValue)
	produce(ch, alphabet, newPrefix, k - 1)
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

    // Stetching
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

func consume(ch <-chan string, cracks int64) {
    defer wg.Done()

    for pass := range ch {
        //Hash, check
        hash := md5crypt([]byte(pass), []byte("hfT7jp2q"), []byte("$1$"))
        atomic.AddInt64(&cracks, 1)
        if string(hash) == len6Target || string(hash) == len8Target {
            elapsed := time.Since(StartTime)
            log.Println("[!] Success: The hash corresponds with ", string(pass))
            log.Println("Cracking took ", elapsed)
            log.Println("Total passwords tried: ", atomic.LoadInt64(&cracks))
        }
    }
}

func main() {
    total := 321272406
    bar := pb.StartNew(total)
    bar.ShowCounters = true
    bar.ShowTimeLeft = true

    ch := make(chan string, 1000000000) // Buffered Channel
    var cracks int64

    for i := 0; i < 4; i++ {
        wg.Add(1)
        go consume(ch, cracks)
    }
    var runeCharArr []rune
    for i := 0; i < 24; i += 6 {
	runeCharArr = nil
	for j := 0; j < 7; j++ {
	    runeCharArr = append(runeCharArr, alphabet[i+j])
	}
	wg.Add(1)
	go Produce(ch, bar, cracks, runeCharArr, 8)
    }
    runeCharArr =  nil

    runeCharArr = append(runeCharArr, alphabet[24])
    runeCharArr = append(runeCharArr, alphabet[25])
    Produce(ch, bar, cracks, alphabet, 8)

    wg.Wait()
}
