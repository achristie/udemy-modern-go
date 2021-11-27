package main

import (
	"fmt"
	"net/http"

	"github.com/achristie/udemy-modern-go/handlers"
)

func main() {
	http.HandleFunc("/", handlers.Home)
	http.HandleFunc("/about", handlers.About)

	fmt.Println("Listening on port 8080")
	_ = http.ListenAndServe(":8080", nil)
}
