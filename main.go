package main

import (
	"html/template"
	"log"
	"net/http"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type BlogPost struct {
	gorm.Model
	Title string
	Body  string
}

func main() {

	http.HandleFunc("/", LandingPageHandeler)
	http.HandleFunc("/submit_blog/", submitBlogPost)
	http.HandleFunc("/add_blog/", func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("add_blog.html"))
		tmpl.Execute(w, nil)
	})

	http.ListenAndServe(":80", nil)

}
func LandingPageHandeler(w http.ResponseWriter, r *http.Request) {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	db.AutoMigrate(&BlogPost{})
	var blogPosts []BlogPost
	result := db.Find(&blogPosts)
	if result.Error != nil {
		http.Error(w, "Error retrieving blog posts", http.StatusInternalServerError)
		return
	}

	tmpl := template.Must(template.ParseFiles("base.html"))

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, blogPosts)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func submitBlogPost(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse form values
	r.ParseForm()
	title := r.FormValue("title")
	content := r.FormValue("content")

	// Create a new BlogPost
	blogPost := BlogPost{Title: title, Body: content}

	// Connect to DB
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed db connection %v", err)
	}
	//Migrate DB
	err = db.AutoMigrate(&BlogPost{})
	if err != nil {
		log.Fatalf("failed to auto-migrate: %v", err)
	}
	// Save to database
	db.Create(&blogPost)

	// Redirect or respond
	http.Redirect(w, r, "/add_blog/", http.StatusSeeOther)
}
