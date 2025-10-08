package main

import (
    "fmt"
    "os"
    "os/signal"
    "syscall"
    "time"
)

func main() {
    sigs := make(chan os.Signal, 1)
    signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

    for {
        select {
        case sig := <-sigs:
            fmt.Println("Received signal:", sig)
            return
        default:
            // No signal received, continue doing other work
            fmt.Println("No signal yet...")
            time.Sleep(1 * time.Second)
        }
    }
}
