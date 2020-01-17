package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/sony/sonyflake"
)

type Article struct {
	ArticleID string  `json:"articleid"`
	Title     string  `json:"title"`
	Desc      string  `json:"desc"`
	Author    *Author `json:"author"`
}

type Author struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

var articles []Article

// Get all Articles
func getArticles(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(articles)
}

// Get single Article
func getArticle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// Get params
	params := mux.Vars(r)
	// Loop through articles and find with id
	for _, item := range articles {
		if item.ArticleID == params["articleid"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	json.NewEncoder(w).Encode(&Article{})
}

// Create a new Article
func createArticle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// New article created with body from POST request ("&" returns memory adress of variable)
	var article Article
	_ = json.NewDecoder(r.Body).Decode(&article)
	// Call genID func to assign random ArticleID
	article.ArticleID = genID(article)
	articles = append(articles, article)
	json.NewEncoder(w).Encode(article)
}

func updateArticle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, item := range articles {
		if item.ArticleID == params["articleid"] {
			// Remove article
			articles = append(articles[:index], articles[index+1:]...)
			var article Article
			_ = json.NewDecoder(r.Body).Decode(&article)
			// Call genID func to assign random ArticleID
			article.ArticleID = params["articleid"]
			// Add article to memory and give response
			articles = append(articles, article)
			json.NewEncoder(w).Encode(article)
		}
	}
}

func deleteArticle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, item := range articles {
		if item.ArticleID == params["articleid"] {
			// Remove article
			articles = append(articles[:index], articles[index+1:]...)
			break
		}
	}
	json.NewEncoder(w).Encode(articles)
}

func handleRequest() {
	// Init Router
	myRouter := mux.NewRouter().StrictSlash(true)
	// Route Handlers/Endpoints
	myRouter.HandleFunc("/api/articles", getArticles).Methods("GET")
	myRouter.HandleFunc("/api/articles/{articleid}", getArticle).Methods("GET")
	myRouter.HandleFunc("/api/articles", createArticle).Methods("POST")
	myRouter.HandleFunc("/api/articles/{articleid}", updateArticle).Methods("PUT")
	myRouter.HandleFunc("/api/articles/{articleid}", deleteArticle).Methods("DELETE")
	// Run server(Log.Fatal throws error if fail)
	err := http.ListenAndServe(":8080", myRouter)
	if err != nil {
		log.Fatalf("ListenAndServe failed with %s\n", err)
	}
}

// Generate random ID and convertToString
func genID(article Article) string {
	flake := sonyflake.NewSonyflake(sonyflake.Settings{})
	id, err := flake.NextID()
	if err != nil {
		log.Fatalf("NextID() failed with %s\n", err)
	}
	newID := strconv.FormatUint(id, 10)
	return newID
}

func main() {

	// @todo - implement DB
	// Mock data
	articles = append(articles, Article{ArticleID: "1", Title: "First Article", Desc: "First Desc",
		Author: &Author{Firstname: "Simon", Lastname: "Sjogedal"}})
	articles = append(articles, Article{ArticleID: "2", Title: "Second Article", Desc: "Second Desc",
		Author: &Author{Firstname: "Ellinor", Lastname: "Svalberg"}})

	handleRequest()
}
