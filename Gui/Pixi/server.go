/*
Serve is a very simple static file server in go
Usage:
	-p="4200": port to serve on
	-d=".":    the directory of static files to host
Navigating to http://localhost:4200 will display the index.html or directory
listing file.
*/
package main

import (
	"flag"
	"log"
	"net/http"
)

func main() {
	port := flag.String("p", "8080", "port to serve on")
	directory := flag.String("d", ".", "the directory of static file to host")
	flag.Parse()

	http.Handle("/", http.FileServer(http.Dir(*directory)))

	log.Printf("Serving %s on HTTP port: %s\n", *directory, *port)
	log.Fatal(http.ListenAndServe(":"+*port, nil))
}
