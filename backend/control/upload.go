package Upload

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

const uploadPath = "./upload/"


type ErrorResponse struct {
	Message string `json:"message"`
}


func Upload(w http.ResponseWriter, r *http.Request) {
	
	w.Header().Set("Content-Type", "application/json")

	
	if r.Method != "POST" {
		RespondWithError(w, http.StatusMethodNotAllowed, "wrong method")
		return
	}

	
	err := r.ParseMultipartForm(32 << 20) 
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	customFilename := r.FormValue("filename")
	if customFilename == "" {
		RespondWithError(w, http.StatusBadRequest, "filename is required")
		return
	}

	
	file,handler, err := r.FormFile("video")
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	fmt.Println(handler.Filename)
	fmt.Println(handler.Size)
	fmt.Println(handler.Header)

	defer file.Close() 

	
	if err := os.MkdirAll(uploadPath, os.ModePerm); err != nil {
		RespondWithError(w, http.StatusInternalServerError, "unable to create upload directory")
		return
	}

	
	filePath := uploadPath + customFilename
	dst, err := os.Create(filePath)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "unable to create file")
		return
	}
	defer dst.Close() 

	
	if _, err := io.Copy(dst, file); err != nil {
		RespondWithError(w, http.StatusInternalServerError, "failed to save file")
		return
	}

	
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"message": "File uploaded successfully: %s"}`, filePath)
}


func RespondWithError(w http.ResponseWriter, code int, message string) {
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(ErrorResponse{Message: message})
}
