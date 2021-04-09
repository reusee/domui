// +build ignore

package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	dirFS := os.DirFS(".")
	http.Handle("/", http.FileServer(http.FS(dirFS)))
	addr := "127.0.0.1:46789"
	fmt.Printf("http://%s/demo.html\n", addr)
	http.ListenAndServe(addr, nil)
}
