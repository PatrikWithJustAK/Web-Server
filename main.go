package main

import (
	"html/template"
	"log"
	"net/http"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type BlogPost struct {
	gorm.Model
	Title string
	Body  string
}

var (
	db   *gorm.DB
	tmpl *template.Template
)

func main() {
	timestamp := time.Now().Format(time.RFC3339)
	log.Printf("HTTP Server started running at %s", timestamp)
	var err error
	db, err = gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database", err)
	}
	db.AutoMigrate(&BlogPost{})
	tmpl = template.Must(template.ParseFiles("base.html", "add_blog.html"))
	http.HandleFunc("/", LandingPageHandeler)
	http.HandleFunc("/submit_blog/", submitBlogPost)
	http.HandleFunc("/add_blog/", addBlogHandler)
	http.HandleFunc("/delete_blog/", deleteBlogHandler)
	log.Fatal(http.ListenAndServe(":80", nil))

}

func LandingPageHandeler(w http.ResponseWriter, r *http.Request) {
	var blogPosts []BlogPost
	result := db.Find(&blogPosts)
	if result.Error != nil {
		http.Error(w, "Error retrieving blog posts", http.StatusInternalServerError)
		return
	}
	tmpl.ExecuteTemplate(w, "base.html", blogPosts)
}

func submitBlogPost(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	logDBConnection(r)
	r.ParseForm()
	title := r.FormValue("title")
	content := r.FormValue("content")

	blogPost := BlogPost{Title: title, Body: content}
	result := db.Create(&blogPost)
	if result.Error != nil {
		http.Error(w, "Error saving blog post", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func addBlogHandler(w http.ResponseWriter, r *http.Request) {

	tmpl := template.Must(template.ParseFiles("add_blog.html"))
	tmpl.Execute(w, nil)

}

func logDBConnection(r *http.Request) {
	ip := r.RemoteAddr
	timestamp := time.Now().Format(time.RFC3339)
	log.Printf("DB connected successfully at %s from IP: %s\n", timestamp, ip)
}
func deleteBlogHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	id := r.FormValue("id")
	result := db.Delete(&BlogPost{}, id)
	if result.Error != nil {
		http.Error(w, "Error deleting blog post", http.StatusInternalServerError)
		return
	}
	log.Printf("Delete succesful")
	http.Redirect(w, r, "/", http.StatusSeeOther)

}
