// smackapi/main.go
package main

import (
	"fmt"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {

    fmt.Fprintf(w, "Hello from server: This is prod\n")
}

func main() {
    http.HandleFunc("/", handler)
    port := "8080"
    fmt.Printf("Listening on port %s\n", port)
    http.ListenAndServe(":" + port, nil)
}
