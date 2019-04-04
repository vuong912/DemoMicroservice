package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func GetIdHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	data := "Book " + params["id"] + "th"
	w.Write([]byte(data))
}
func GetHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Book story...."))
}
func main() {
	router := mux.NewRouter()
	router.HandleFunc("/books", GetHandler).Methods("GET")
	router.HandleFunc("/books/{id}", GetIdHandler).Methods("GET")
	log.Fatal(http.ListenAndServe("0.0.0.0:8080", router))
}
