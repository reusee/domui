// +build ignore

package main

import (
	"net/http"
	"os"
)

func main() {
	dirFS := os.DirFS(".")
	http.Handle("/", http.FileServer(http.FS(dirFS)))
	http.ListenAndServe(":46789", nil)
}
