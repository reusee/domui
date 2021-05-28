// +build ignore

package main

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
)

func main() {

	cmd := exec.Command("go", "build", "-o", "demo.wasm", "demo.go")
	cmd.Env = append(os.Environ(),
		"GOOS=js",
		"GOARCH=wasm",
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		panic(err)
	}

	dirFS := os.DirFS(".")
	http.Handle("/", http.FileServer(http.FS(dirFS)))
	addr := "127.0.0.1:46789"
	fmt.Printf("http://%s/demo.html\n", addr)
	http.ListenAndServe(addr, nil)
}
