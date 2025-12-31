package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
)

var connectionCounter int64

func main() {
        sigs := make(chan os.Signal, 1)
        signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

        controlListener, err := net.Listen("tcp", "0.0.0.0:9000")
        if err != nil {
                fmt.Printf("Control xatosi: %v\n", err)
                os.Exit(1)
        }

        externalListener, err := net.Listen("tcp", "0.0.0.0:8000")
        if err != nil {
                fmt.Printf("External xatosi: %v\n", err)
                os.Exit(1)
        }

        fmt.Println("============================================")
        fmt.Println("NGROK CLONE SERVER ISHGA TUSHDI")
        fmt.Println("Control Port (Client uchun): 9000")
        fmt.Println("Public Port (Brauzer uchun): 8000")
        fmt.Println("============================================")

        tunnels := make(chan net.Conn, 100)

        // Tunnel zaxirasini yig'ish
        go func() {
                for {
                        conn, err := controlListener.Accept()
                        if err != nil {
                                return
                        }
                        tunnels <- conn
                }
        }()

        // Tashqi so'rovlarni boshqarish
        go func() {
           for {
                        userConn, err := externalListener.Accept()
                        if err != nil {
                                return
                        }

                        id := atomic.AddInt64(&connectionCounter, 1)
                        fmt.Printf("[Req #%d] Yangi ulanish: %s\n", id, userConn.RemoteAddr())

                        select {
                        case tunnelConn := <-tunnels:
                                go handleTraffic(id, userConn, tunnelConn)
                        default:
                                fmt.Printf("[Req #%d] XATO: Bo'sh tunnel yo'q!\n", id)
                                userConn.Write([]byte("HTTP/1.1 503 Service Unavailable\r\n\r\nZaxirada tunnel yo'q."))
                                userConn.Close()
                        }
                }
        }()

        <-sigs
        fmt.Println("\nServer to'xtatilmoqda...")
}

func handleTraffic(id int64, user, tunnel net.Conn) {
        defer user.Close()
        defer tunnel.Close()

        // Ma'lumot hajmini hisoblash uchun wrapper
        copyAndLog := func(dst io.Writer, src io.Reader, direction string) {
                n, _ := io.Copy(dst, src)
                if n > 0 {
                        fmt.Printf("[Req #%d] %s: %d bytes uzatildi\n", id, direction, n)
                }
        }

        done := make(chan struct{}, 2)

        // Request (User -> Tunnel -> Laravel)
        go func() {
                copyAndLog(tunnel, user, "REQUEST ")
                done <- struct{}{}
        }()

        // Response (Laravel -> Tunnel -> User)
        go func() {
                copyAndLog(user, tunnel, "RESPONSE")
                done <- struct{}{}
        }()

        <-done
        fmt.Printf("[Req #%d] Ulanish yakunlandi.\n", id)
}


