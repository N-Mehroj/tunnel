#!/usr/bin/env python3
"""
TCP Port Forwarder: 8080 -> 9000
"""
import socket
import threading

def forward(src, dst):
    try:
        while True:
            data = src.recv(4096)
            if not data:
                break
            # Print data for debugging (as hex and utf-8 if possible)
            try:
                print(f"[FORWARD] {len(data)} bytes: {data.decode(errors='replace')}")
            except Exception:
                print(f"[FORWARD] {len(data)} bytes (binary)")
            dst.sendall(data)
    except Exception:
        pass
    finally:
        src.close()
        dst.close()

def handle_client(client_sock, target_host, target_port):
    client_addr = None
    try:
        client_addr = client_sock.getpeername()
    except Exception:
        client_addr = '(unknown)'
    try:
        print(f"[NEW CONNECTION] {client_addr} -> {target_host}:{target_port}")
        server_sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        server_sock.connect((target_host, target_port))
        t1 = threading.Thread(target=forward, args=(client_sock, server_sock))
        t2 = threading.Thread(target=forward, args=(server_sock, client_sock))
        t1.start()
        t2.start()
        t1.join()
        t2.join()
        print(f"[CONNECTION CLOSED] {client_addr} -> {target_host}:{target_port}")
    except Exception as e:
        print(f"[ERROR] {e}")
        client_sock.close()

def main():
    listen_host = '127.0.0.1'
    listen_port = 8000
    target_host = '127.0.0.1'
    target_port = 9000
    
    server = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    server.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEADDR, 1)
    server.bind((listen_host, listen_port))
    server.listen(5)
    print(f"TCP Forwarder: {listen_host}:{listen_port} -> {target_host}:{target_port}")
    while True:
        client_sock, addr = server.accept()
        threading.Thread(target=handle_client, args=(client_sock, target_host, target_port), daemon=True).start()

if __name__ == "__main__":
    main()
