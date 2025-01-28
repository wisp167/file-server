package main

import (
	"os"
	"log"
	"context"
	"strconv"
	"net/http"
	"time"
	
	"github.com/joho/godotenv"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/google/uuid"
	"github.com/wisp167/file-server/internal/storage"
	"github.com/wisp167/file-server/internal/dbQueries"
	_ "github.com/lib/pq"
)

func main(){
	
	godotenv.Load(".env")
	
	PortString := os.Getenv("PORT")
	if(PortString == ""){
		log.Fatal("PORT is not found in .env")
	}
	DBport, err:=strconv.Atoi(os.Getenv("Port"))
	if err != nil{
		log.Fatal("Database port is not found in .env file")
	}
	
	
	store, err := storage.NewMySqlStorage(storage.MySqlConfig{
		Username: os.Getenv("Username"),
		Password: os.Getenv("Password"),
		DbName: os.Getenv("DbName"),
		Port: uint(DBport),
		Host: os.Getenv("Host"),
	})
	if err != nil {
		log.Fatalf("impossible to create sql storage: %s\n", err)
	}
	
	// Testing Sample
	b, err := store.DB.CreateFile(context.Background(), dbQueries.CreateFileParams{
		ID: uuid.New(),
		FileName: "Practical Go Lessons.txt",
		CreateTime: time.Now().UTC(),
		UpdateTime: time.Now().UTC(),
	})
	if err != nil{
		log.Fatalf("impossible to insert %w", err)
	}
	log.Printf("Inserted book: %v\n", b)
	//
	
	
	
	
	
	
	columns := []string{"id", "file_name", "create_time"}
	files, err := store.Select(context.Background(), columns)
	if err != nil {
		log.Fatalf("Failed to select books: %v", err)
	}
	for _, file := range files {
		log.Printf("Files: %+v\n\n", file)
	}
	
	
	
	router := chi.NewRouter()
	
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"https://*", "http://*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"*"},
		ExposedHeaders: []string{"Link"},
		MaxAge:         300,
	}))
	
	v1Router := chi.NewRouter()
	
	v1Router.Get("/healthz", handlerReadiness)
	v1Router.Get("/err", handlerErr)
	v1Router.Post("/file", store.HandlerCreateFile)
	v1Router.Get("/file_get", store.HandlerSelectFile)
	v1Router.Post("/upload", store.HandlerUploadFile)
	v1Router.Get("/download", store.HandlerGetByID)
	v1Router.Get("/list", store.HandlerListFiles)
	v1Router.Get("/file_by_name", store.HandlerGetFileByName)
	v1Router.Post("/update_file", store.HandlerUpdateFile)
	v1Router.Post("/del", store.HandlerDeleteFile)
	v1Router.Post("/del_by_name", store.HandlerDeleteFileByName)
	v1Router.Get("/count_total", store.HandlerCountFiles)
	v1Router.Get("/count_by_name", store.HandlerCountFilesByName)
	v1Router.Get("/search_by_name", store.HandlerSearchFiles)
	
	
	
	router.Mount("/v1", v1Router)
	
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!\n"))
	})
	
	srv := &http.Server{
		Handler: router,
		Addr: ":"+PortString,
	}
	log.Printf("server starting on port %v", PortString)
	
	err = srv.ListenAndServe()
	
	if err != nil{
		log.Fatal(err)
	}
	
}
