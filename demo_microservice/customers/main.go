package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type servicesHost struct {
	BooksHost string
}

var ServicesHost = servicesHost{
	BooksHost: "books.local",
}

func GetIdHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	data := "User " + params["id"] + "th have "

	resp, err := Request("GET", "http://192.168.99.100/books/"+params["id"], nil, ServicesHost.BooksHost)
	//resp, err := Request("GET", "http://127.0.0.1:80/books", nil, ServicesHost.BooksHost)
	if err == nil {
		w.Write([]byte(data + string(resp) + "~~~"))
	} else {
		w.Write([]byte("Error"))
	}

}

var client = &http.Client{}

func Request(method string, path string, body io.Reader, host string) ([]byte, error) {
	req, _ := http.NewRequest(method, path, body)
	if host != "" {
		req.Host = host
	}
	resp, err := client.Do(req)

	if err == nil && resp.StatusCode == 200 {
		bytes, _ := ioutil.ReadAll(resp.Body)
		return bytes, nil
	} else {
		return nil, fmt.Errorf("Some error")
	}
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/customers/{id}", GetIdHandler).Methods("GET")
	log.Fatal(http.ListenAndServe(":8080", router))
}
