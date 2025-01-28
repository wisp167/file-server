
package storage

import (
	"net/http"
	"log"
	"context"
	"encoding/json"
	"time"
	"io"
	"archive/zip"
	"bytes"
	"fmt"
	"database/sql"
	"github.com/wisp167/file-server/internal/dbQueries"
	"github.com/google/uuid"
)

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}){
	dat, err := json.Marshal(payload)
	if err != nil{
		log.Printf("Failer to marshal JSON response %v\n", payload)
		w.WriteHeader(500)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(dat)
}

func respondWithError(w http.ResponseWriter, code int, msg string){
	if code > 499{
		log.Println("Responding with 5xx error:", msg)
	}
	type errResponse struct{
		Error string "json: 'error'"
	}
	
	respondWithJSON(w, code, errResponse{Error: msg,})
}

func (s MySqlStorage) HandlerListFiles(w http.ResponseWriter, r *http.Request) {
	files, err := s.DB.ListFiles(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Couldn't retrieve files: %s", err))
		return
	}

	respondWithJSON(w, http.StatusOK, files)
}

func (s MySqlStorage) HandlerCreateFile(w http.ResponseWriter, r *http.Request){
	type parameters struct{
		Name string `json:"file_data"`
		Data []byte `json:"file_data"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil{
		respondWithError(w, 400, fmt.Sprintf("Error parsing JSON"))
		return
	}
	
	
	new_book, err := s.DB.CreateFile(r.Context(), dbQueries.CreateFileParams{
		ID: uuid.New(),
		FileName: params.Name,
		CreateTime: time.Now().UTC(),
		FileData: params.Data,
		UpdateTime: time.Now().UTC(),
	})
	if err != nil{
		respondWithError(w, 400, fmt.Sprintf("Couldn't create file: %s", err))
		return
	}
	
	respondWithJSON(w, 200, new_book)
	
}

func (s MySqlStorage) HandlerGetByID(w http.ResponseWriter, r *http.Request) {
	fileID := r.URL.Query().Get("id")
	if fileID == "" {
		http.Error(w, "File ID is required", http.StatusBadRequest)
		return
	}
	
	uuid, err := uuid.Parse(fileID)
	if err != nil {
		http.Error(w, "Invalid file ID", http.StatusBadRequest)
		return
	}

	file, err := s.DB.GetFileByID(context.Background(), uuid)
	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Disposition", "attachment; filename="+file.FileName)

	contentType := http.DetectContentType(file.FileData)
	w.Header().Set("Content-Type", contentType)

	w.Write(file.FileData)
}


func (s MySqlStorage) HandlerUploadFile(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20) // 10 MB
	if err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Unable to retrieve file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	fileData, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Unable to read file data", http.StatusInternalServerError)
		return
	}
	id := uuid.New()
	_, err = s.DB.CreateFile(context.Background(), dbQueries.CreateFileParams{
		ID:         id,
		FileName:   handler.Filename,
		FileData:   fileData,
		CreateTime: time.Now().UTC(),
		UpdateTime: time.Now().UTC(),
	})
	if err != nil {
		http.Error(w, "Failed to store file in database", http.StatusInternalServerError)
		return
	}
	log.Printf("file id: %s", id)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("File uploaded successfully"))
}



func (s MySqlStorage) HandlerSelectFile(w http.ResponseWriter, r *http.Request){
	/*
	type parameters struct{
		col []string `json:"columns"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil{
		respondWithError(w, 400, fmt.Sprintf("Error parsing JSON"))
		return
	}
	*/
	columns := r.URL.Query()["col"]
	//recode_string(columns)
	fmt.Println(columns)
	if len(columns) == 0{
		respondWithError(w, http.StatusBadRequest, "No columns specified")
		return
	}
	
	books, err := s.Select(r.Context(), columns)
	if err != nil{
		respondWithError(w, 400, fmt.Sprintf("Couldn't create file: %s", err))
		return
	}
	respondWithJSON(w, 200, books)
	
}



func (s MySqlStorage) HandlerGetFileByName(w http.ResponseWriter, r *http.Request) {
	fileName := r.URL.Query().Get("name")
	if fileName == "" {
		respondWithError(w, http.StatusBadRequest, "File name is required")
		return
	}
	nullFileName := sql.NullString{String: fileName, Valid: true}
	files, err := s.DB.GetFileByName(r.Context(), nullFileName)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Couldn't retrieve files: %s", err))
		return
	}
	if len(files) == 0 {
		respondWithError(w, http.StatusNotFound, "No files found")
		return
	}

	var zipBuffer bytes.Buffer
	zipWriter := zip.NewWriter(&zipBuffer)

	for _, file := range files {
		fileWriter, err := zipWriter.Create(file.FileName)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Couldn't create ZIP entry: %s", err))
			return
		}
		_, err = fileWriter.Write(file.FileData)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Couldn't write file data to ZIP: %s", err))
			return
		}
	}

	err = zipWriter.Close()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Couldn't finalize ZIP file: %s", err))
		return
	}

	w.Header().Set("Content-Disposition", "attachment; filename=files.zip")
	w.Header().Set("Content-Type", "application/zip")

	// Stream the ZIP file to the client
	w.Write(zipBuffer.Bytes())
}

