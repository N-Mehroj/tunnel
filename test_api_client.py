#!/usr/bin/env python3
"""
Go HTTP API Client
Ngrok orqali testing uchun
"""

import requests
import json
import sys

# Ngrok URL - o'z URL ingizni o'rnatish kerak!
# Misol: http://abc123def456.ngrok.io
BASE_URL = "http://localhost:8080"  # Lokal testing uchun

class APIClient:
    def __init__(self, base_url):
        self.base_url = base_url
        self.token = None
        self.headers = {"Content-Type": "application/json"}
        
    def register(self, email, password):
        """Ro'yxatdan o'tish"""
        url = f"{self.base_url}/auth/register"
        data = {"email": email, "password": password}
        try:
            response = requests.post(url, json=data, headers=self.headers, timeout=5)
            print(f"✓ Register: {response.status_code}")
            print(f"  Response: {response.json()}\n")
            return response.json()
        except Exception as e:
            print(f"✗ Register xatosi: {e}\n")
            return None
    
    def login(self, email, password):
        """Login"""
        url = f"{self.base_url}/auth/login"
        data = {"email": email, "password": password}
        try:
            response = requests.post(url, json=data, headers=self.headers, timeout=5)
            result = response.json()
            if result.get('success'):
                self.token = result.get('data', {}).get('token')
                self.headers['Authorization'] = f"Bearer {self.token}"
            print(f"✓ Login: {response.status_code}")
            print(f"  Response: {result}\n")
            return result
        except Exception as e:
            print(f"✗ Login xatosi: {e}\n")
            return None
    
    def get_users(self):
        """Barcha foydalanuvchilarni olish"""
        url = f"{self.base_url}/api/users"
        try:
            response = requests.get(url, headers=self.headers, timeout=5)
            print(f"✓ Get Users: {response.status_code}")
            print(f"  Response: {response.json()}\n")
            return response.json()
        except Exception as e:
            print(f"✗ Get Users xatosi: {e}\n")
            return None
    
    def create_user(self, email, password):
        """Yangi foydalanuvchi yaratish"""
        url = f"{self.base_url}/api/users"
        data = {"email": email, "password": password}
        try:
            response = requests.post(url, json=data, headers=self.headers, timeout=5)
            print(f"✓ Create User: {response.status_code}")
            print(f"  Response: {response.json()}\n")
            return response.json()
        except Exception as e:
            print(f"✗ Create User xatosi: {e}\n")
            return None
    
    def get_user(self, user_id):
        """ID bo'yicha foydalanuvchi olish"""
        url = f"{self.base_url}/api/users/{user_id}"
        try:
            response = requests.get(url, headers=self.headers, timeout=5)
            print(f"✓ Get User {user_id}: {response.status_code}")
            print(f"  Response: {response.json()}\n")
            return response.json()
        except Exception as e:
            print(f"✗ Get User xatosi: {e}\n")
            return None
    
    def update_user(self, user_id, email, password):
        """Foydalanuvchini yangilash"""
        url = f"{self.base_url}/api/users/{user_id}"
        data = {"email": email, "password": password}
        try:
            response = requests.put(url, json=data, headers=self.headers, timeout=5)
            print(f"✓ Update User {user_id}: {response.status_code}")
            print(f"  Response: {response.json()}\n")
            return response.json()
        except Exception as e:
            print(f"✗ Update User xatosi: {e}\n")
            return None
    
    def delete_user(self, user_id):
        """Foydalanuvchini o'chirish"""
        url = f"{self.base_url}/api/users/{user_id}"
        try:
            response = requests.delete(url, headers=self.headers, timeout=5)
            print(f"✓ Delete User {user_id}: {response.status_code}")
            print(f"  Response: {response.json()}\n")
            return response.json()
        except Exception as e:
            print(f"✗ Delete User xatosi: {e}\n")
            return None
    
    def health_check(self):
        """Server holatini tekshirish"""
        url = f"{self.base_url}/health"
        try:
            response = requests.get(url, timeout=5)
            print(f"✓ Health Check: {response.status_code}")
            print(f"  Response: {response.json()}\n")
            return response.json()
        except Exception as e:
            print(f"✗ Health Check xatosi: {e}\n")
            return None
    
    def logout(self):
        """Logout"""
        url = f"{self.base_url}/auth/logout"
        try:
            response = requests.post(url, headers=self.headers, timeout=5)
            print(f"✓ Logout: {response.status_code}")
            print(f"  Response: {response.json()}\n")
            return response.json()
        except Exception as e:
            print(f"✗ Logout xatosi: {e}\n")
            return None


def main():
    print("=" * 70)
    print(" " * 15 + "GO HTTP API CLIENT TEST")
    print("=" * 70 + "\n")
    
    print(f"Server: {BASE_URL}\n")
    
    client = APIClient(BASE_URL)
    
    # 1. Health check
    print("[1] Server holatini tekshirish:")
    print("-" * 70)
    client.health_check()
    
    # 2. Ro'yxatdan o'tish
    print("[2] Foydalanuvchi ro'yxatdan o'tish:")
    print("-" * 70)
    client.register("testuser@example.com", "password123")
    
    # 3. Login
    print("[3] Login qilish:")
    print("-" * 70)
    client.login("testuser@example.com", "password123")
    
    # 4. Barcha foydalanuvchilarni olish
    print("[4] Barcha foydalanuvchilarni olish:")
    print("-" * 70)
    client.get_users()
    
    # 5. Yangi foydalanuvchi yaratish
    print("[5] Yangi foydalanuvchi yaratish:")
    print("-" * 70)
    client.create_user("newuser@example.com", "pass456")
    
    # 6. ID bo'yicha foydalanuvchi olish
    print("[6] ID bo'yicha foydalanuvchi olish (ID=1):")
    print("-" * 70)
    client.get_user(1)
    
    # 7. Foydalanuvchini yangilash
    print("[7] Foydalanuvchini yangilash (ID=1):")
    print("-" * 70)
    client.update_user(1, "updated@example.com", "newpass123")
    
    # 8. Logout
    print("[8] Logout qilish:")
    print("-" * 70)
    client.logout()
    
    print("=" * 70)
    print(" " * 20 + "TEST YAKUNLANDI")
    print("=" * 70 + "\n")


if __name__ == "__main__":
    main()
