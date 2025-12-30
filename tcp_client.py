#!/usr/bin/env python3
"""
TCP Client - Native socket connection
"""

import socket
import json
import time

class TCPClient:
    def __init__(self, host="localhost", port=8080):
        self.host = host
        self.port = port
        self.socket = None
        
    def connect(self):
        """Serverga ulanish"""
        try:
            self.socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
            self.socket.settimeout(5)
            self.socket.connect((self.host, self.port))
            print(f"✓ {self.host}:{self.port} ga ulandi\n")
            return True
        except Exception as e:
            print(f"✗ Connection xatosi: {e}")
            print(f"  Server {self.host}:{self.port} da ishga tushganini tekshir\n")
            return False
            
    def send_request(self, method, data=None, token=None):
        """Request yuborish"""
        if not self.socket:
            print("✗ Server bilan bog'lanmagan!")
            return None
            
        message = {
            "method": method,
            "data": data or {},
            "token": token or ""
        }
        
        try:
            # Message yuborish
            json_str = json.dumps(message) + '\n'
            self.socket.send(json_str.encode())
            
            # Response qabul qilish
            response_data = self.socket.recv(4096).decode().strip()
            response = json.loads(response_data)
            
            print(f"📤 Request: {method}")
            print(f"📥 Response: {json.dumps(response, indent=2)}\n")
            return response
        except json.JSONDecodeError as e:
            print(f"✗ JSON Parse xatosi: {e}")
            return None
        except Exception as e:
            print(f"✗ Request xatosi: {e}\n")
            return None
            
    def close(self):
        """Ulanishni yopish"""
        if self.socket:
            self.socket.close()
            print("✓ Ulanish yopildi\n")


def main():
    print("=" * 70)
    print(" " * 18 + "TCP CLIENT TEST (Port 8080 -> 9000)")
    print("=" * 70 + "\n")
    
    # 8080 portga ulanish, lekin har doim 9000 ga yo'naltiriladi
    client = TCPClient("localhost", 9000)
    
    # Serverga ulanish
    if not client.connect():
        return
    
    time.sleep(0.5)
    
    # Test 1: Login
    print("[1] Login request:")
    print("-" * 70)
    client.send_request("LOGIN", {
        "email": "test@example.com",
        "password": "password123"
    })
    time.sleep(0.5)
    
    # Test 2: Register
    print("[2] Register request:")
    print("-" * 70)
    client.send_request("REGISTER", {
        "email": "newuser@example.com",
        "password": "pass456"
    })
    time.sleep(0.5)
    
    # Test 3: Get Users
    print("[3] Get Users request:")
    print("-" * 70)
    client.send_request("GET_USERS")
    time.sleep(0.5)
    
    # Test 4: Unknown method
    print("[4] Unknown method request:")
    print("-" * 70)
    client.send_request("INVALID_METHOD")
    time.sleep(0.5)
    
    # Ulanishni yopish
    # client.close()
    
    print("=" * 70)
    print(" " * 20 + "TEST YAKUNLANDI")
    print("=" * 70 + "\n")


if __name__ == "__main__":
    main()
