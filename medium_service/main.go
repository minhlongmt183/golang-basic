package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
)

var counter int
var mutex = &sync.Mutex{}

func main() {
	fmt.Println("Starting....")
	//http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
	//	//fmt.Fprintf(writer, "Hello,%q", html.EscapeString(request.URL.Path))
	//	http.ServeFile(writer, request, request.URL.Path[1:])
	//})
	//
	//http.HandleFunc("/hi", func(writer http.ResponseWriter, request *http.Request) {
	//	fmt.Fprintf(writer, "Hi, I am Edisc")
	//
	//})

	//http.HandleFunc("/increment", incrementCounter)

	// tra noi dung tu 1 thu muc
	http.Handle("/", http.FileServer(http.Dir("./static")))

	//log.Fatal(http.ListenAndServe(":8081", nil))

	// https
	log.Fatal(http.ListenAndServeTLS(":443", "certs/server.crt", "certs/server.key", nil))
	fmt.Println("Stopped!")

}

func incrementCounter(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()
	counter++
	fmt.Fprintf(w, strconv.Itoa(counter))
	mutex.Unlock()
}
