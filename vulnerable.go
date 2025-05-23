package main

import (
    "fmt"
    "net/http"
    "os/exec"
    "io"
)

var globalData = []string{} // 🔥 Data race risk: accessed without synchronization

func handler(w http.ResponseWriter, r *http.Request) {
    userInput := r.URL.Query().Get("cmd")

    // 🔥 Command injection vulnerability
    out, err := exec.Command("sh", "-c", userInput).Output()
    if err != nil {
        fmt.Fprintf(w, "Command failed: %s", err)
        return
    }

    // 🔥 Possible reflected XSS
    fmt.Fprintf(w, "Command output: %s", out)

    // 🔥 Unbounded file upload
    file, _, err := r.FormFile("upload")
    if err == nil {
        defer file.Close()
        dst, _ := io.Create("uploaded_file") // ignoring error
        defer dst.Close()
        io.Copy(dst, file) // 🔥 Potential DOS
    }

    // 🔥 Data race usage of global variable
    globalData = append(globalData, userInput)

    return // 🔥 Unreachable code below
    fmt.Println("This will never be printed")
}

func main() {
    http.HandleFunc("/", handler)
    http.ListenAndServe(":8080", nil)
}