func (s MySqlStorage) HandlerGetFileByID(w http.ResponseWriter, r *http.Request) {
	fileID := r.URL.Query().Get("id")
	if fileID == "" {
		respondWithError(w, http.StatusBadRequest, "File ID is required")
		return
	}

	uuid, err := uuid.Parse(fileID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid file ID")
		return
	}

	file, err := s.DB.GetFileByID(r.Context(), uuid)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "File not found")
		return
	}
	w.Header().Set("Content-Disposition", "attachment; filename="+file.FileName)

	contentType := http.DetectContentType(file.FileData)
	w.Header().Set("Content-Type", contentType)

	w.Write(file.FileData)
	respondWithJSON(w, http.StatusOK, file)
}


func (s MySqlStorage) HandlerUpdateFile(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		ID   string `json:"id"`
		Name string `json:"file_name"`
		Data []byte `json:"file_data"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Error parsing JSON")
		return
	}

	uuid, err := uuid.Parse(params.ID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid file ID")
		return
	}

	updatedFile, err := s.DB.UpdateFile(r.Context(), dbQueries.UpdateFileParams{
		ID:         uuid,
		FileName:   params.Name,
		FileData:   params.Data,
		UpdateTime: time.Now().UTC(),
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Couldn't update file: %s", err))
		return
	}

	respondWithJSON(w, http.StatusOK, updatedFile)
}

func (s MySqlStorage) HandlerDeleteFileByName(w http.ResponseWriter, r *http.Request) {
	fileName := r.URL.Query().Get("name")
	if fileName == "" {
		respondWithError(w, http.StatusBadRequest, "File name is required")
		return
	}
	nullFileName := sql.NullString{String: fileName, Valid: true}
	files, err := s.DB.GetFileByName(r.Context(), nullFileName)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Couldn't retrieve files: %s", err))
		return
	}
	if len(files) == 0 {
		respondWithError(w, http.StatusNotFound, "No files found")
		return
	}
	for _, val := range files{
		err = s.DB.DeleteFile(r.Context(), val.ID)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Couldn't delete file: %s", err))
			return
		}
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"message": "File(s) deleted successfully"})
}

func (s MySqlStorage) HandlerDeleteFile(w http.ResponseWriter, r *http.Request) {
	fileID := r.URL.Query().Get("id")
	if fileID == "" {
		respondWithError(w, http.StatusBadRequest, "File ID is required")
		return
	}

	uuid, err := uuid.Parse(fileID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid file ID")
		return
	}

	err = s.DB.DeleteFile(r.Context(), uuid)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Couldn't delete file: %s", err))
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"message": "File deleted successfully"})
}

func (s MySqlStorage) HandlerCountFiles(w http.ResponseWriter, r *http.Request) {
	count, err := s.DB.CountFiles(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Couldn't count files: %s", err))
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]int64{"count": count})
}

func (s MySqlStorage) HandlerCountFilesByName(w http.ResponseWriter, r *http.Request) {
	fileName := r.URL.Query().Get("name")
	if fileName == "" {
		respondWithError(w, http.StatusBadRequest, "File name is required")
		return
	}
	nullFileName := sql.NullString{String: fileName, Valid: true}
	count, err := s.DB.CountFilesByName(r.Context(), nullFileName)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Couldn't count files: %s", err))
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]int64{"count": count})
}

func (s MySqlStorage) HandlerSearchFiles(w http.ResponseWriter, r *http.Request) {
	fileName := r.URL.Query().Get("name")
	if fileName == "" {
		respondWithError(w, http.StatusBadRequest, "File name is required")
		return
	}
	nullFileName := sql.NullString{String: fileName, Valid: true}
	files, err := s.DB.SearchFiles(r.Context(), nullFileName)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Couldn't retrieve files: %s", err))
		return
	}

	respondWithJSON(w, http.StatusOK, files)
}
