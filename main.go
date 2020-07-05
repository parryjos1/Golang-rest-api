package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// Book Struct (Model)
type Book struct {
	ID     string  `json: "id"`
	Isbn   string  `json: "isbn"`
	Title  string  `json: "title"`
	Author *Author `json: "author"`
}

// Author Struct
type Author struct {
	Firstname string `json: "firstname"`
	Lastname  string `json: "lastname"`
}

// Init books var as a slice Book struct
// a slice in go is a variable length array
var books []Book

// GET All books
func getBooks(w http.ResponseWriter, r *http.Request) {
	// above is like our response and request
	// anuy route handler function has to take in a response and request
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}

// GET Single Book
func getBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r) //Get params
	// Loop through books and find the correct ID
	for _, item := range books {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	json.NewEncoder(w).Encode(&Book{})
}

// Create a new Book
func createBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var book Book
	_ = json.NewDecoder(r.Body).Decode(&book)
	book.ID = strconv.Itoa(rand.Intn(10000000)) //Mock ID - not safe
	books = append(books, book)
	json.NewEncoder(w).Encode(book)

}

// Update a Book
// Not sure if best way but it works
func updateBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, item := range books {
		if item.ID == params["id"] {
			books = append(books[:index], books[index+1:]...)
			var book Book
			_ = json.NewDecoder(r.Body).Decode(&book)
			book.ID = params["id"]
			books = append(books, book)
			json.NewEncoder(w).Encode(book)
			return
		}
	}

}

// Delete a Book
func deleteBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, item := range books {
		if item.ID == params["id"] {
			books = append(books[:index], books[index+1:]...)
			break
		}
	}
	// respond with
	json.NewEncoder(w).Encode(books)

}

func main() {
	// Init Mux router
	router := mux.NewRouter()

	// Mock data - @todo implement database
	books = append(books, Book{ID: "1", Isbn: "4325332", Title: "Harry Potter", Author: &Author{Firstname: "J.K.", Lastname: "Rowling"}})
	books = append(books, Book{ID: "2", Isbn: "979796", Title: "Lord of the Rings", Author: &Author{Firstname: "T.", Lastname: "Tolkein"}})

	// Route handlers
	// Endpoints for our API
	router.HandleFunc("/api/books", getBooks).Methods("GET")
	router.HandleFunc("/api/books/{id}", getBook).Methods("GET")
	router.HandleFunc("/api/books", createBook).Methods("POST")
	router.HandleFunc("/api/books/{id}", updateBook).Methods("PUT")
	router.HandleFunc("/api/books/{id}", deleteBook).Methods("DELETE")

	// To run the server
	// wrapped in a log to show if it fails
	log.Fatal(http.ListenAndServe(":4001", router))

	fmt.Println("Server running on PORT:4001")
}
