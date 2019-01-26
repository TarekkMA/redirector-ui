package main

import (
	"log"
	"net/http"
	"sync"

	_ "github.com/jinzhu/gorm/dialects/sqlite"

	"github.com/jinzhu/gorm"

	"github.com/go-chi/chi/middleware"

	"github.com/TarekkMA/redirector-ui/controllers"
	"github.com/TarekkMA/redirector-ui/data"
	"github.com/TarekkMA/redirector-ui/datasource"
	"github.com/gorilla/mux"
)

func main() {

	ds := datasource.RedirectsDataSource{
		DB: NewDatabase(),
	}

	r := NewRouter()
	r.Use(middleware.Recoverer)
	r.Use(middleware.Logger)

	homeController := controllers.HomeController{
		DS: ds,
	}

	r.Handle("/login", http.HandlerFunc(controllers.LoginPage)).Methods("GET")
	r.Handle("/login", http.HandlerFunc(controllers.LoginAction)).Methods("POST")
	r.Handle("/home", http.HandlerFunc(homeController.HomePage)).Methods("GET")
	r.Handle("/home", http.HandlerFunc(homeController.HomeAction)).Methods("POST")
	r.Handle("/", http.FileServer(http.Dir("./public")))
	//Run this in the background
	go http.ListenAndServe(":8877", r)
	log.Println("d")

	redirects := sync.Map{}

	ds.Subscribe(func(d []*data.Redirect) {
		redirects.Range(func(key interface{}, value interface{}) bool {
			redirects.Delete(key)
			return true
		})
		for i := 0; i < len(d); i++ {
			redirects.Store(d[i].From, d[i].To)
		}
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		to, found := redirects.Load(r.Host)

		if !found {
			log.Printf("Client %s %s -> Not Found (404)", r.RemoteAddr, r.Host)
			http.NotFound(w, r)
			return
		}

		To := to.(string) + r.URL.RequestURI()

		log.Printf("Client: %s %s -> %s\n", r.RemoteAddr, r.Host, To)
		http.Redirect(w, r, To, 302)
		return
	})

	http.ListenAndServe(":80", nil)
}

func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)

	// Server CSS, JS & Images Statically.
	router.
		PathPrefix("/public/").
		Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("."+"/public/"))))

	return router
}

func NewDatabase() *gorm.DB {
	db, err := gorm.Open("sqlite3", "data.db")
	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&data.Redirect{})

	return db
}
