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


package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

var users = []User{
	{ID: 1, Name: "Alice", Email: "alice@example.com"},
	{ID: 2, Name: "Bob", Email: "bob@example.com"},
}

func main() {
	router := gin.Default()

	router.GET("/users", getUsers)
	router.POST("/users", createUser)
	router.GET("/users/:id", getUserByID)
	router.PUT("/users/:id", updateUser)
	router.DELETE("/users/:id", deleteUser)

	router.Run(":8080")
}

func getUsers(c *gin.Context) {
	c.JSON(http.StatusOK, users)
}

func createUser(c *gin.Context) {
	var newUser User
	if err := c.BindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	newUser.ID = len(users) + 1
	users = append(users, newUser)
	c.JSON(http.StatusCreated, newUser)
}

func getUserByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	for _, user := range users {
		if user.ID == id {
			c.JSON(http.StatusOK, user)
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
}

func updateUser(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	for i, user := range users {
		if user.ID == id {
			var updatedUser User
			if err := c.BindJSON(&updatedUser); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			updatedUser.ID = id
			users[i] = updatedUser
			c.JSON(http.StatusOK, updatedUser)
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
}

func deleteUser(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	for i, user := range users {
		if user.ID == id {
			users = append(users[:i], users[i+1:]...)
			c.JSON(http.StatusOK, gin.H{"message": "User deleted"})
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
}
