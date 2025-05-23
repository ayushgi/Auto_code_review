package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type Book struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
}

var books = []Book{
	{ID: 1, Title: "Go Programming", Author: "Alan"},
}

func main() {
	http.HandleFunc("/books/", bookHandler)
	log.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func bookHandler(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/books/")
	if idStr == "" || idStr == "/" {
		switch r.Method {
		case http.MethodGet:
			json.NewEncoder(w).Encode(books)
		case http.MethodPost:
			var book Book
			json.NewDecoder(r.Body).Decode(&book)
			book.ID = len(books) + 1
			books = append(books, book)
			json.NewEncoder(w).Encode(book)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
		return
	}

	id, err := strconv.Atoi(strings.Trim(idStr, "/"))
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	for i, book := range books {
		if book.ID == id {
			switch r.Method {
			case http.MethodGet:
				json.NewEncoder(w).Encode(book)
			case http.MethodPut:
				var updated Book
				json.NewDecoder(r.Body).Decode(&updated)
				updated.ID = id
				books[i] = updated
				json.NewEncoder(w).Encode(updated)
			case http.MethodDelete:
				books = append(books[:i], books[i+1:]...)
				w.WriteHeader(http.StatusNoContent)
			default:
				http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			}
			return
		}
	}

	http.Error(w, "Book Not Found", http.StatusNotFound)
}
