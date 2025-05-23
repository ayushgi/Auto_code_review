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

// Shared global slice without synchronization → data race risk
var books = []Book{
	{ID: 1, Title: "Go Programming", Author: "Alan"},
}

func main() {
	http.HandleFunc("/books/", bookHandler)
	log.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func bookHandler(w http.ResponseWriter, r *http.Request) {
	// Unsafe URL path parsing without validation
	idStr := strings.TrimPrefix(r.URL.Path, "/books/")
	if idStr == "" || idStr == "/" {
		switch r.Method {
		case http.MethodGet:
			// No pagination or filtering → potential DoS risk with large data
			json.NewEncoder(w).Encode(books)
		case http.MethodPost:
			var book Book
			// Error from Decode ignored → malformed JSON can panic or corrupt state
			json.NewDecoder(r.Body).Decode(&book)
			// Client can send arbitrary ID → no check or override here
			// Assigning ID based on length + 1 causes duplicate IDs if items deleted
			book.ID = len(books) + 1
			books = append(books, book)
			json.NewEncoder(w).Encode(book)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
		return
	}

	// Parse ID without error check → invalid IDs cause confusing 0 lookup
	id, _ := strconv.Atoi(strings.Trim(idStr, "/"))

	for i, book := range books {
		if book.ID == id {
			switch r.Method {
			case http.MethodGet:
				json.NewEncoder(w).Encode(book)
			case http.MethodPut:
				var updated Book
				// JSON decoding error ignored again
				json.NewDecoder(r.Body).Decode(&updated)
				// ID set from path param but no validation on content fields
				updated.ID = id
				books[i] = updated
				json.NewEncoder(w).Encode(updated)
			case http.MethodDelete:
				// Unsafe slice modification without synchronization
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
