//go:build ignore

package main

import (
	_ "embed"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

//go:embed wasm_exec.js
var wasm_exec_js string

var pt = fmt.Printf

func main() {

	srcFile := "demo.go"
	if len(os.Args) > 1 {
		srcFile = filepath.Clean(os.Args[1])
	}
	wasmFile := srcFile[:strings.LastIndex(srcFile, ".go")] + ".wasm"
	pt("src %s, wasm %s\n", srcFile, wasmFile)

	cmd := exec.Command("go", "build", "-o", wasmFile, srcFile)
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
	http.HandleFunc("/wasm_exec.js", func(w http.ResponseWriter, req *http.Request) {
		io.WriteString(w, wasm_exec_js)
	})
	http.HandleFunc("/demo.html", func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte(`
<!DOCTYPE html>
<html>
  <head>
    <meta charset="UTF-8">
    <script src="/wasm_exec.js"></script>
  </head>
  <body>
    <div id="app"></div>
    <script>

      (async function exec() {
        const go = new Go();
        const result = await WebAssembly.instantiateStreaming(
          fetch("` + wasmFile + `"), go.importObject);
        await go.run(result.instance);
      })()

    </script>
  </body>
</html>
    `))
	})
	addr := "127.0.0.1:46789"
	fmt.Printf("http://%s/demo.html\n", addr)
	http.ListenAndServe(addr, nil)

}
