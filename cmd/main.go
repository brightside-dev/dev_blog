package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/brightside-dev/dev-blog/database"
	"github.com/brightside-dev/dev-blog/database/client"
	"github.com/brightside-dev/dev-blog/internal/template"

	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

func gracefulShutdown(apiServer *http.Server, done chan bool) {
	// Create context that listens for the interrupt signal from the OS.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Listen for the interrupt signal.
	<-ctx.Done()

	log.Println("shutting down gracefully, press Ctrl+C again to force")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := apiServer.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown with error: %v", err)
	}

	log.Println("Server exiting")

	// Notify the main goroutine that the shutdown is complete
	done <- true
}

func main() {
	port, _ := strconv.Atoi(os.Getenv("HTTP_PORT"))

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      RegisterRoutes(database.New()),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	// Create a done channel to signal when the shutdown is complete
	done := make(chan bool, 1)

	// Run graceful shutdown in a separate goroutine
	go gracefulShutdown(server, done)

	err := server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		panic(fmt.Sprintf("http server error: %s", err))
	}

	// Wait for the graceful shutdown to complete
	<-done
	log.Println("Graceful shutdown complete.")
}

func RegisterRoutes(db client.DatabaseService) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	sessionManager := scs.New()
	sessionManager.Store = mysqlstore.New(db.GetDB())
	sessionManager.Lifetime = 12 * time.Hour
	sessionManager.Cookie.Secure = true

	r.Use(sessionManager.LoadAndSave)

	// Web
	fileServer := http.FileServer(http.Dir("./ui/assets/"))
	r.Handle("/assets/*", http.StripPrefix("/assets/", fileServer))

	dummyData := []map[string]any{
		{
			"id":      1,
			"title":   "I gave it a cold?",
			"excerpt": "I gave it a cold? I gave it a virus. A computer virus.",
			"content": `
				I gave it a cold? I gave it a virus. A computer virus. Yeah, but John, if The Pirates of the Caribbean breaks down, the pirates don’t eat the tourists. What do they got in there? King Kong? They're using our own satellites against us. And the clock is ticking.
				We gotta burn the rain forest, dump toxic waste, pollute the air, and rip up the OZONE! 'Cause maybe if we screw up this planet enough, they won't want it anymore! Life finds a way. Do you have any idea how long it takes those cups to decompose. Hey, take a look at the earthlings. Goodbye!
			`,
			"date": "10-10-2005 13:15",
		},
		{
			"id":      2,
			"title":   "Another title",
			"excerpt": "Another excerpt goes here...",
			"content": `
				I gave it a cold? I gave it a virus. A computer virus. Yeah, but John, if The Pirates of the Caribbean breaks down, the pirates don’t eat the tourists. What do they got in there? King Kong? They're using our own satellites against us. And the clock is ticking.
				We gotta burn the rain forest, dump toxic waste, pollute the air, and rip up the OZONE! 'Cause maybe if we screw up this planet enough, they won't want it anymore! Life finds a way. Do you have any idea how long it takes those cups to decompose. Hey, take a look at the earthlings. Goodbye!
			`,
			"date": "11-10-2005 14:20",
		},
	}

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		data := template.NewTemplateData(r, sessionManager)
		data.Data = &dummyData
		template.Render(w, r, "home", data)
	})

	r.Get("/article/{id}", func(w http.ResponseWriter, r *http.Request) {
		idParam := chi.URLParam(r, "id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			http.Error(w, "Invalid article ID", http.StatusBadRequest)
			return
		}

		var article map[string]any
		// Loop through dummyData to find the article by id
		for _, a := range dummyData {
			if a["id"] == id {
				article = a
				break
			}
		}

		// If the article doesn't exist, return a 404
		if len(article) == 0 {
			http.NotFound(w, r)
			return
		}

		// Prepare the article data with a structured type
		articleData := struct {
			Title   string
			Content string
			Date    string
		}{
			Title:   article["title"].(string),
			Content: article["content"].(string),
			Date:    article["date"].(string),
		}

		// Prepare the data to be passed to the template
		data := template.NewTemplateData(r, sessionManager)
		data.Data = articleData

		// Render the template
		template.Render(w, r, "post", data)
	})

	return r
}
