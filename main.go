package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

func home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello"))
}

func main() {
	addr := flag.String("addr", ":4000", "set the port on which to listen")
	flag.Parse()

	mux := http.NewServeMux()

	mux.HandleFunc("/", home)

	fmt.Printf("Listening on: http://localhost%v\n", *addr)
	log.Fatal(http.ListenAndServe(*addr, mux))
}
