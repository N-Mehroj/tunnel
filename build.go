package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

type Target struct {
	os   string
	arch string
	path string
}

func main() {
	osFlag := flag.String("os", "all", "Build uchun operatsion tizim (windows, linux, mac, all)")
	flag.Parse()

	targets := []Target{
		{"windows", "amd64", "builds/win/64/mytunnel-win64.exe"},
		{"windows", "386",   "builds/win/32/mytunnel-win32.exe"},
		{"linux",   "amd64", "builds/linux/64/mytunnel-linux-amd64"},
		{"linux",   "386",   "builds/linux/32/mytunnel-linux-386"},
		{"linux",   "arm",   "builds/linux/arm/mytunnel-linux-arm"},
		{"darwin",  "amd64", "builds/mac/intel/mytunnel-macos-intel"},
		{"darwin",  "arm64", "builds/mac/arm/mytunnel-macos-arm"},
	}

	for _, t := range targets {
		if *osFlag != "all" && *osFlag != t.os && !(*osFlag == "mac" && t.os == "darwin") {
			continue
		}
		dir := filepath.Dir(t.path)
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			fmt.Printf("Papka yaratishda xato (%s): %v\n", dir, err)
			continue
		}

		fmt.Printf("Building for %s (%s) -> %s\n", t.os, t.arch, t.path)
		
		cmd := exec.Command("go", "build", "-o", t.path, "client/client.go")
		cmd.Env = append(os.Environ(), "GOOS="+t.os, "GOARCH="+t.arch)
		
		cmd.Stderr = os.Stderr
		
		if err := cmd.Run(); err != nil {
			fmt.Printf("Xatolik [%s/%s]: %v\n", t.os, t.arch, err)
		}
	}
	fmt.Println("\nBuild jarayoni yakunlandi. Fayllar 'builds/' papkasida.")
}