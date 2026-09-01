package main

import (
"fmt"
"os"
"os/exec"
)

func main() {
fmt.Println("=== RusPass Password Manager ===")
fmt.Println("Frontend UI Placeholder")
fmt.Println("")
fmt.Println("The Go Qt frontend requires significant memory to compile (4GB+ RAM recommended).")
fmt.Println("Your system currently has limited RAM which causes out-of-memory errors during compilation.")
fmt.Println("")
fmt.Println("To build the frontend, you need to:")
fmt.Println("1. Add more swap space or RAM (at least 4GB total)")
fmt.Println("2. Run: make install-go-qt (to setup Qt bindings)")
fmt.Println("3. Run: make build-frontend")
fmt.Println("")
fmt.Println("For now, the Rust backend is fully functional.")
fmt.Println("Run the backend with: ./ruspass_backend")
fmt.Println("")

if _, err := os.Stat("./ruspass_backend"); err == nil {
fmt.Println("Starting backend demo...")
cmd := exec.Command("./ruspass_backend")
cmd.Stdout = os.Stdout
cmd.Stderr = os.Stderr
if err := cmd.Run(); err != nil {
fmt.Printf("Backend error: %v\n", err)
}
} else {
fmt.Println("Backend not found. Run 'make build-backend' first.")
}
}
