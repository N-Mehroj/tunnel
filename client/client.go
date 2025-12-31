package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"time"
	"sync/atomic"
)

var requestID int64

func main() {
	serverAddr := flag.String("server", "192.168.122.91:9000", "Server IP va port")
	localPort := flag.Int("local", 8000, "Laravel porti")
	flag.Parse()

	localAddr := fmt.Sprintf("localhost:%d", *localPort)

	fmt.Println("============================================")
	fmt.Printf("TUNNEL CLIENT: %s -> %s\n", *serverAddr, localAddr)
	fmt.Println("Monitoring yoqildi. Har bir so'rov shu yerda ko'rinadi.")
	fmt.Println("============================================")

	// Tunnel pool (zaxira ulanishlar)
	for i := 0; i < 25; i++ {
		go createTunnelWorker(*serverAddr, localAddr)
	}

	select {}
}

func createTunnelWorker(server, local string) {
	for {
		tunnel, err := net.Dial("tcp", server)
		if err != nil {
			time.Sleep(3 * time.Second)
			continue
		}

		// Serverdan birinchi bayt kelishini kutamiz (Request boshlanishi)
		buf := make([]byte, 1024)
		n, err := tunnel.Read(buf)
		if err != nil {
			tunnel.Close()
			continue
		}

		id := atomic.AddInt64(&requestID, 1)
		fmt.Printf("[Conn #%d] Tunnel ishga tushdi. Clientga yo'naltirilmoqda...\n", id)

		localConn, err := net.Dial("tcp", local)
		if err != nil {
			fmt.Printf("[Conn #%d] XATO: Clientga ulanib bo'lmadi: %v\n", id, err)
			tunnel.Close()
			continue
		}

		// Birinchi olingan paketni yuborish
		localConn.Write(buf[:n])

		// Ma'lumot almashinuvi
		handle(id, tunnel, localConn)
	}
}

func handle(id int64, tunnel, local net.Conn) {
	defer tunnel.Close()
	defer local.Close()

	done := make(chan struct{}, 2)

	// Client -> Server (Laraveldan javobni Serverga)
	go func() {
		n, _ := io.Copy(tunnel, local)
		if n > 0 {
			fmt.Printf("[Conn #%d] Clientdan javob: %d bytes serverga ketdi\n", id, n)
		}
		done <- struct{}{}
	}()

	// Server -> Client (Serverdan so'rovni Laravelga)
	go func() {
		n, _ := io.Copy(local, tunnel)
		if n > 0 {
			fmt.Printf("[Conn #%d] Serverdan so'rov: %d bytes Clientga keldi\n", id, n)
		}
		done <- struct{}{}
	}()

	<-done
	fmt.Printf("[Conn #%d] Tugallandi.\n", id)
}

